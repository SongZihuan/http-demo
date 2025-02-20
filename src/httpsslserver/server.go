package httpsslserver

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/SongZihuan/http-demo/src/certssl"
	"github.com/SongZihuan/http-demo/src/engine"
	"github.com/SongZihuan/http-demo/src/flagparser"
	"github.com/pires/go-proxyproto"
	"net"
	"net/http"
	"sync"
	"time"
)

var HttpSSLServer *http.Server = nil
var HttpSSLListener net.Listener = nil
var HttpSSLAddress string
var HttpSSLDomain string
var HttpSSLEmail string
var HttpSSLCertDir string
var HttpSSLAliyunAccessKey string
var HttpSSLAliyunAccessSecret string

var PrivateKey crypto.PrivateKey
var Certificate *x509.Certificate
var IssuerCertificate *x509.Certificate

var ErrStop = fmt.Errorf("https server error")
var ReloadMutex sync.Mutex

func InitHttpSSLServer() (err error) {
	HttpSSLAddress = flagparser.HttpsAddress
	HttpSSLDomain = flagparser.HttpsDomain
	HttpSSLEmail = flagparser.HttpsEmail
	HttpSSLCertDir = flagparser.HttpsCertDir
	HttpSSLAliyunAccessKey = flagparser.HttpsAliyunKey
	HttpSSLAliyunAccessSecret = flagparser.HttpsAliyunSecret

	PrivateKey, Certificate, IssuerCertificate, err = certssl.GetCertificateAndPrivateKey(HttpSSLCertDir, HttpSSLEmail, HttpSSLAliyunAccessKey, HttpSSLAliyunAccessSecret, HttpSSLDomain)
	if err != nil {
		return fmt.Errorf("init htttps cert ssl server error: %s", err.Error())
	} else if PrivateKey == nil || Certificate == nil || IssuerCertificate == nil {
		return fmt.Errorf("init https server error: get key and cert error, return nil, unknown reason")
	}

	err = initHttpSSLServer()
	if err != nil {
		return fmt.Errorf("init htttps error: %s", err.Error())
	}

	return nil
}

func initHttpSSLServer() (err error) {
	if PrivateKey == nil || Certificate == nil || IssuerCertificate == nil {
		return fmt.Errorf("init https server error: get key and cert error, return nil, unknown reason")
	}

	if Certificate.Raw == nil || len(Certificate.Raw) == 0 || IssuerCertificate.Raw == nil || len(IssuerCertificate.Raw) == 0 {
		return fmt.Errorf("init https server error: get cert.raw error, return nil, unknown reason")
	}

	HttpSSLServer = &http.Server{
		Addr:    HttpSSLAddress,
		Handler: engine.Engine,
	}

	return nil
}

func RunServer() error {
	watchStopChan := make(chan bool)
	defer close(watchStopChan)

	watchCertificate(watchStopChan)

	err := runServer()
	if err != nil {
		return err
	}

	return nil
}

func StopServer() (err error) {
	if HttpSSLServer == nil {
		return nil
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	err = HttpSSLServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	HttpSSLServer = nil
	HttpSSLListener = nil

	return nil
}

func loadListener() (err error) {
	if PrivateKey == nil || Certificate == nil || IssuerCertificate == nil {
		return fmt.Errorf("init https server error: get key and cert error, return nil, unknown reason")
	}

	if Certificate.Raw == nil || len(Certificate.Raw) == 0 || IssuerCertificate.Raw == nil || len(IssuerCertificate.Raw) == 0 {
		return fmt.Errorf("init https server error: get cert.raw error, return nil, unknown reason")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{Certificate.Raw, IssuerCertificate.Raw}, // Raw包含 DER 编码的证书
			PrivateKey:  PrivateKey,
			Leaf:        Certificate,
		}},
		MinVersion: tls.VersionTLS12,
	}

	tcpListener, err := net.Listen("tcp", HttpSSLServer.Addr)
	if err != nil {
		return err
	}

	proxyListener := &proxyproto.Listener{
		Listener:          tcpListener,
		ReadHeaderTimeout: 10 * time.Second,
	}

	tlsListener := tls.NewListener(proxyListener, tlsConfig)
	HttpSSLListener = tlsListener

	return nil
}

func runServer() error {
	defer func() {
		HttpSSLServer = nil
		HttpSSLListener = nil
	}()

	for {
		err := func() error {
			err := loadListener()
			if err != nil {
				return err
			}
			defer func() {
				_ = HttpSSLListener.Close()
			}()

			fmt.Printf("https server start at %s\n", HttpSSLAddress)
			err = HttpSSLServer.Serve(HttpSSLListener)
			if err != nil && errors.Is(err, http.ErrServerClosed) {
				if ReloadMutex.TryLock() {
					ReloadMutex.Unlock()
					return ErrStop
				}
				ReloadMutex.Lock()
				ReloadMutex.Unlock() // 等待证书更换完毕
				return nil
			} else if err != nil {
				return fmt.Errorf("https server error: %s", err.Error())
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}
}

func watchCertificate(stopchan chan bool) {
	newCertChan := make(chan certssl.NewCert)

	go func() {
		err := certssl.WatchCertificate(HttpSSLCertDir, HttpSSLEmail, HttpSSLAliyunAccessKey, HttpSSLAliyunAccessSecret, HttpSSLDomain, Certificate, stopchan, newCertChan)
		if err != nil {
			fmt.Printf("watch https cert server error: %s", err.Error())
		}
	}()

	go func() {
		defer close(newCertChan)

		for {
			select {
			case <-stopchan:
				return
			case res := <-newCertChan:
				if res.Error != nil {
					fmt.Printf("https cert reload server error: %s", res.Error.Error())
				} else if res.PrivateKey != nil && res.Certificate != nil && res.IssuerCertificate != nil {
					func() {
						ReloadMutex.Lock()
						defer ReloadMutex.Unlock()

						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()

						err := HttpSSLServer.Shutdown(ctx)
						if err != nil {
							fmt.Printf("https server reload shutdown error: %s", err.Error())
						}

						PrivateKey = res.PrivateKey
						Certificate = res.Certificate
						IssuerCertificate = res.IssuerCertificate
						err = initHttpSSLServer()
						if err != nil {
							fmt.Printf("https server reload init error: %s", err.Error())
						}
					}()
				}
			}
		}
	}()
}

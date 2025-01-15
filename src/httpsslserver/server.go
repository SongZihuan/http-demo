package httpsslserver

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/SongZihuan/Http-Demo/src/certssl"
	"github.com/SongZihuan/Http-Demo/src/engine"
	"github.com/SongZihuan/Http-Demo/src/flagparser"
	"net/http"
	"sync"
	"time"
)

var HttpSSLServer *http.Server = nil
var HttpSSLAddress string
var HttpSSLDomain string
var HttpSSLEmail string
var HttpSSLCertDir string

var PrivateKey crypto.PrivateKey
var Certificate *x509.Certificate

var ErrStop = fmt.Errorf("http server error")
var ReloadMutex sync.Mutex

func InitHttpSSLServer() (err error) {
	HttpSSLAddress = flagparser.HttpsAddress
	HttpSSLDomain = flagparser.HttpsDomain
	HttpSSLEmail = flagparser.HttpsEmail
	HttpSSLCertDir = flagparser.HttpsCertDir

	PrivateKey, Certificate, err = certssl.GetCertificateAndPrivateKey(HttpSSLCertDir, HttpSSLEmail, HttpSSLAddress, HttpSSLDomain)
	if err != nil {
		return err
	}

	return initHttpSSLServer()
}

func initHttpSSLServer() (err error) {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{Certificate.Raw}, // Raw包含 DER 编码的证书
			PrivateKey:  PrivateKey,
			Leaf:        Certificate,
		}},
	}

	HttpSSLServer = &http.Server{
		Addr:      HttpSSLAddress,
		Handler:   engine.Engine,
		TLSConfig: tlsConfig,
	}

	return nil
}

func RunServer() error {
	stopchan := make(chan bool)
	WatchCert(stopchan)
	err := runServer()
	stopchan <- true
	return err
}

func runServer() error {
	fmt.Printf("https server start at %s\n", HttpSSLAddress)
ListenCycle:
	for {
		err := HttpSSLServer.ListenAndServeTLS("", "")
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			if ReloadMutex.TryLock() {
				ReloadMutex.Unlock()
				return ErrStop
			}
			ReloadMutex.Lock()
			ReloadMutex.Unlock() // 等待证书更换完毕
			continue ListenCycle
		} else if err != nil {
			return err
		}
	}
}

func WatchCert(stopchan chan bool) {
	newchan := make(chan certssl.NewCert)

	go func() {
		err := certssl.WatchCertificateAndPrivateKey(HttpSSLCertDir, HttpSSLEmail, HttpSSLAddress, HttpSSLDomain, PrivateKey, Certificate, stopchan, newchan)
		if err != nil {
			fmt.Printf("watch cert error: %s", err.Error())
		}
	}()

	go func() {
		select {
		case res := <-newchan:
			if res.Certificate == nil && res.PrivateKey == nil && res.Error == nil {
				close(newchan)
				return
			} else if res.Error != nil {
				fmt.Printf("watch cert error: %s", res.Error.Error())
			} else if res.PrivateKey == nil && res.Certificate == nil {
				func() {
					ReloadMutex.Lock()
					defer ReloadMutex.Unlock()

					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					err := HttpSSLServer.Shutdown(ctx)
					if err != nil {
						fmt.Printf("reload error: %s", err.Error())
					}

					PrivateKey = res.PrivateKey
					Certificate = res.Certificate
				}()
			}
		}
	}()
}

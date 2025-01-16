package certssl

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/SongZihuan/Http-Demo/src/certssl/aliyunclear"
	"github.com/SongZihuan/Http-Demo/src/certssl/applycert"
	"github.com/SongZihuan/Http-Demo/src/utils"
	"time"
)

func InitCertSSL(aliyunAccessKey string, aliyunAccessSecret string, domain string) error {
	err := aliyunclear.InitAliyun(aliyunAccessKey, aliyunAccessSecret, domain)
	if err != nil {
		return fmt.Errorf("init aliyun failed: %s", err.Error())
	}
	return nil
}

func GetCertificateAndPrivateKey(basedir string, email string, aliyunAccessKey string, aliyunAccessSecret string, domain string) (crypto.PrivateKey, *x509.Certificate, error) {
	if email == "" {
		email = "no-reply@example.com"
	}

	if !utils.IsValidEmail(email) {
		return nil, nil, fmt.Errorf("not a valid email")
	}

	if !utils.IsValidDomain(domain) {
		return nil, nil, fmt.Errorf("not a valid domain")
	}

	privateKey, cert, err := applycert.ReadLocalCertificateAndPrivateKey(basedir)
	if err == nil && utils.CheckCertWithDomain(cert, domain) && utils.CheckCertWithTime(cert, 5*24*time.Hour) {
		return privateKey, cert, nil
	}

	privateKey, resource, err := applycert.ApplyCert(basedir, email, aliyunAccessKey, aliyunAccessSecret, domain)
	if err != nil {
		return nil, nil, fmt.Errorf("apply cert failed: %s", err.Error())
	} else if privateKey == nil || cert == nil {
		return nil, nil, fmt.Errorf("read cert failed: private key or certificate (resource) is nil, unknown reason")
	}

	cert, err = utils.ReadCertificate(resource.Certificate)
	if err != nil {
		return nil, nil, fmt.Errorf("read cert failed: %s", err.Error())
	}

	return privateKey, cert, nil
}

type NewCert struct {
	PrivateKey  crypto.PrivateKey
	Certificate *x509.Certificate
	Error       error
}

func WatchCertificateAndPrivateKey(dir string, email string, aliyunAccessKey string, aliyunAccessSecret string, domain string, oldCert *x509.Certificate, stopchan chan bool, newchan chan NewCert) error {
	for {
		select {
		case <-stopchan:
			newchan <- NewCert{
				PrivateKey:  nil,
				Certificate: nil,
				Error:       nil,
			}
			close(stopchan)
			return nil
		default:
			privateKey, cert, err := watchCertificateAndPrivateKey(dir, email, aliyunAccessKey, aliyunAccessSecret, domain, oldCert)
			if err != nil {
				newchan <- NewCert{
					Error: fmt.Errorf("watch cert failed: %s", err.Error()),
				}
			} else if privateKey != nil || cert != nil {
				newchan <- NewCert{
					PrivateKey:  privateKey,
					Certificate: cert,
				}
			}
		}
	}
}

func watchCertificateAndPrivateKey(dir string, email string, aliyunAccessKey string, aliyunAccessSecret string, domain string, oldCert *x509.Certificate) (crypto.PrivateKey, *x509.Certificate, error) {
	if email == "" {
		email = "no-reply@example.com"
	}

	if !utils.IsValidEmail(email) {
		return nil, nil, fmt.Errorf("not a valid email")
	}

	if !utils.IsValidDomain(domain) {
		return nil, nil, fmt.Errorf("not a valid domain")
	}

	if utils.CheckCertWithDomain(oldCert, domain) && utils.CheckCertWithTime(oldCert, 5*24*time.Hour) {
		return nil, nil, nil
	}

	privateKey, resource, err := applycert.ApplyCert(dir, email, aliyunAccessKey, aliyunAccessSecret, domain)
	if err != nil {
		return nil, nil, fmt.Errorf("apply cert fail: %s", err.Error())
	}

	cert, err := utils.ReadCertificate(resource.Certificate)
	if err != nil {
		return nil, nil, fmt.Errorf("read cert failed: %s", err.Error())
	}

	return privateKey, cert, nil
}

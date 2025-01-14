package acme

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/SongZihuan/Http-Demo/utils"
	"time"
)

func GetCertificateAndPrivateKey(dir string, email string, httpsAddress string, domain string) (crypto.PrivateKey, *x509.Certificate, error) {
	if email == "" {
		email = "no-reply@example.com"
	}

	if !utils.IsValidEmail(email) {
		return nil, nil, fmt.Errorf("not a valid email")
	}

	if !utils.IsValidDomain(domain) {
		return nil, nil, fmt.Errorf("not a valid domain")
	}

	privateKey, cert, err := ReadLocalCertificateAndPrivateKey(dir)
	if err != nil {
		return nil, nil, err
	}

	if checkCertWithDomain(cert, domain) && checkCertWithTime(cert, 5*24*time.Hour) {
		return privateKey, cert, nil
	}

	privateKey, resource, err := newCert(email, httpsAddress, domain)
	if err != nil {
		return nil, nil, err
	}

	err = writerWithDate(dir, resource)
	if err != nil {
		return nil, nil, err
	}

	err = writer(dir, resource)
	if err != nil {
		return nil, nil, err
	}

	cert, err = getCert(resource)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, cert, nil
}

type NewCert struct {
	PrivateKey  crypto.PrivateKey
	Certificate *x509.Certificate
	Error       error
}

func WatchCertificateAndPrivateKey(dir string, email string, httpsAddress string, domain string, oldPrivateKey crypto.PrivateKey, oldCert *x509.Certificate, stopchan chan bool, newchan chan NewCert) error {
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
			privateKey, cert, err := watchCertificateAndPrivateKey(dir, email, httpsAddress, domain, oldPrivateKey, oldCert)
			if err != nil {
				newchan <- NewCert{
					Error: err,
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

func watchCertificateAndPrivateKey(dir string, email string, httpsAddress string, domain string, oldPrivateKey crypto.PrivateKey, oldCert *x509.Certificate) (crypto.PrivateKey, *x509.Certificate, error) {
	if email == "" {
		email = "no-reply@example.com"
	}

	if !utils.IsValidEmail(email) {
		return nil, nil, fmt.Errorf("not a valid email")
	}

	if !utils.IsValidDomain(domain) {
		return nil, nil, fmt.Errorf("not a valid domain")
	}

	if checkCertWithDomain(oldCert, domain) && checkCertWithTime(oldCert, 5*24*time.Hour) {
		return nil, nil, nil
	}

	privateKey, resource, err := newCert(email, httpsAddress, domain)
	if err != nil {
		return nil, nil, err
	}

	err = writerWithDate(dir, resource)
	if err != nil {
		return nil, nil, err
	}

	err = writer(dir, resource)
	if err != nil {
		return nil, nil, err
	}

	cert, err := getCert(resource)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, cert, nil
}

func checkCertWithDomain(cert *x509.Certificate, domain string) bool {
	// 遍历主题备用名称查找匹配的域名
	for _, name := range cert.DNSNames {
		if name == domain {
			return true // 找到了匹配的域名
		}
	}

	// 检查通用名作为回退，虽然现代实践倾向于使用SAN
	if cert.Subject.CommonName != "" && cert.Subject.CommonName == domain {
		return true // 通用名匹配
	}

	// 如果没有找到匹配，则返回错误
	return false
}

func checkCertWithTime(cert *x509.Certificate, gracePeriod time.Duration) bool {
	now := time.Now()
	nowWithGracePeriod := now.Add(gracePeriod)

	if now.Before(cert.NotBefore) {
		return false
	} else if nowWithGracePeriod.After(cert.NotBefore) {
		return false
	}

	return false
}

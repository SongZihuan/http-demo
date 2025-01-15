package applycert

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/SongZihuan/Http-Demo/src/certssl/filename"
	"os"
	"path"
)

func ReadLocalCertificateAndPrivateKey(basedir string) (crypto.PrivateKey, *x509.Certificate, error) {
	cert, err := readCertificate(basedir)
	if err != nil {
		return nil, nil, fmt.Errorf("read certificate failed: %s", err.Error())
	}

	privateKey, err := readPrivateKey(basedir)
	if err != nil {
		return nil, nil, fmt.Errorf("read private key failed: %s", err.Error())
	}

	return privateKey, cert, nil
}

func readCertificate(basedir string) (*x509.Certificate, error) {
	// 请替换为你的证书文件路径
	certPath := path.Join(basedir, filename.FileCertificate)

	// 读取PEM编码的证书文件
	pemData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %v", err)
	}

	// 解析PEM编码的数据
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	// 解析证书
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return cert, nil
}

func readPrivateKey(basedir string) (crypto.PrivateKey, error) {
	// 请替换为你的RSA私钥文件路径
	keyPath := path.Join(basedir, filename.FilePrivateKey)

	// 读取PEM编码的私钥文件
	pemData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %v", err)
	}

	// 解析PEM块
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	// 根据PEM块类型解析私钥
	var privateKey crypto.PrivateKey
	if block.Type == "RSA PRIVATE KEY" {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else if block.Type == "EC PRIVATE KEY" {
		privateKey, err = x509.ParseECPrivateKey(block.Bytes)
	} else if block.Type == "PRIVATE KEY" {
		privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	} else {
		return nil, fmt.Errorf("unknown private key type: %s", block.Type)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return privateKey, nil
}

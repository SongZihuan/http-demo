package applycert

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/SongZihuan/Http-Demo/src/certssl/filename"
	"github.com/SongZihuan/Http-Demo/src/utils"
	"os"
	"path"
)

func ReadLocalCertificateAndPrivateKey(basedir string) (crypto.PrivateKey, *x509.Certificate, *x509.Certificate, error) {
	cert, err := readCertificate(basedir)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read certificate failed: %s", err.Error())
	}

	cacert, err := readCACertificate(basedir)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read certificate failed: %s", err.Error())
	}

	privateKey, err := readPrivateKey(basedir)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read private key failed: %s", err.Error())
	}

	return privateKey, cert, cacert, nil
}

func readCertificate(basedir string) (*x509.Certificate, error) {
	filepath := path.Join(basedir, filename.FileCertificate)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %v", err)
	}

	cert, err := utils.ReadCertificate(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parser certificate file: %v", err)
	}

	return cert, nil
}

func readCACertificate(basedir string) (*x509.Certificate, error) {
	filepath := path.Join(basedir, filename.FileIssuerCertificate)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %v", err)
	}

	cert, err := utils.ReadCertificate(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parser certificate file: %v", err)
	}

	return cert, nil
}

func readPrivateKey(basedir string) (crypto.PrivateKey, error) {
	filepath := path.Join(basedir, filename.FilePrivateKey)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %v", err)
	}

	privateKey, err := utils.ReadPrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parser key file: %v", err)
	}

	return privateKey, nil
}

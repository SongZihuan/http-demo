package applycert

import (
	"fmt"
	"github.com/SongZihuan/Http-Demo/src/certssl/filename"
	"github.com/SongZihuan/Http-Demo/utils"
	"github.com/go-acme/lego/v4/certificate"
	"os"
	"path"
)

func writerWithDate(baseDir string, resource *certificate.Resource) error {
	cert, err := utils.ReadCertificate(resource.Certificate)
	if err != nil {
		return err
	}

	domain := cert.Subject.CommonName
	if domain == "" && len(cert.DNSNames) == 0 {
		return fmt.Errorf("no domains in certificate")
	}
	domain = cert.DNSNames[0]

	year := fmt.Sprintf("%d", cert.NotBefore.Year())
	month := fmt.Sprintf("%d", cert.NotBefore.Month())
	day := fmt.Sprintf("%d", cert.NotBefore.Day())

	dir := path.Join(baseDir, domain, year, month, day)
	err = os.MkdirAll(dir, 0775)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, filename.FilePrivateKey), resource.PrivateKey, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, filename.FileCertificate), resource.Certificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, filename.FileIssuerCertificate), resource.IssuerCertificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, filename.FileCSR), resource.CSR, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func writer(dir string, resource *certificate.Resource) error {
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, filename.FilePrivateKey), resource.PrivateKey, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, filename.FileCertificate), resource.Certificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, filename.FileIssuerCertificate), resource.IssuerCertificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, filename.FileCSR), resource.CSR, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

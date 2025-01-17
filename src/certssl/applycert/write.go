package applycert

import (
	"fmt"
	"github.com/SongZihuan/http-demo/src/certssl/filename"
	"github.com/SongZihuan/http-demo/src/utils"
	"github.com/go-acme/lego/v4/certificate"
	"os"
	"path"
)

func writerWithDate(dir string, resource *certificate.Resource) error {
	cert, err := utils.ReadCertificate(resource.Certificate)
	if err != nil {
		return fmt.Errorf("failed to read certificate: %s", err.Error())
	}

	domain := cert.Subject.CommonName
	if domain == "" && len(cert.DNSNames) == 0 {
		return fmt.Errorf("no domains in certificate")
	}
	domain = cert.DNSNames[0]

	year := fmt.Sprintf("%d", cert.NotBefore.Year())
	month := fmt.Sprintf("%d", cert.NotBefore.Month())
	day := fmt.Sprintf("%d", cert.NotBefore.Day())

	backupdir := path.Join(dir, domain, year, month, day)
	err = os.MkdirAll(backupdir, 0775)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(backupdir, filename.FilePrivateKey), resource.PrivateKey, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(backupdir, filename.FileCertificate), resource.Certificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(backupdir, filename.FileIssuerCertificate), resource.IssuerCertificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(backupdir, filename.FileCSR), resource.CSR, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func writer(basedir string, resource *certificate.Resource) error {
	err := os.MkdirAll(basedir, 0775)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %s", basedir, err.Error())
	}

	err = os.WriteFile(path.Join(basedir, filename.FilePrivateKey), resource.PrivateKey, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(basedir, filename.FileCertificate), resource.Certificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(basedir, filename.FileIssuerCertificate), resource.IssuerCertificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(basedir, filename.FileCSR), resource.CSR, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

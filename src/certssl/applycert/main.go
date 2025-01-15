package applycert

import (
	"crypto"
	"fmt"
	"github.com/SongZihuan/Http-Demo/src/certssl/account"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"net"
	"path"
	"time"
)

func ApplyCert(basedir string, email string, httpsAddress string, domain string) (crypto.PrivateKey, *certificate.Resource, error) {
	privateKey, err := certcrypto.GeneratePrivateKey(certcrypto.RSA4096)
	if err != nil {
		return nil, nil, fmt.Errorf("generate private key failed: %s", err.Error())
	}

	user := newUser(email, privateKey)

	config := lego.NewConfig(user)
	config.Certificate.KeyType = certcrypto.RSA4096
	config.Certificate.Timeout = 30 * 24 * time.Hour
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, nil, fmt.Errorf("new client failed: %s", err.Error())
	}

	iface, port, err := net.SplitHostPort(httpsAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("split host port failed: %s", err.Error())
	} else if port == "" {
		port = "443"
	}

	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer(domain, port))
	if err != nil {
		return nil, nil, fmt.Errorf("set http01 provider failed: %s", err.Error())
	}

	reg, err := account.GetAccount(path.Join(basedir, "account"), user.GetEmail(), client)
	if err != nil {
		return nil, nil, fmt.Errorf("get account failed: %s", err.Error())
	} else if reg == nil {
		return nil, nil, fmt.Errorf("get account failed: return nil account.resurce, unknown reason")
	}
	user.setRegistration(reg)

	if domain == "" {
		domain = iface
	}

	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}

	resource, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, nil, fmt.Errorf("obtain certificate failed: %s", err.Error())
	}

	err = writerWithDate(path.Join(basedir, "cert-backup"), resource)
	if err != nil {
		return nil, nil, fmt.Errorf("writer certificate backup failed: %s", err.Error())
	}

	err = writer(basedir, resource)
	if err != nil {
		return nil, nil, fmt.Errorf("writer certificate failed: %s", err.Error())
	}

	return privateKey, resource, nil
}

package applycert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/SongZihuan/Http-Demo/src/certssl/account"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"net"
	"time"
)

func ApplyCert(basedir string, email string, httpsAddress string, domain string) (crypto.PrivateKey, *certificate.Resource, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	user := newUser(email, privateKey)

	config := lego.NewConfig(user)
	config.Certificate.KeyType = certcrypto.RSA4096
	config.Certificate.Timeout = 30 * 24 * time.Hour
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, nil, err
	}

	iface, port, err := net.SplitHostPort(httpsAddress)
	if err != nil {
		return nil, nil, err
	}

	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer(domain, port))
	if err != nil {
		return nil, nil, err
	}

	reg, err := account.GetAccount(basedir, user.GetEmail(), client)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	err = writerWithDate(basedir, resource)
	if err != nil {
		return nil, nil, err
	}

	err = writer(basedir, resource)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, resource, nil
}

package applycert

import (
	"crypto"
	"fmt"
	"github.com/SongZihuan/Http-Demo/src/certssl/account"
	"github.com/SongZihuan/Http-Demo/src/utils"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"path"
	"time"
)

func ApplyCert(basedir string, email string, aliyunAccessKey string, aliyunAccessSecret string, domain string) (crypto.PrivateKey, *certificate.Resource, error) {
	if domain == "" || !utils.IsValidDomain(domain) {
		return nil, nil, fmt.Errorf("domain is invalid")
	}

	privateKey, err := certcrypto.GeneratePrivateKey(certcrypto.RSA4096)
	if err != nil {
		return nil, nil, fmt.Errorf("generate private key failed: %s", err.Error())
	}

	user := newUser(email, privateKey)

	config := lego.NewConfig(user)
	config.Certificate.KeyType = certcrypto.RSA4096
	config.Certificate.Timeout = 30 * 24 * time.Hour
	config.CADirURL = "https://acme-v02.api.letsencrypt.org/directory"
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, nil, fmt.Errorf("new client failed: %s", err.Error())
	}

	aliyunDnsConfig := alidns.NewDefaultConfig()
	aliyunDnsConfig.APIKey = aliyunAccessKey
	aliyunDnsConfig.SecretKey = aliyunAccessSecret

	provider, err := alidns.NewDNSProviderConfig(aliyunDnsConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize AliDNS provider: %s", err.Error())
	}

	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		return nil, nil, fmt.Errorf("set challenge dns1 provider failed: %s", err.Error())
	}

	reg, err := account.GetAccount(path.Join(basedir, "account"), user.GetEmail(), client)
	if err != nil {
		return nil, nil, fmt.Errorf("get account failed: %s", err.Error())
	} else if reg == nil {
		return nil, nil, fmt.Errorf("get account failed: return nil account.resurce, unknown reason")
	}
	user.setRegistration(reg)

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

package account

import (
	"errors"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"os"
	"path"
)

func GetAccount(basedir string, email string, client *lego.Client) (*registration.Resource, error) {
	dir := path.Join(basedir, "account")
	err := os.MkdirAll(basedir, 0775)
	if err != nil {
		return nil, err
	}

	account, err := loadAccount(dir, email)
	if err != nil && errors.Is(err, ErrExpiredAccount) {
		account, err = newAccount(email, client)
		if err != nil {
			return nil, err
		}

		err = saveAccount(dir, account)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, nil
	}

	return &account.Resource, nil
}

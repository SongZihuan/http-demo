package account

import (
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"os"
)

func GetAccount(dir string, email string, client *lego.Client) (*registration.Resource, error) {
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %s", dir, err.Error())
	}

	account, err := loadAccount(dir, email)
	if err != nil && errors.Is(err, ErrExpiredAccount) {
		account, err = newAccount(email, client)
		if err != nil {
			return nil, fmt.Errorf("not local account, new account failed: %s", err.Error())
		}

		err = saveAccount(dir, account)
		if err != nil {
			return nil, fmt.Errorf("not local account, save account failed: %s", err.Error())
		}

		fmt.Printf("account register success %s\n", email)
	} else if err != nil {
		return nil, nil
	}

	return &account.Resource, nil
}

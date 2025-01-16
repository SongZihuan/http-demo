package account

import (
	"fmt"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"os"
)

var resource *registration.Resource

func GetAccount(dir string, email string, client *lego.Client) (*registration.Resource, error) {
	if resource != nil {
		return resource, nil
	}

	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %s", dir, err.Error())
	}

	account, err := loadAccount(dir, email)
	if err != nil {
		fmt.Printf("load local account failed, register a ew on for %s: %s\n", email, err.Error())

		account, err = newAccount(email, client)
		if err != nil {
			return nil, fmt.Errorf("not local account, new account failed: %s", err.Error())
		} else if account.Email == "" {
			return nil, fmt.Errorf("not local account, new account failed: return empty account, not email, unknown reason")
		}

		err = saveAccount(dir, account)
		if err != nil {
			return nil, fmt.Errorf("not local account, save account failed: %s", err.Error())
		}

		fmt.Printf("account register success for email: %s\n", email)
	} else {
		fmt.Printf("load local account success for email: %s\n", email)
	}

	resource = account.Resource
	return resource, nil
}

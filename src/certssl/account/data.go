package account

import (
	"fmt"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"time"
)

const DefaultAccountExp = 24 * time.Hour

// Account 不得包含指针
type Account struct {
	Resource       registration.Resource // 避免使用指针
	Email          string
	RegisterTime   int64
	ExpirationTime int64
}

func newAccount(email string, client *lego.Client) (Account, error) {
	res, err := register(client)
	if err != nil {
		return Account{}, fmt.Errorf("new account failed: %s", err.Error())
	} else if res == nil {
		return Account{}, fmt.Errorf("new account failed: register return nil, unknown error")
	}

	now := time.Now()
	return Account{
		Resource:       *res,
		Email:          email,
		RegisterTime:   now.Unix(),
		ExpirationTime: now.Add(DefaultAccountExp).Unix(),
	}, nil
}

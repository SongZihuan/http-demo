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
	Resource       *registration.Resource `json:"resource,omitempty"`
	Email          string                 `json:"email,omitempty"`
	RegisterTime   int64                  `json:"register-time,omitempty"`
	ExpirationTime int64                  `json:"expiration-time,omitempty"`
}

func newAccount(email string, client *lego.Client) (*Account, error) {
	res, err := register(client)
	if err != nil {
		return nil, fmt.Errorf("new account failed: %s", err.Error())
	} else if res == nil {
		return nil, fmt.Errorf("new account failed: register return nil, unknown error")
	}

	now := time.Now()
	return &Account{
		Resource:       res,
		Email:          email,
		RegisterTime:   now.Unix(),
		ExpirationTime: now.Add(DefaultAccountExp).Unix(),
	}, nil
}

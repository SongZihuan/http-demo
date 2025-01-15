package account

import (
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"time"
)

const DefaultAccountExp = 24 * time.Hour

// Account 不得包含指针
type Account struct {
	Resource       registration.Resource // 避免使用指针
	Email          string
	RegisterTime   time.Time
	ExpirationTime time.Time
}

func newAccount(email string, client *lego.Client) (Account, error) {
	res, err := register(client)
	if err != nil {
		return Account{}, err
	}

	now := time.Now()
	return Account{
		Resource:       *res,
		Email:          email,
		RegisterTime:   now,
		ExpirationTime: now.Add(DefaultAccountExp),
	}, nil
}

package acme

import (
	"crypto"
	"github.com/go-acme/lego/v4/registration"
)

type CertUser struct {
	email        string
	registration *registration.Resource
	key          crypto.PrivateKey
}

func NewUser(email string, key crypto.PrivateKey) *CertUser {
	return &CertUser{
		email: email,
		key:   key,
	}
}

func (u *CertUser) GetEmail() string {
	return u.email
}

func (u *CertUser) GetRegistration() *registration.Resource {
	return u.registration
}

func (u *CertUser) setRegistration(res *registration.Resource) {
	u.registration = res
}

func (u *CertUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

package account

import (
	"crypto"
	"encoding/json"
	"fmt"
	"github.com/SongZihuan/Http-Demo/src/utils"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"os"
	"path"
	"time"
)

const DefaultAccountExp = 24 * time.Hour

var ErrExpiredAccount = fmt.Errorf("account expired")
var ErrNotValidAccount = fmt.Errorf("account not valid")

type Data struct {
	Resource       *registration.Resource `json:"resource,omitempty"`
	Email          string                 `json:"email,omitempty"`
	RegisterTime   int64                  `json:"register-time,omitempty"`
	ExpirationTime int64                  `json:"expiration-time,omitempty"`
}

// Account 不得包含指针
type Account struct {
	data        Data
	key         crypto.PrivateKey
	accountpath string
	keypath     string
}

func NewAccount(basedir string, email string, key crypto.PrivateKey) (*Account, error) {
	dir := path.Join(basedir, "account", email)
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return nil, fmt.Errorf("create account dir failed: %s", err.Error())
	}

	now := time.Now()
	return &Account{
		data: Data{
			Email:          email,
			Resource:       nil,
			RegisterTime:   now.Unix(),
			ExpirationTime: now.Add(DefaultAccountExp).Unix(),
		},
		key:         key,
		accountpath: path.Join(dir, "account.json"),
		keypath:     path.Join(dir, "account.key"),
	}, nil
}

func LoadAccount(basedir string, email string) (*Account, error) {

	dir := path.Join(basedir, "account", email)
	accountpath := path.Join(dir, "account.json")
	keypath := path.Join(dir, "account.key")

	dataAccount, err := os.ReadFile(accountpath)
	if err != nil {
		return nil, fmt.Errorf("read account file failed: %s", err.Error())
	}

	var data Data
	err = json.Unmarshal(dataAccount, &data)
	if err != nil {
		return nil, fmt.Errorf("load account error")
	}

	dataKey, err := os.ReadFile(keypath)
	if err != nil {
		return nil, fmt.Errorf("read account key file failed: %s", err.Error())
	}

	privateKey, err := utils.ReadPrivateKey(dataKey)
	if err != nil {
		return nil, fmt.Errorf("read account key failed: %s", err.Error())
	}

	if time.Now().After(time.Unix(data.ExpirationTime, 0)) {
		return nil, ErrExpiredAccount
	}

	if data.Resource == nil || data.Resource.Body.Status != "valid" {
		return nil, ErrNotValidAccount
	}

	return &Account{
		data:        data,
		key:         privateKey,
		accountpath: accountpath,
		keypath:     keypath,
	}, nil
}

func (u *Account) GetEmail() string {
	return u.data.Email
}

func (u *Account) GetRegistration() *registration.Resource {
	return u.data.Resource
}

func (u *Account) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func (u *Account) SaveAccount() error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	err = os.WriteFile(u.accountpath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write account %s: %s", u.accountpath, err.Error())
	}

	privateKeyData, err := utils.EncodePrivateKeyToPEM(u.key)
	if err != nil {
		return fmt.Errorf("failed to read account private %s: %s", u.accountpath, err.Error())
	}

	err = os.WriteFile(u.keypath, privateKeyData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write account %s: %s", u.keypath, err.Error())
	}

	return nil
}

func (u *Account) Register(client *lego.Client) (*registration.Resource, error) {
	if u.data.Resource != nil {
		return u.data.Resource, nil
	}

	res, err := register(client)
	if err != nil {
		return nil, fmt.Errorf("new account failed: %s", err.Error())
	} else if res == nil {
		return nil, fmt.Errorf("new account failed: register return nil, unknown error")
	}

	u.data.Resource = res
	return u.data.Resource, nil
}

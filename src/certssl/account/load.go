package account

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"time"
)

var ErrExpiredAccount = fmt.Errorf("account not found")

func loadAccount(dir string, email string) (Account, error) {
	filepath := path.Join(dir, fmt.Sprintf("%s.account", email))

	file, err := os.Open(filepath)
	if err != nil {
		return Account{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	var account Account
	dec := gob.NewDecoder(file)

	err = dec.Decode(&account)
	if err != nil {
		return Account{}, err
	}

	if time.Now().After(account.ExpirationTime) {
		return Account{}, ErrExpiredAccount
	}

	return account, nil
}

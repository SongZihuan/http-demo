package account

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

var ErrExpiredAccount = fmt.Errorf("account not found")

func loadAccount(dir string, email string) (*Account, error) {
	filepath := path.Join(dir, fmt.Sprintf("%s.account.json", email))
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("read account file failed: %s", err.Error())
	}

	account := new(Account)
	err = json.Unmarshal(data, account)
	if err != nil {
		return nil, fmt.Errorf("load account error")
	}

	if time.Now().After(time.Unix(account.ExpirationTime, 0)) {
		return nil, ErrExpiredAccount
	}

	return account, nil
}

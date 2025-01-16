package account

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func saveAccount(dir string, account *Account) error {
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %s", dir, err.Error())
	}
	filepath := path.Join(dir, fmt.Sprintf("%s.account.jsom", account.Email))

	data, err := json.Marshal(account)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write account %s: %s", filepath, err.Error())
	}

	return nil
}

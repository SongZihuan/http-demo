package account

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path"
)

func saveAccount(dir string, account Account) error {
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %s", dir, err.Error())
	}
	filepath := path.Join(dir, fmt.Sprintf("%s.account", account.Email))

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err = enc.Encode(account)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, buff.Bytes(), 0644)
}

package flagparser

import "fmt"

func InitFlagParser() error {
	err := initEnv()
	if err != nil {
		return fmt.Errorf("init env error: %v", err)
	}

	err = initFlag()
	if err != nil {
		return fmt.Errorf("init flag error: %v", err)
	}

	return nil
}

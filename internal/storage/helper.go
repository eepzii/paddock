package storage

import (
	"os"
	"time"
)

func tryRemoveAll(path string) error {
	var err error

	for range 3 {
		err = os.RemoveAll(path)
		if err == nil {
			return nil
		}

		if !os.IsPermission(err) && !os.IsExist(err) {
			return err
		}
		time.Sleep(250 * time.Millisecond)
	}

	return err
}

package files

import (
	"os"
	"path/filepath"
)

func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func OpenOrCreateFile(path string, flag int) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			dir := filepath.Dir(path)
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return nil, err
			}
			file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|flag, os.ModePerm)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return file, nil
}

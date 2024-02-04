package config

import (
	"os"

	"github.com/GDGVIT/configgy-backend/constants"
)

var Password *constants.PasswordConfig

func LoadPasswordConfig() {
	panic("not implemented")
}

func LoadFileStoragePath() string {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	if filePath == "" {
		filePath = constants.DefaultFileStoragePath
	}
	_, err := os.Stat(filePath)
	if err != os.ErrNotExist {
		os.Mkdir(filePath, 0755)
	}
	return filePath
}

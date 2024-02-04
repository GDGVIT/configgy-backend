package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"strings"

	"fmt"

	"github.com/GDGVIT/configgy-backend/config"
	"github.com/GDGVIT/configgy-backend/constants"
	"golang.org/x/crypto/argon2"
)

func HashPassword(password string, salt []byte) (string, error) {
	passwordConfig := config.Password
	if passwordConfig == nil {
		passwordConfig = &constants.DefaultConfig
	}
	saltLength := passwordConfig.SaltLength

	if len(salt) == 0 {
		salt = make([]byte, saltLength)
		if _, err := rand.Read(salt); err != nil {
			return "", err
		}
	}
	hash := argon2.IDKey([]byte(password), salt, passwordConfig.Iterations, passwordConfig.Memory, passwordConfig.Parallelism, passwordConfig.HashLength)

	return getFormattedPassword(constants.PasswordHash{
		Hash:   hash,
		Salt:   salt,
		Config: *passwordConfig,
	}), nil
}

func HashBytes(secret []byte, salt []byte) (string, error) {
	passwordConfig := config.Password
	if passwordConfig == nil {
		passwordConfig = &constants.DefaultConfig
	}
	saltLength := config.Password.SaltLength

	if len(salt) == 0 {
		salt = make([]byte, saltLength)
		if _, err := rand.Read(salt); err != nil {
			return "", err
		}
	}
	hash := argon2.IDKey(secret, salt, config.Password.Iterations, config.Password.Memory, config.Password.Parallelism, config.Password.HashLength)

	return getFormattedPassword(constants.PasswordHash{
		Hash:   hash,
		Salt:   salt,
		Config: *passwordConfig,
	}), nil
}

func getFormattedPassword(passwordHash constants.PasswordHash) string {
	b64Hash := base64.StdEncoding.EncodeToString(passwordHash.Hash)
	b64Salt := base64.StdEncoding.EncodeToString(passwordHash.Salt)
	return fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		passwordHash.Config.Algorithm,
		passwordHash.Config.Argon2Version,
		passwordHash.Config.Memory,
		passwordHash.Config.Iterations,
		passwordHash.Config.Parallelism,
		b64Salt,
		b64Hash)
}

func GetParamsFromPassword(password string) (constants.PasswordHash, error) {
	var passwordHash constants.PasswordHash
	var err error

	passwordConfig := config.Password
	if passwordConfig == nil {
		passwordConfig = &constants.DefaultConfig
	}

	splitPassword := strings.Split(password, "$")
	algorithm := splitPassword[1]
	version, err := strconv.ParseUint(splitPassword[2][2:], 10, 32)
	if err != nil {
		return passwordHash, err
	}
	params := strings.Split(splitPassword[3], ",")
	if err != nil {
		return passwordHash, err
	}
	memory, err := strconv.ParseUint(params[0][2:], 10, 32)
	if err != nil {
		return passwordHash, err
	}
	iterations, err := strconv.ParseUint(params[1][2:], 10, 32)
	if err != nil {
		return passwordHash, err
	}
	parallelism, err := strconv.ParseUint(params[2][2:], 10, 8)
	if err != nil {
		return passwordHash, err
	}
	salt := splitPassword[4]
	hash := splitPassword[5]

	passwordHash.Config.Algorithm = algorithm
	passwordHash.Config.Argon2Version = uint32(version)
	passwordHash.Config.Memory = uint32(memory)
	passwordHash.Config.Iterations = uint32(iterations)
	passwordHash.Config.Parallelism = uint8(parallelism)
	passwordHash.Salt, err = base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return passwordHash, err
	}
	passwordHash.Hash, err = base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return passwordHash, err
	}

	return passwordHash, nil
}

func VerifyPassword(passwordHash string, password string) (bool, error) {
	passwordConfig, err := GetParamsFromPassword(passwordHash)
	if err != nil {
		return false, err
	}

	challengeHash, err := HashPassword(password, passwordConfig.Salt)
	if err != nil {
		return false, err
	}

	return challengeHash == passwordHash, nil
}

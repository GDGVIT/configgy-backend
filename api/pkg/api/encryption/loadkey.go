package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
	"os"
)

func GetSecretKey() cipher.Block {
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		log.Fatal("SECRET_KEY not found")
	}

	cipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatal("Error in creating cipher", err)
	}
	return cipher
}

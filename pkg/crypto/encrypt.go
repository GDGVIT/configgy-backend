package crypto

import (
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func EncryptBytes(data []byte, key cipher.Block) ([]byte, error) {
	// making a random initialization vector
	iv := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, err
	}
	// Use Galois/Counter Mode (GCM) for authenticated encryption
	aesGCM, err := cipher.NewGCM(key)
	if err != nil {
		return []byte{}, err
	}
	// Encrypt the plaintext with the AES-GCM cipher
	encryptedData := aesGCM.Seal(nil, iv, data, nil)
	// Combine IV and ciphertext
	return append(iv, encryptedData...), nil
}

func DecryptBytes(data []byte, key cipher.Block) ([]byte, error) {
	// retrieve the initialization vector and the encrypted data
	iv := data[:12]
	encryptedData := data[12:]

	// Use Galois/Counter Mode (GCM) for authenticated decryption
	aesGCM, err := cipher.NewGCM(key)
	if err != nil {
		return []byte{}, err
	}
	// Decrypt the ciphertext
	decryptedData, err := aesGCM.Open(nil, iv, encryptedData, nil)
	if err != nil {
		return []byte{}, err
	}
	return decryptedData, nil
}

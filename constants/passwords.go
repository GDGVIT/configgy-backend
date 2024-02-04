package constants

import "golang.org/x/crypto/argon2"

type PasswordHash struct {
	Hash   []byte
	Salt   []byte
	Config PasswordConfig
}

// specify the password hashing config
type PasswordConfig struct {
	HashLength    uint32
	SaltLength    uint32
	Iterations    uint32
	Memory        uint32
	Parallelism   uint8
	Algorithm     string
	Argon2Version uint32
}

var DefaultConfig = PasswordConfig{
	HashLength:    32,
	SaltLength:    16,
	Iterations:    1,
	Memory:        2 * 2 * 1024,
	Parallelism:   4,
	Algorithm:     "argon2id",
	Argon2Version: argon2.Version,
}

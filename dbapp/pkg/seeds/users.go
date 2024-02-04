package seeds

import (
	"github.com/GDGVIT/configgy-backend/pkg/crypto"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
	"gorm.io/gorm"
)

func Users(db *gorm.DB) error {
	// Seed 1
	seedPassword, err := crypto.HashPassword("password", []byte{})
	if err != nil {
		return err
	}
	err = db.Create(&tables.Users{
		PID:        tables.UUIDWithPrefix("usr"),
		Name:       "Admin",
		Email:      "admin@admin.com",
		Password:   []byte(seedPassword),
		IsAdmin:    true,
		IsVerified: true,
		IsEnabled:  true,
		PublicKey:  []byte("publickey"),
	}).Error

	if err != nil {
		return err
	}

	return nil
}

package pkg

import (
	"github.com/GDGVIT/configgy-backend/pkg/tables"
	"gorm.io/gorm"
)

type Migrate struct {
	TableName string
	Run       func(*gorm.DB) error
}

func AutoMigrate(db *gorm.DB) []Migrate {
	var users tables.Users
	var userVerifications tables.UserVerification
	var vaults tables.Vault
	var vaultCredentials tables.VaultCredentials
	var permissionAssignments tables.PermissionAssignments
	var credentials tables.Credential
	var passwordCredentials tables.PasswordCredentials
	var fileCredentials tables.FileCredentials
	var featureFlagCredentials tables.FeatureFlagCredentials
	var featureFlagData tables.FeatureFlagData

	usersM := Migrate{TableName: "users",
		Run: func(d *gorm.DB) error { return db.AutoMigrate(&users) }}
	userVerificationsM := Migrate{TableName: "user_verifications",
		Run: func(d *gorm.DB) error { return db.AutoMigrate(&userVerifications) }}
	vaultsM := Migrate{TableName: "vaults",
		Run: func(d *gorm.DB) error { return db.AutoMigrate(&vaults) }}
	vaultCredentialsM := Migrate{TableName: "vault_credentials",
		Run: func(d *gorm.DB) error { return db.AutoMigrate(&vaultCredentials) }}
	permissionAssignmentsM := Migrate{TableName: "permission_assignments",
		Run: func(d *gorm.DB) error { return db.AutoMigrate(&permissionAssignments) }}
	credentialsM := Migrate{TableName: "credentials",
		Run: func(d *gorm.DB) error { return db.AutoMigrate(&credentials) }}
	passwordCredentialsM := Migrate{TableName: "password_credentials",
		Run: func(d *gorm.DB) error { return db.AutoMigrate(&passwordCredentials) }}
	fileCredentialsM := Migrate{TableName: "file_credentials",
		Run: func(d *gorm.DB) error { return db.AutoMigrate(&fileCredentials) }}
	featureFlagCredentialsM := Migrate{TableName: "feature_flag_credentials",
		Run: func(d *gorm.DB) error { return db.AutoMigrate(&featureFlagCredentials) }}
	featureFlagDataM := Migrate{TableName: "feature_flag_data",
		Run: func(d *gorm.DB) error { return db.AutoMigrate(&featureFlagData) }}

	return []Migrate{usersM, userVerificationsM, vaultsM, vaultCredentialsM, permissionAssignmentsM, credentialsM, passwordCredentialsM, fileCredentialsM, featureFlagCredentialsM, featureFlagDataM}
}

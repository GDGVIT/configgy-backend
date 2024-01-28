package tables

import (
	"time"

	"gorm.io/gorm"
)

type Vault struct {
	ID          int    `gorm:"column:vault_id;primaryKey;autoIncrement"`
	PID         string `gorm:"column:vault_pid;unique;type:varchar(40)"`
	Name        string `gorm:"column:vault_name;not null;type:varchar(100)"`
	Description string `gorm:"column:vault_description;not null;type:varchar(20000)"`
	PublicKey   []byte `gorm:"column:public_key;not null"`
	IsPersonal  bool   `gorm:"column:is_personal;not null;default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type VaultCredentials struct {
	ID           int `gorm:"column:vault_credential_id;primaryKey;autoIncrement"`
	VaultID      int `gorm:"column:vault_id;not null"`
	CredentialID int `gorm:"column:credential_id;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (t *Vault) TableName() string {
	return "vault"
}

func (t *VaultCredentials) TableName() string {
	return "vault_credentials"
}

// Create a new vault
func (db *DB) CreateVault(vault Vault, userPID string) error {
	user := Users{}
	err := db.gormDB.Where("user_pid = ?", userPID).First(&user).Error
	if err != nil {
		return err
	}
	vaultCreateTx := db.gormDB.Create(&vault)
	if vaultCreateTx.Error != nil {
		vaultCreateTx.Rollback()
		return vaultCreateTx.Error
	}

	permissionAssignmentCreateTx := db.gormDB.Create(&PermissionAssignments{
		PermissionName: OwnerPermission,
		PID:            UUIDWithPrefix("permissionassignment"),
		VaultID:        vault.ID,
		UserID:         user.ID,
		ResourcePID:    vault.PID,
		ResourceType:   VaultResource,
		IdentityPID:    userPID,
		IdentityType:   UserIdentity,
	})
	if permissionAssignmentCreateTx.Error != nil {
		permissionAssignmentCreateTx.Rollback()
		vaultCreateTx.Rollback()
		return permissionAssignmentCreateTx.Error
	}
	return nil
}

func (db *DB) EditVault(vaultID int, vaultContent Vault) *gorm.DB {
	return db.gormDB.Model(&Vault{}).Where("vault_id = ?", vaultID).Updates(vaultContent)
}

// Delete a vault
func (db *DB) DeleteVault(vaultID int) error {
	txns := []*gorm.DB{}
	vaultCredentials := []*VaultCredentials{}
	err := db.gormDB.Where("vault_id = ?", vaultID).Find(&vaultCredentials).Error
	if err != nil {
		return err
	}
	// delete all credential data that is in the vault
	for _, vaultCredential := range vaultCredentials {
		tx := db.gormDB.Where("credential_id = ?", vaultCredential.CredentialID).Delete(&Credential{})
		if tx.Error != nil {
			db.RollbackTxns(txns)
			return tx.Error
		}
		txns = append(txns, tx)
	}

	vaultCredentialDeleteTx := db.gormDB.Where("vault_id = ?", vaultID).Delete(&VaultCredentials{})
	if vaultCredentialDeleteTx.Error != nil {
		db.RollbackTxns(txns)
		return vaultCredentialDeleteTx.Error
	}
	txns = append(txns, vaultCredentialDeleteTx)

	vaultDeleteTx := db.gormDB.Where("vault_id = ?", vaultID).Delete(&Vault{})
	if vaultDeleteTx.Error != nil {
		db.RollbackTxns(txns)
		return vaultDeleteTx.Error
	}
	txns = append(txns, vaultDeleteTx)
	// delete all permissionassignments
	permissionAssignmentDeleteTx := db.gormDB.Where("vault_id = ?", vaultID).Delete(&PermissionAssignments{})
	if permissionAssignmentDeleteTx.Error != nil {
		db.RollbackTxns(txns)
		return permissionAssignmentDeleteTx.Error
	}
	return nil
}

// Get a vault by id
func (db *DB) GetVaultByID(id int) (*Vault, error) {
	var vault Vault
	err := db.gormDB.Where("id = ?", id).First(&vault).Error
	return &vault, err
}

func (db *DB) GetVaultByPID(pid string) (*Vault, error) {
	var vault Vault
	err := db.gormDB.Where("vault_pid = ?", pid).First(&vault).Error
	return &vault, err
}

// create a new vault - credential mapping
func (db *DB) CreateVaultCredential(vault_credential *VaultCredentials) error {
	return db.gormDB.Create(vault_credential).Error
}

// Get a vault - credential mapping by id
func (db *DB) GetVaultCredentialByID(id int) (*VaultCredentials, error) {
	var vault_credential VaultCredentials
	err := db.gormDB.Where("id = ?", id).First(&vault_credential).Error
	return &vault_credential, err
}

func (db *DB) AddCredentialToVault(credentialID int, vaultID int) *gorm.DB {
	return db.gormDB.Create(&VaultCredentials{
		VaultID:      vaultID,
		CredentialID: credentialID,
	})
}

func (db *DB) GetCredentialsForVault(vaultID int) ([]*VaultCredentials, error) {
	var vaultCredentials []*VaultCredentials
	err := db.gormDB.Where("vault_id = ?", vaultID).Find(&vaultCredentials).Error
	return vaultCredentials, err
}

func (db *DB) GetVaultIDForCredential(credentialID int) (*VaultCredentials, error) {
	var vaultCredentials VaultCredentials
	err := db.gormDB.Where("credential_id = ?", credentialID).First(&vaultCredentials).Error
	if err != nil {
		return nil, err
	}
	return &vaultCredentials, nil
}

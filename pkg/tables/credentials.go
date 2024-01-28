package tables

import (
	"time"

	"gorm.io/gorm"
)

type CredentialType string

const (
	Password       CredentialType = "password"
	FileCredential CredentialType = "file"
	FeatureFlag    CredentialType = "feature_flags"
)

type Credential struct {
	ID             int            `gorm:"column:credential_id;primaryKey;autoIncrement"`
	PID            string         `gorm:"column:credential_pid;unique;type:varchar(40)"`
	CredentialName string         `gorm:"column:credential_name;not null"`
	Notes          string         `gorm:"column:notes;not null;type:varchar(20000)"`
	CredentialType CredentialType `gorm:"column:credential_type;not null"`
	CredentialID   int            `gorm:"column:credential_data_id;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type PasswordCredentials struct {
	ID               int       `gorm:"column:password_credential_id;primaryKey;autoIncrement"`
	Username         string    `gorm:"column:username;not null"`
	Password         []byte    `gorm:"column:password;not null"`
	PasswordStrength int       `gorm:"column:password_strength;not null"`
	ExpiresAt        time.Time `gorm:"column:expires_at"`
	TOTPKey          []byte    `gorm:"column:totp_key"`
	TOTPLength       int       `gorm:"column:totp_length"`
	TOTPPeriod       int       `gorm:"column:totp_period"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type FileCredentials struct {
	ID        int       `gorm:"column:file_credential_id;primaryKey;autoIncrement"`
	FilePath  string    `gorm:"column:file_path;not null"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FeatureFlagCredentials struct {
	ID                         int    `gorm:"column:feature_flag_credential_id;primaryKey;autoIncrement"`
	FeatureFlagsetName         string `gorm:"column:feature_flagset_name;not null"`
	FeatureFlagEnvironmentName string `gorm:"column:feature_flag_environment;not null"`
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

type FeatureFlagData struct {
	ID                      int    `gorm:"column:feature_flag_data_id;primaryKey;autoIncrement"`
	FeatureFlagCredentialID int    `gorm:"column:feature_flag_credential_id;not null"`
	Name                    string `gorm:"column:feature_flag_name"`
	FlagKey                 string `gorm:"column:feature_flag_key;not null"`
	FlagValue               string `gorm:"column:feature_flag_value;not null"`
}

func (t *Credential) TableName() string {
	return "credentials"
}

// Create a new credential
func (db *DB) CreateCredential(credential *Credential) *gorm.DB {
	return db.gormDB.Create(credential)
}

// Get a credential by id
func (db *DB) GetCredentialByID(id int) (*Credential, error) {
	var credential Credential
	err := db.gormDB.Where("credential_id = ?", id).First(&credential).Error
	return &credential, err
}

// Get a credential by pid
func (db *DB) GetCredentialByPID(pid string) (*Credential, error) {
	var credential Credential
	err := db.gormDB.Where("credential_pid = ?", pid).First(&credential).Error
	return &credential, err
}

func (db *DB) UpdateCredentialByID(id int, credential *Credential) *gorm.DB {
	return db.gormDB.Model(&Credential{}).Where("credential_id = ?", id).Updates(&credential)
}

func (t *PasswordCredentials) TableName() string {
	return "password_credentials"
}

// Create a new password credential
func (db *DB) CreatePasswordCredential(credential *PasswordCredentials) *gorm.DB {
	return db.gormDB.Create(credential)
}

func (db *DB) GetPasswordCredentialByID(id int) (*PasswordCredentials, error) {
	var credential PasswordCredentials
	err := db.gormDB.Where("password_credential_id = ?", id).First(&credential).Error
	return &credential, err
}

func (db *DB) UpdatePasswordCredentialByID(id int, passwordValue *PasswordCredentials) *gorm.DB {
	return db.gormDB.Model(&PasswordCredentials{}).Where("password_credential_id = ?", id).Updates(&passwordValue)
}

func (t *FileCredentials) TableName() string {
	return "file_credentials"
}

// Create a new file credential
func (db *DB) CreateFileCredential(credential *FileCredentials) *gorm.DB {
	return db.gormDB.Create(credential)
}

func (db *DB) GetFileCredentialByID(id int) (*FileCredentials, error) {
	var credential FileCredentials
	err := db.gormDB.Where("file_credential_id = ?", id).First(&credential).Error
	return &credential, err
}

func (t *FeatureFlagCredentials) TableName() string {
	return "feature_flag_credentials"
}

// Create a new feature flag credential
func (db *DB) CreateFeatureFlagCredential(credential *FeatureFlagCredentials) *gorm.DB {
	return db.gormDB.Create(credential)
}

func (db *DB) GetFeatureFlagCredentialByID(id int) (*FeatureFlagCredentials, error) {
	var credential FeatureFlagCredentials
	err := db.gormDB.Where("feature_flag_credential_id = ?", id).First(&credential).Error
	return &credential, err
}

func (t *FeatureFlagData) TableName() string {
	return "feature_flag_data"
}

func (db *DB) AddFeatureFlag(featureFlagData *FeatureFlagData) (*gorm.DB, error) {
	// check if feature flag already exists
	var existingFeatureFlagData FeatureFlagData
	err := db.gormDB.Where("feature_flag_credential_id = ? AND feature_flag_key = ?", featureFlagData.FeatureFlagCredentialID, featureFlagData.FlagKey).First(&existingFeatureFlagData).Error
	if err == gorm.ErrRecordNotFound {
		tx := db.gormDB.Create(&featureFlagData)
		return tx, tx.Error
	}
	return &gorm.DB{}, err
}

func (db *DB) GetFeatureFlagsByID(FeatureFlagCredentialsID int) ([]*FeatureFlagData, error) {
	var featureFlagData []*FeatureFlagData
	tx := db.gormDB.Where("feature_flag_credential_id = ?", FeatureFlagCredentialsID).Find(&featureFlagData)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return featureFlagData, nil
}

func (db *DB) GetFeatureFlagValueForKey(featureFlagCredentialID int, key string) (FeatureFlagData, error) {
	var featureFlagData FeatureFlagData
	err := db.gormDB.Where("feature_flag_credential_id = ? AND feature_flag_key = ?", featureFlagCredentialID, key).First(&featureFlagData).Error
	return featureFlagData, err
}

func (db *DB) EditFeatureFlagValueForKey(featureFlagCredentialID int, key string, value string) *gorm.DB {
	return db.gormDB.Model(&FeatureFlagData{}).Where("feature_flag_credential_id = ? AND feature_flag_key = ?", featureFlagCredentialID, key).Update("feature_flag_value", value)
}

func (db *DB) DeleteFeatureFlagValueForKey(featureFlagCredentialID int, key string) *gorm.DB {
	return db.gormDB.Where("feature_flag_credential_id = ? AND feature_flag_key = ?", featureFlagCredentialID, key).Delete(&FeatureFlagData{})
}

func (db *DB) AddMultipleFeatureFlags(featureFlagCredentialID int, featureFlagData []*FeatureFlagData) *gorm.DB {
	return db.gormDB.Create(featureFlagData)
}

func (db *DB) GetCredentialDataForCredential(credentialID int, credentialType CredentialType) (interface{}, error) {
	switch credentialType {
	case Password:
		return db.GetPasswordCredentialByID(credentialID)
	case FileCredential:
		return db.GetFileCredentialByID(credentialID)
	case FeatureFlag:
		return db.GetFeatureFlagCredentialByID(credentialID)
	default:
		return nil, nil
	}
}

func (db *DB) DeleteCredentialByID(id int, credential Credential) error {
	// handles deletions in other tables to avoid orphaned resources
	var txns []*gorm.DB
	credentialVaultDeleteTx := db.gormDB.Where("credential_id = ?", id).Delete(&VaultCredentials{})
	txns = append(txns, credentialVaultDeleteTx)
	if credentialVaultDeleteTx.Error != nil {
		db.RollbackTxns(txns)
		return credentialVaultDeleteTx.Error
	}
	switch credential.CredentialType {
	case Password:
		credentialPasswordDeleteTx := db.gormDB.Where("credential_id = ?", id).Delete(&PasswordCredentials{})
		txns = append(txns, credentialPasswordDeleteTx)
		if credentialPasswordDeleteTx.Error != nil {
			db.RollbackTxns(txns)
			return credentialPasswordDeleteTx.Error
		}
	case FileCredential:
		credentialFileDeleteTx := db.gormDB.Where("credential_id = ?", id).Delete(&FileCredentials{})
		txns = append(txns, credentialFileDeleteTx)
		if credentialFileDeleteTx.Error != nil {
			db.RollbackTxns(txns)
			return credentialFileDeleteTx.Error
		}
	case FeatureFlag:
		featureFlag, err := db.GetFeatureFlagCredentialByID(id)
		if err != nil {
			db.RollbackTxns(txns)
			return err
		}
		featureFlagID := featureFlag.ID

		featureFlagCredentialDataDeleteTx := db.gormDB.Where("feature_flag_credential_id = ?", featureFlagID).Delete(&FeatureFlagData{})
		txns = append(txns, featureFlagCredentialDataDeleteTx)
		if featureFlagCredentialDataDeleteTx.Error != nil {
			db.RollbackTxns(txns)
			return featureFlagCredentialDataDeleteTx.Error
		}

		credentialFeatureFlagDeleteTx := db.gormDB.Where("credential_id = ?", id).Delete(&FeatureFlagCredentials{})
		txns = append(txns, credentialFeatureFlagDeleteTx)
		if credentialFeatureFlagDeleteTx.Error != nil {
			db.RollbackTxns(txns)
			return credentialFeatureFlagDeleteTx.Error
		}
	}
	credentialPermissionAssignmentsDeleteTx := db.gormDB.Where("credential_id = ?", id).Delete(&PermissionAssignments{})
	txns = append(txns, credentialPermissionAssignmentsDeleteTx)
	if credentialPermissionAssignmentsDeleteTx.Error != nil {
		db.RollbackTxns(txns)
		return credentialPermissionAssignmentsDeleteTx.Error
	}
	credentialDeleteTx := db.gormDB.Where("credential_id = ?", id).Delete(&Credential{})
	txns = append(txns, credentialDeleteTx)
	if credentialDeleteTx.Error != nil {
		db.RollbackTxns(txns)
		return credentialDeleteTx.Error
	}
	return nil
}

func (db *DB) RollbackTxns(txns []*gorm.DB) {
	// rolllback all transactions
	for _, txn := range txns {
		txn.Rollback()
	}
}

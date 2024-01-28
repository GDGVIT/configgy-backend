package credentialsvc

import (
	"context"
	"crypto/cipher"

	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/logger"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
	"gorm.io/gorm"
)

type CredentialServiceImpl struct {
	DB                CredentialDb
	secretKey         cipher.Block
	logger            logger.Logger
	messageBroker     MessageBroker
	accesscontrolrsvc AccessControlSvc
	authnsvc          AuthnSvc
}

type AccessControlSvc interface {
	CreatePermissionAssignment(permission_assignment *tables.PermissionAssignments) *gorm.DB
	UserHasPermissionToResource(userID, resourceID int, resourceType tables.ResourceTypes, action accesscontrolsvc.CRUDOperation) (bool, error)
	GroupHasPermissionToResource(groupID, resourceID int, resourceType tables.ResourceTypes, action accesscontrolsvc.CRUDOperation) (bool, error)
}

type AuthnSvc interface {
	ValidateToken(signedToken string) error
}

type CredentialDb interface {
	CreateCredential(credential *tables.Credential) *gorm.DB
	GetCredentialByID(id int) (*tables.Credential, error)
	GetCredentialByPID(pid string) (*tables.Credential, error)
	UpdateCredentialByID(id int, credential *tables.Credential) *gorm.DB
	CreatePasswordCredential(credential *tables.PasswordCredentials) *gorm.DB
	GetPasswordCredentialByID(id int) (*tables.PasswordCredentials, error)
	UpdatePasswordCredentialByID(id int, credential *tables.PasswordCredentials) *gorm.DB
	CreateFileCredential(credential *tables.FileCredentials) *gorm.DB
	GetFileCredentialByID(id int) (*tables.FileCredentials, error)
	CreateFeatureFlagCredential(credential *tables.FeatureFlagCredentials) *gorm.DB
	GetFeatureFlagCredentialByID(id int) (*tables.FeatureFlagCredentials, error)
	AddFeatureFlag(featureFlagData *tables.FeatureFlagData) (*gorm.DB, error)
	GetFeatureFlagsByID(FeatureFlagCredentialsID int) ([]*tables.FeatureFlagData, error)
	GetFeatureFlagValueForKey(featureFlagCredentialID int, key string) (tables.FeatureFlagData, error)
	EditFeatureFlagValueForKey(featureFlagCredentialID int, key string, value string) *gorm.DB
	AddMultipleFeatureFlags(featureFlagCredentialID int, featureFlagData []*tables.FeatureFlagData) *gorm.DB
	DeleteFeatureFlagValueForKey(featureFlagCredentialID int, key string) *gorm.DB
	GetUserByPID(pid string) (*tables.Users, error)
	GetVaultByPID(pid string) (*tables.Vault, error)
	AddCredentialToVault(vaultID, credentialID int) *gorm.DB
	GetCredentialDataForCredential(id int, credentialType tables.CredentialType) (interface{}, error)
	DeleteCredentialByID(id int, credential tables.Credential) error

	RollbackTxns(txns []*gorm.DB)
}

type MessageBroker interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
}

func Handler(credentialDb CredentialDb, secretKey cipher.Block, logger logger.Logger, messageBroker MessageBroker, accessControlrSvc AccessControlSvc, authnSvc AuthnSvc) *CredentialServiceImpl {
	return &CredentialServiceImpl{
		DB:                credentialDb,
		secretKey:         secretKey,
		logger:            logger,
		messageBroker:     messageBroker,
		authnsvc:          authnSvc,
		accesscontrolrsvc: accessControlrSvc,
	}
}

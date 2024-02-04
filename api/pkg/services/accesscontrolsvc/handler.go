package accesscontrolsvc

import (
	"github.com/GDGVIT/configgy-backend/pkg/tables"
	"gorm.io/gorm"
)

type accessControlSvcImpl struct {
	DB DB
}

func Handler(db DB) *accessControlSvcImpl {
	return &accessControlSvcImpl{DB: db}
}

type DB interface {
	GetPermissionsForUser(int) ([]*tables.PermissionAssignments, error)
	GetPermissionsForUserOnResource(userID int, resourceID int, resourceType tables.ResourceTypes) (*tables.PermissionAssignments, error)
	GetPermissionsForGroupOnResource(groupID int, resourceID int, resourceType tables.ResourceTypes) (*tables.PermissionAssignments, error)
	GetPermissionsForResource(int, tables.ResourceTypes) ([]*tables.PermissionAssignments, error)
	CreatePermissionAssignment(*tables.PermissionAssignments) *gorm.DB

	GetGroupsByUserID(userID int) ([]*tables.PermissionAssignments, error)
	GetCredentialByID(id int) (*tables.Credential, error)
	GetVaultIDForCredential(credentialID int) (*tables.VaultCredentials, error)
}

type CRUDOperation string

const (
	CreateOperation CRUDOperation = "create"
	ReadOperation   CRUDOperation = "read"
	UpdateOperation CRUDOperation = "update"
	DeleteOperation CRUDOperation = "delete"
)

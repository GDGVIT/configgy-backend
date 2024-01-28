package authzsvc

import (
	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/logger"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
	"gorm.io/gorm"
)

type AuthzSvcImpl struct {
	DB               DB
	accessControlSvc AccessControlrSvc
	authnSvc         AuthnSvc
	logger           logger.Logger
}

type DB interface {
	GetUserByPID(pid string) (*tables.Users, error)
	GetGroupByPID(pid string) (*tables.Groups, error)
	GetVaultByPID(pid string) (*tables.Vault, error)
	GetCredentialByPID(pid string) (*tables.Credential, error)
	IsUserInGroup(userID int, groupID int) (bool, error)

	GetPermissionAssignmentByPID(pid string) (*tables.PermissionAssignments, error)
	GetPermissionsForResource(resourceID int, resourceType tables.ResourceTypes) ([]*tables.PermissionAssignments, error)
	GetPermissionsForUser(userID int) ([]*tables.PermissionAssignments, error)
	GetPermissionsForGroup(groupID int) ([]*tables.PermissionAssignments, error)

	CreatePermissionAssignment(permissionAssignment *tables.PermissionAssignments) *gorm.DB

	DeletePermissionAssignment(permissionAssignment *tables.PermissionAssignments) *gorm.DB

	RollbackTxns(txns []*gorm.DB)
}

type AccessControlrSvc interface {
	UserHasPermissionToResource(userID int, resourceID int, resourceType tables.ResourceTypes, action accesscontrolsvc.CRUDOperation) (bool, error)
}

type AuthnSvc interface {
	ValidateToken(signedToken string) error
}

func Handler(db DB, accessControlSvc AccessControlrSvc, authnSvc AuthnSvc, logger logger.Logger) *AuthzSvcImpl {
	return &AuthzSvcImpl{DB: db, accessControlSvc: accessControlSvc, authnSvc: authnSvc, logger: logger}
}

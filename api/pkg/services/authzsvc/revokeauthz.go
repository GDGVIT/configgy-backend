package authzsvc

import (
	"context"
	"errors"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

func (svc *AuthzSvcImpl) AuthzDelete(ctx context.Context, AuthzPID string, userPID string) (api.GenericMessageResponse, int, error) {
	var err error
	permissionAssignment, err := svc.DB.GetPermissionAssignmentByPID(AuthzPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	var hasPermission bool
	switch permissionAssignment.ResourceType {
	case tables.VaultResource:
		hasPermission, err = svc.accessControlSvc.UserHasPermissionToResource(permissionAssignment.UserID, permissionAssignment.VaultID, permissionAssignment.ResourceType, accesscontrolsvc.DeleteOperation)
	case tables.GroupResource:
		hasPermission, err = svc.accessControlSvc.UserHasPermissionToResource(permissionAssignment.UserID, permissionAssignment.GroupID, permissionAssignment.ResourceType, accesscontrolsvc.DeleteOperation)
	case tables.CredentialResource:
		hasPermission, err = svc.accessControlSvc.UserHasPermissionToResource(permissionAssignment.UserID, permissionAssignment.CredentialID, permissionAssignment.ResourceType, accesscontrolsvc.DeleteOperation)
	default:
		return api.GenericMessageResponse{}, 0, err
	}
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	if !hasPermission {
		return api.GenericMessageResponse{}, 0, errors.New("user does not have permission to delete this resource")
	}
	tx := svc.DB.DeletePermissionAssignment(permissionAssignment)
	if tx.Error != nil {
		return api.GenericMessageResponse{}, 0, tx.Error
	}

	return api.GenericMessageResponse{}, 0, nil
}

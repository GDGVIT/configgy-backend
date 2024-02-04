package authzsvc

import (
	"context"
	"errors"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

func (svc *AuthzSvcImpl) AuthzEdit(ctx context.Context, AuthzUpdated api.AuthzPermission, AuthzPID string, UserPID string) (api.GenericMessageResponse, int, error) {
	permissionAssignment, err := svc.DB.GetPermissionAssignmentByPID(AuthzPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	var hasPermission bool
	switch permissionAssignment.ResourceType {
	case tables.VaultResource:
		hasPermission, err = svc.accessControlSvc.UserHasPermissionToResource(permissionAssignment.UserID, permissionAssignment.VaultID, permissionAssignment.ResourceType, accesscontrolsvc.UpdateOperation)
	case tables.GroupResource:
		hasPermission, err = svc.accessControlSvc.UserHasPermissionToResource(permissionAssignment.UserID, permissionAssignment.GroupID, permissionAssignment.ResourceType, accesscontrolsvc.UpdateOperation)
	case tables.CredentialResource:
		hasPermission, err = svc.accessControlSvc.UserHasPermissionToResource(permissionAssignment.UserID, permissionAssignment.CredentialID, permissionAssignment.ResourceType, accesscontrolsvc.UpdateOperation)
	default:
		return api.GenericMessageResponse{}, 0, err
	}
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	if !hasPermission {
		return api.GenericMessageResponse{}, 0, errors.New("user does not have permission to update this resource")
	}
	if AuthzUpdated.AccessLevel != "" {
		permissionAssignment.PermissionName = tables.Permission(AuthzUpdated.AccessLevel)
	}
	return api.GenericMessageResponse{}, 0, nil
}

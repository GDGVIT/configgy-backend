package authzsvc

import (
	"context"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

func (svc *AuthzSvcImpl) AuthzGet(ctx context.Context, AuthzPID string, UserPID string) (api.AuthzPermission, int, error) {
	user, err := svc.DB.GetUserByPID(UserPID)
	if err != nil {
		return api.AuthzPermission{}, 0, err
	}
	permissionAssignment, err := svc.DB.GetPermissionAssignmentByPID(AuthzPID)
	if err != nil {
		return api.AuthzPermission{}, 0, err
	}
	hasPermission := false
	switch permissionAssignment.ResourceType {
	case tables.VaultResource:
		hasPermission, err = svc.accessControlSvc.UserHasPermissionToResource(user.ID, permissionAssignment.VaultID, permissionAssignment.ResourceType, accesscontrolsvc.ReadOperation)
	case tables.GroupResource:
		hasPermission, err = svc.accessControlSvc.UserHasPermissionToResource(user.ID, permissionAssignment.GroupID, permissionAssignment.ResourceType, accesscontrolsvc.ReadOperation)
	case tables.CredentialResource:
		hasPermission, err = svc.accessControlSvc.UserHasPermissionToResource(user.ID, permissionAssignment.CredentialID, permissionAssignment.ResourceType, accesscontrolsvc.ReadOperation)
	}
	if err != nil {
		return api.AuthzPermission{}, 0, err
	}

	if !hasPermission {
		return api.AuthzPermission{}, 404, nil
	}
	return api.AuthzPermission{}, 0, nil
}

func (svc *AuthzSvcImpl) AuthzUserGet(ctx context.Context, UserPID string) (api.AuthzUserGetResponse, int, error) {
	response := api.AuthzUserGetResponse{}
	user, err := svc.DB.GetUserByPID(UserPID)
	if err != nil {
		return response, 0, err
	}

	permissionAssignments, err := svc.DB.GetPermissionsForUser(user.ID)
	if err != nil {
		return response, 0, err
	}

	for _, permissionAssignment := range permissionAssignments {
		identityType := api.AuthzPermissionIdentityType(permissionAssignment.IdentityType)
		permission := api.AuthzPermission{IdentityPid: &permissionAssignment.IdentityPID, IdentityType: &identityType, AccessLevel: api.AuthzPermissionAccessLevel(permissionAssignment.PermissionName), ResourcePid: &permissionAssignment.ResourcePID, ResourceType: api.AuthzPermissionResourceType(permissionAssignment.ResourceType)}
		response = append(response, permission)
	}

	return response, 0, nil
}

func (svc *AuthzSvcImpl) AuthzGroupGet(ctx context.Context, GroupPID string, UserPID string) (api.AuthzGroupGetResponse, int, error) {
	response := api.AuthzGroupGetResponse{}
	user, err := svc.DB.GetUserByPID(UserPID)
	if err != nil {
		return response, 0, err
	}
	group, err := svc.DB.GetGroupByPID(GroupPID)
	if err != nil {
		return response, 0, err
	}

	ok, err := svc.DB.IsUserInGroup(user.ID, group.ID)
	if err != nil {
		return response, 0, err
	}
	if !ok {
		return response, 0, err
	}

	permissionAssignments, err := svc.DB.GetPermissionsForGroup(group.ID)
	if err != nil {
		return response, 0, err
	}

	for _, permissionAssignment := range permissionAssignments {
		identityType := api.AuthzPermissionIdentityType(permissionAssignment.IdentityType)
		permission := api.AuthzPermission{IdentityPid: &permissionAssignment.IdentityPID, IdentityType: &identityType, AccessLevel: api.AuthzPermissionAccessLevel(permissionAssignment.PermissionName), ResourcePid: &permissionAssignment.ResourcePID, ResourceType: api.AuthzPermissionResourceType(permissionAssignment.ResourceType)}
		response = append(response, permission)
	}

	return response, 0, nil
}

func (svc *AuthzSvcImpl) AuthzVaultGet(ctx context.Context, VaultPID string, UserPID string) (api.AuthzVaultGetResponse, int, error) {
	response := api.AuthzVaultGetResponse{}
	user, err := svc.DB.GetUserByPID(UserPID)
	if err != nil {
		return response, 0, err
	}
	vault, err := svc.DB.GetVaultByPID(VaultPID)
	if err != nil {
		return response, 0, err
	}
	permission, err := svc.accessControlSvc.UserHasPermissionToResource(user.ID, vault.ID, tables.VaultResource, accesscontrolsvc.UpdateOperation)
	if err != nil {
		return response, 0, err
	}
	if !permission {
		return response, 0, err
	}
	permissions, err := svc.DB.GetPermissionsForResource(vault.ID, tables.VaultResource)
	if err != nil {
		return response, 0, nil
	}
	for _, permission := range permissions {
		authzpermission := api.AuthzPermission{
			AccessLevel:  api.AuthzPermissionAccessLevel(permission.PermissionName),
			IdentityPid:  &permission.IdentityPID,
			IdentityType: (*api.AuthzPermissionIdentityType)(&permission.IdentityType),
			ResourcePid:  &VaultPID,
			ResourceType: api.AuthzPermissionResourceTypeVault,
		}
		response = append(response, authzpermission)
	}

	return response, 0, nil

}

func (svc *AuthzSvcImpl) AuthzCredentialGet(ctx context.Context, CredentialPID string, UserPID string) (api.AuthzCredentialGetResponse, int, error) {
	response := api.AuthzCredentialGetResponse{}
	user, err := svc.DB.GetUserByPID(UserPID)
	if err != nil {
		return response, 0, err
	}
	credential, err := svc.DB.GetCredentialByPID(CredentialPID)
	if err != nil {
		return response, 0, err
	}
	permission, err := svc.accessControlSvc.UserHasPermissionToResource(user.ID, credential.CredentialID, tables.CredentialResource, accesscontrolsvc.UpdateOperation)
	if err != nil {
		return response, 0, err
	}
	if !permission {
		return response, 0, nil
	}
	permissions, err := svc.DB.GetPermissionsForResource(credential.CredentialID, tables.CredentialResource)
	if err != nil {
		return response, 0, nil
	}
	for _, permission := range permissions {
		authzpermission := api.AuthzPermission{
			AccessLevel:  api.AuthzPermissionAccessLevel(permission.PermissionName),
			IdentityPid:  &permission.IdentityPID,
			IdentityType: (*api.AuthzPermissionIdentityType)(&permission.IdentityType),
			ResourcePid:  &CredentialPID,
			ResourceType: api.AuthzPermissionResourceTypeCredential,
		}
		response = append(response, authzpermission)
	}
	return response, 0, nil
}

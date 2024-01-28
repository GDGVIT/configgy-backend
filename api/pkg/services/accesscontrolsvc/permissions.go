package accesscontrolsvc

import (
	"errors"
	"strings"

	"github.com/GDGVIT/configgy-backend/constants"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
	"gorm.io/gorm"
)

func (svc *accessControlSvcImpl) CreatePermissionAssignment(permission_assignment *tables.PermissionAssignments) *gorm.DB {
	return svc.DB.CreatePermissionAssignment(permission_assignment)
}

func (svc *accessControlSvcImpl) GetPermissionsForUser(userID int) ([]*tables.PermissionAssignments, error) {
	return svc.DB.GetPermissionsForUser(userID)
}

func (svc *accessControlSvcImpl) UserHasPermissionToResource(userID int, resourceID int, resourceType tables.ResourceTypes, action CRUDOperation) (bool, error) {
	// permission order (of importance)
	// 1. Direct permission on resource
	// 2. Direct Group permission on resource
	// 3. User Vault permission on resource
	// 4. Group Vault permission on resource
	var userPermissions []*tables.PermissionAssignments
	// direct permission on resource
	resourcePermissionForUser, err := svc.DB.GetPermissionsForUserOnResource(userID, resourceID, resourceType)
	if err != nil {
		return false, err
	}
	userPermissions = append(userPermissions, resourcePermissionForUser)
	userGroups, err := svc.DB.GetGroupsByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, group := range userGroups {
		// direct group permission on resource
		groupPermissionForUser, err := svc.DB.GetPermissionsForGroupOnResource(group.GroupID, resourceID, resourceType)
		if err != nil {
			return false, err
		}
		if groupPermissionForUser != nil {
			userPermissions = append(userPermissions, groupPermissionForUser)
		}
	}

	// handle credential-vault relations
	if resourceType == tables.CredentialResource {
		// user vault permission on resource
		vaultCredential, err := svc.DB.GetVaultIDForCredential(resourceID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		// credential is not in a vault so no need to check permissions on the vault
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			vaultPermissions, err := svc.DB.GetPermissionsForResource(vaultCredential.VaultID, tables.VaultResource)
			if err != nil {
				return false, err
			}
			userPermissions = append(userPermissions, vaultPermissions...)
			for _, group := range userGroups {
				// group vault permission on resource
				groupPermissions, err := svc.DB.GetPermissionsForGroupOnResource(group.GroupID, vaultCredential.VaultID, tables.VaultResource)
				if err != nil {
					return false, err
				}
				userPermissions = append(userPermissions, groupPermissions)
			}
		}
	}
	// check if user has permission
	for _, permission := range userPermissions {
		switch {
		case action == CreateOperation:
			// does the user have permissions to create a resource
			if strings.Contains(StringifyListOfPermissions(constants.CreatePermissions), string(permission.PermissionName)) {
				return true, nil
			}
			continue
		case action == ReadOperation:
			// does the user have permissions to read a resource
			if strings.Contains(StringifyListOfPermissions(constants.ReadPermissions), string(permission.PermissionName)) {
				return true, nil
			}
			continue
		case action == UpdateOperation:
			// does the user have permissions to update a resource
			if strings.Contains(StringifyListOfPermissions(constants.UpdatePermissions), string(permission.PermissionName)) {
				return true, nil
			}
			continue
		case action == DeleteOperation:
			// does the user have permissions to delete a resource
			if strings.Contains(StringifyListOfPermissions(constants.DeletePermissions), string(permission.PermissionName)) {
				return true, nil
			}
			continue
		}
	}
	return false, nil

}

func (svc *accessControlSvcImpl) GroupHasPermissionToResource(groupID int, resourceID int, resourceType tables.ResourceTypes, action CRUDOperation) (bool, error) {
	// group direct permission for resource
	resourcePermissionForGroup, err := svc.DB.GetPermissionsForGroupOnResource(groupID, resourceID, resourceType)
	if err != nil {
		return false, err
	}
	switch {
	case action == CreateOperation:
		// does the group have permissions to create a resource
		return strings.Contains(StringifyListOfPermissions(constants.CreatePermissions), string(resourcePermissionForGroup.PermissionName)), nil
	case action == ReadOperation:
		// does the group have permissions to read a resource
		return strings.Contains(StringifyListOfPermissions(constants.ReadPermissions), string(resourcePermissionForGroup.PermissionName)), nil
	case action == UpdateOperation:
		// does the group have permissions to update a resource
		return strings.Contains(StringifyListOfPermissions(constants.UpdatePermissions), string(resourcePermissionForGroup.PermissionName)), nil
	case action == DeleteOperation:
		// does the group have permissions to delete a resource
		return strings.Contains(StringifyListOfPermissions(constants.DeletePermissions), string(resourcePermissionForGroup.PermissionName)), nil
	default:
		return false, nil
	}
}

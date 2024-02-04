package vaultsvc

import (
	"context"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

func (svc *VaultServiceImpl) VaultDelete(c context.Context, vaultPID string, userPID string) (api.GenericMessageResponse, int, error) {
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	vault, err := svc.DB.GetVaultByPID(vaultPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	hasPermission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, vault.ID, tables.VaultResource, accesscontrolsvc.DeleteOperation)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	if !hasPermission {
		message := "User does not have permission to delete this vault"
		return api.GenericMessageResponse{Message: &message}, 0, nil
	}
	// delete the vault
	err = svc.DB.DeleteVault(vault.ID)
	// delete the permission assignment
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	return api.GenericMessageResponse{}, 0, nil
}

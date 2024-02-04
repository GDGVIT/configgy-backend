package vaultsvc

import (
	"context"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

func (svc *VaultServiceImpl) VaultGet(c context.Context, vaultPID string, userPID string) (api.VaultCreateRequest, int, error) {
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.VaultCreateRequest{}, 0, err
	}
	vault, err := svc.DB.GetVaultByPID(vaultPID)
	if err != nil {
		return api.VaultCreateRequest{}, 0, err
	}
	permission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, vault.ID, tables.VaultResource, accesscontrolsvc.ReadOperation)
	if err != nil {
		return api.VaultCreateRequest{}, 0, err
	}
	if !permission {
		return api.VaultCreateRequest{}, 0, nil
	}
	return api.VaultCreateRequest{
		Name:        vault.Name,
		Description: vault.Description,
		PublicKey:   string(vault.PublicKey),
	}, 0, nil
}

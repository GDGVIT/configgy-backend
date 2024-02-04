package vaultsvc

import (
	"context"
	"net/http"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

func (svc *VaultServiceImpl) VaultCreate(c context.Context, req api.VaultCreateRequest, userPID string) (api.GenericMessageResponse, int, error) {
	vault := tables.Vault{
		Name:        req.Name,
		PID:         tables.UUIDWithPrefix("vault"),
		Description: req.Description,
		PublicKey:   []byte(req.PublicKey),
		IsPersonal:  true,
	}
	err := svc.DB.CreateVault(vault, userPID)
	if err != nil {
		return api.GenericMessageResponse{}, http.StatusInternalServerError, err
	}
	message := "Vault created successfully"
	return api.GenericMessageResponse{Message: &message}, http.StatusOK, nil
}

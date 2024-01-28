package credentialsvc

import (
	"context"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

func (svc *CredentialServiceImpl) CredentialDelete(c context.Context, credentialPID string, userPID string) (api.GenericMessageResponse, int, error) {
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{Success: false}, 0, err
	}
	credential, err := svc.DB.GetCredentialByPID(credentialPID)
	if err != nil {
		return api.GenericMessageResponse{Success: false}, 0, err
	}

	hasPermission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, credential.ID, tables.CredentialResource, accesscontrolsvc.DeleteOperation)
	if err != nil {
		return api.GenericMessageResponse{Success: false}, 0, err
	}
	if !hasPermission {
		return api.GenericMessageResponse{Success: false}, 403, nil
	}
	err = svc.DB.DeleteCredentialByID(credential.ID, *credential)
	if err != nil {
		return api.GenericMessageResponse{Success: false}, 0, err
	}
	return api.GenericMessageResponse{Success: true}, 0, nil
}

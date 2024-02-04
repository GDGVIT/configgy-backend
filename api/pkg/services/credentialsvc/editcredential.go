package credentialsvc

import (
	"context"
	"errors"
	"os"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/crypto"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
	"gorm.io/gorm"
)

// Edit a credential
// (POST /credentials/edit/{pid})
func (svc *CredentialServiceImpl) CredentialEdit(ctx context.Context, credentialPID string, req api.CredentialCreateRequest, userPID string) (api.GenericMessageResponse, int, error) {
	credential, err := svc.DB.GetCredentialByPID(credentialPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	// check if user has permission to edit credential
	hasPermission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, credential.ID, tables.CredentialResource, accesscontrolsvc.UpdateOperation)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	if !hasPermission {
		message := "User does not have permission to edit this credential"
		return api.GenericMessageResponse{Message: &message}, 0, nil
	}
	if req.Name != "" {
		credential.CredentialName = req.Name
	}
	if req.Notes != nil {
		credential.Notes = *req.Notes
	}
	tx := svc.DB.UpdateCredentialByID(credential.ID, credential)
	if tx.Error != nil {
		return api.GenericMessageResponse{}, 0, tx.Error
	}
	return api.GenericMessageResponse{}, 0, nil
}

// (POST /credentials/featureflag/add/{pid})
func (svc *CredentialServiceImpl) CredentialFeatureFlagAdd(ctx context.Context, req api.CredentialFeatureFlagAddJSONBody, credentialPID string, userPID string) (api.GenericMessageResponse, int, error) {
	credential, err := svc.DB.GetCredentialByPID(credentialPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}

	// check if user has permission to edit credential
	hasPermission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, credential.ID, tables.CredentialResource, accesscontrolsvc.UpdateOperation)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	if !hasPermission {
		message := "User does not have permission to edit this credential"
		return api.GenericMessageResponse{Message: &message}, 0, nil
	}
	txns := []*gorm.DB{}
	for _, featureFlag := range req.FeatureFlags {
		tx, err := svc.DB.AddFeatureFlag(&tables.FeatureFlagData{
			FeatureFlagCredentialID: credential.ID,
			Name:                    featureFlag.Name,
			FlagKey:                 featureFlag.Key,
			FlagValue:               featureFlag.Value,
		})
		if err != nil {
			svc.DB.RollbackTxns(txns)
			return api.GenericMessageResponse{}, 0, tx.Error
		}
		txns = append(txns, tx)
	}
	return api.GenericMessageResponse{}, 0, nil
}

// Remove Feature Flag K:V Pairs from existing Feature Flag Credential
// (POST /credentials/featureflag/remove/{pid})
func (svc *CredentialServiceImpl) CredentialFeatureFlagRemove(ctx context.Context, credentialPID string, req api.CredentialFeatureFlagRemoveJSONBody, userPID string) (api.GenericMessageResponse, int, error) {
	credential, err := svc.DB.GetCredentialByPID(credentialPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	// check if user has permission to edit credential
	hasPermission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, credential.ID, tables.CredentialResource, accesscontrolsvc.DeleteOperation)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	if !hasPermission {
		message := "User does not have permission to edit this credential"
		return api.GenericMessageResponse{Message: &message}, 0, nil
	}
	txns := []*gorm.DB{}
	for _, featureFlag := range req.FeatureFlags {
		tx := svc.DB.DeleteFeatureFlagValueForKey(credential.ID, featureFlag)
		if tx.Error != nil {
			svc.DB.RollbackTxns(txns)
			return api.GenericMessageResponse{}, 0, tx.Error
		}
		txns = append(txns, tx)
	}
	return api.GenericMessageResponse{}, 0, nil
}

// Get Feature Flag value for key
// (POST /credentials/featureflag/{pid}/{key})
func (svc *CredentialServiceImpl) CredentialFeatureFlagGet(ctx context.Context, credentialPID string, key string, userPID string) (api.FeatureFlag, int, error) {
	credential, err := svc.DB.GetCredentialByPID(credentialPID)
	if err != nil {
		return api.FeatureFlag{}, 0, err
	}
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.FeatureFlag{}, 0, err
	}
	// check if user has permission to read credential
	hasPermission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, credential.ID, tables.CredentialResource, accesscontrolsvc.ReadOperation)
	if err != nil {
		return api.FeatureFlag{}, 0, err
	}
	if !hasPermission {
		return api.FeatureFlag{}, 0, errors.New("user does not have permission to read this credential")
	}
	featureFlag, err := svc.DB.GetFeatureFlagValueForKey(credential.ID, key)
	if err != nil {
		return api.FeatureFlag{}, 0, err
	}
	return api.FeatureFlag{Key: featureFlag.FlagKey, Value: featureFlag.FlagValue, Name: featureFlag.Name}, 0, nil
}

// Edit a file credential
// (POST /credentials/file/edit/{pid})
func (svc *CredentialServiceImpl) CredentialFileEdit(ctx context.Context, req api.CredentialFileEditJSONRequestBody, credentialPID string, userPID string) (api.GenericMessageResponse, int, error) {
	credential, err := svc.DB.GetCredentialByPID(credentialPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	credentialData, err := svc.DB.GetFileCredentialByID(credential.CredentialID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	// check if user has permission to edit credential
	hasPermission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, credential.ID, tables.CredentialResource, accesscontrolsvc.UpdateOperation)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	if !hasPermission {
		message := "User does not have permission to edit this credential"
		return api.GenericMessageResponse{Message: &message}, 0, nil
	}
	if req.File != nil {
		// read old data
		oldData, err := os.ReadFile(credentialData.FilePath)
		if err != nil {
			return api.GenericMessageResponse{}, 0, err
		}
		// clear old file
		err = os.WriteFile(credentialData.FilePath, []byte{}, 0644)
		if err != nil {
			// rollback
			os.WriteFile(credentialData.FilePath, oldData, 0644)
			return api.GenericMessageResponse{}, 0, err
		}
		encryptedData, err := crypto.EncryptBytes(req.File, svc.secretKey)
		if err != nil {
			// rollback
			os.WriteFile(credentialData.FilePath, oldData, 0644)
			return api.GenericMessageResponse{}, 0, err
		}
		err = os.WriteFile(credentialData.FilePath, encryptedData, 0644)
		if err != nil {
			// rollback
			os.WriteFile(credentialData.FilePath, oldData, 0644)
			return api.GenericMessageResponse{}, 0, err
		}
	}
	return api.GenericMessageResponse{}, 0, nil
}

// Edit a password credential
// (POST /credentials/password/edit/{pid})
func (svc *CredentialServiceImpl) CredentialPasswordEdit(ctx context.Context, req api.CredentialPasswordEditJSONRequestBody, credentialPID string, userPID string) (api.GenericMessageResponse, int, error) {
	credential, err := svc.DB.GetCredentialByPID(credentialPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	credentialData, err := svc.DB.GetPasswordCredentialByID(credential.CredentialID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	// check if user has permission to edit credential
	hasPermission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, credential.ID, tables.CredentialResource, accesscontrolsvc.UpdateOperation)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	if !hasPermission {
		message := "User does not have permission to edit this credential"
		return api.GenericMessageResponse{Message: &message}, 0, nil
	}
	if req.Password != nil {
		// read old data
		encryptedPassword, err := crypto.EncryptBytes(req.Password, svc.secretKey)
		if err != nil {
			return api.GenericMessageResponse{}, 0, err
		}
		credentialData.Password = encryptedPassword
	}
	if req.Username != "" {
		credentialData.Username = req.Username
	}

	if req.PasswordStrength != nil {
		credentialData.PasswordStrength = *req.PasswordStrength
	}
	if req.Totp != nil {
		encryptedTotpSecret, err := crypto.EncryptBytes(req.Totp.TotpSecret, svc.secretKey)
		if err != nil {
			return api.GenericMessageResponse{}, 0, err
		}
		credentialData.TOTPKey = encryptedTotpSecret
		credentialData.TOTPLength = req.Totp.TotpLength
		credentialData.TOTPPeriod = req.Totp.TotpPeriod
	}
	if req.ExpiresAt != nil {
		credentialData.ExpiresAt = *req.ExpiresAt
	}
	tx := svc.DB.UpdatePasswordCredentialByID(credential.CredentialID, credentialData)
	if tx.Error != nil {
		return api.GenericMessageResponse{}, 0, tx.Error
	}
	return api.GenericMessageResponse{}, 0, nil
}

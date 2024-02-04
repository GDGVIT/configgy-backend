package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CredentialService interface {
	CredentialCreate(c context.Context, req CredentialCreateRequest, userPID string) (GenericMessageResponse, int, error)
	CredentialDelete(c context.Context, credentialPID string, userPID string) (GenericMessageResponse, int, error)
	CredentialGet(c context.Context, credentialPID string, userPID string) (CredentialCreateRequest, int, error)
	CredentialEdit(c context.Context, credentialPID string, req CredentialCreateRequest, userPID string) (GenericMessageResponse, int, error)
	// Add Feature Flag K:V Pairs to existing Feature Flag Credential
	// (POST /credentials/featureflag/add/{pid})
	CredentialFeatureFlagAdd(ctx context.Context, req CredentialFeatureFlagAddJSONBody, credentialPID string, userPID string) (GenericMessageResponse, int, error)
	// Remove Feature Flag K:V Pairs from existing Feature Flag Credential
	// (POST /credentials/featureflag/remove/{pid})
	CredentialFeatureFlagRemove(ctx context.Context, credentialPID string, req CredentialFeatureFlagRemoveJSONBody, userPID string) (GenericMessageResponse, int, error)
	// Get Feature Flag value for key
	// (POST /credentials/featureflag/{pid}/{key})
	CredentialFeatureFlagGet(ctx context.Context, credentialPID string, key string, userPID string) (FeatureFlag, int, error)
	// Edit a file credential
	// (POST /credentials/file/edit/{pid})
	CredentialFileEdit(ctx context.Context, req CredentialFileEditJSONRequestBody, credentialPID string, userPID string) (GenericMessageResponse, int, error)
	// Edit a password credential
	// (POST /credentials/password/edit/{pid})
	CredentialPasswordEdit(ctx context.Context, req CredentialPasswordEditJSONRequestBody, credentialPID string, userPID string) (GenericMessageResponse, int, error)
}

// CredentialCreate - Handler for credential creation
// (POST "/credential/create")
func (svc *Service) CredentialCreate(c echo.Context) error {
	userPID, err := getUserPIDFromAuthToken(c)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	// Parse the request body into the SignupRequest struct
	request := &CredentialCreateRequest{}
	if err := c.Bind(request); err != nil {
		svc.logger.Error("Failed to parse request body:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}

	// You can now perform your credential creation logic here
	response, httpCode, err := svc.Services.CredentialSvc.CredentialCreate(c.Request().Context(), *request, userPID)
	if err != nil {
		svc.logger.Error("Failed to create credential:", err)
		return echo.NewHTTPError(httpCode, "Failed to create credential")
	}
	svc.logger.Info("Credential creation request received")

	// Return a response (e.g., a success message)
	return c.JSON(http.StatusOK, response)
}

// CredentialDelete - Handler for credential deletion
// (DELETE "/credential/delete/{id}")
func (svc *Service) CredentialDelete(c echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(c)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	response, httpCode, err := svc.Services.CredentialSvc.CredentialDelete(c.Request().Context(), pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to delete credential:", err)
		return echo.NewHTTPError(httpCode, "Failed to create credential")
	}
	svc.logger.Info("Credential delete request received")

	// Return a response (e.g., a success message)
	return c.JSON(http.StatusOK, response)
}

// CredentialGet - Handler for credential getting
// (GET "/credential/get/{id}")
func (svc *Service) CredentialGet(c echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(c)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	// You can now perform your signup logic here
	response, httpCode, err := svc.Services.CredentialSvc.CredentialGet(c.Request().Context(), pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to get credential:", err)
		return echo.NewHTTPError(httpCode, "Failed to get credential")
	}
	svc.logger.Info("Credential get request received")

	// Return a response (e.g., a success message)
	return c.JSON(http.StatusOK, response)
}

// CredentialEdit
// (POST "/credential/edit/{id}")
func (svc *Service) CredentialEdit(c echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(c)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	// Parse the request body into the SignupRequest struct
	request := &CredentialCreateRequest{}
	if err := c.Bind(request); err != nil {
		svc.logger.Error("Failed to parse request body:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}
	response, status, err := svc.Services.CredentialSvc.CredentialEdit(c.Request().Context(), pid, *request, userPID)
	if err != nil {
		svc.logger.Error("Failed to edit credential:", err)
		return echo.NewHTTPError(status, "Failed to edit credential")
	}
	svc.logger.Info("Credential edit request received")
	return c.JSON(http.StatusOK, response)
}

func (svc *Service) CredentialPasswordEdit(c echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(c)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	request := &CredentialPasswordEditJSONRequestBody{}
	if err := c.Bind(request); err != nil {
		svc.logger.Error("Failed to parse request body:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}
	response, status, err := svc.Services.CredentialSvc.CredentialPasswordEdit(c.Request().Context(), *request, pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to edit credential:", err)
		return echo.NewHTTPError(status, "Failed to edit credential")
	}
	svc.logger.Info("Credential edit request received")
	return c.JSON(http.StatusOK, response)
}

func (svc *Service) CredentialFileEdit(c echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(c)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	request := &CredentialFileEditJSONRequestBody{}
	if err := c.Bind(request); err != nil {
		svc.logger.Error("Failed to parse request body:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}
	response, status, err := svc.Services.CredentialSvc.CredentialFileEdit(c.Request().Context(), *request, pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to edit credential:", err)
		return echo.NewHTTPError(status, "Failed to edit credential")
	}
	svc.logger.Info("Credential edit request received")
	return c.JSON(http.StatusOK, response)
}

func (svc *Service) CredentialFeatureFlagGet(c echo.Context, pid string, featureFlagKey string) error {
	userPID, err := getUserPIDFromAuthToken(c)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	response, status, err := svc.Services.CredentialSvc.CredentialFeatureFlagGet(c.Request().Context(), pid, featureFlagKey, userPID)
	if err != nil {
		svc.logger.Error("Failed to get credential:", err)
		return echo.NewHTTPError(status, "Failed to get credential")
	}
	svc.logger.Info("Credential get request received")
	return c.JSON(http.StatusOK, response)
}

func (svc *Service) CredentialFeatureFlagAdd(c echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	request := &CredentialFeatureFlagAddJSONBody{}
	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}
	response, status, err := svc.Services.CredentialSvc.CredentialFeatureFlagAdd(c.Request().Context(), *request, pid, userPID)
	if err != nil {
		return echo.NewHTTPError(status, "Failed to add feature flag")
	}
	return c.JSON(http.StatusOK, response)
}

func (svc *Service) CredentialFeatureFlagRemove(c echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	request := &CredentialFeatureFlagRemoveJSONBody{}
	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}
	response, status, err := svc.Services.CredentialSvc.CredentialFeatureFlagRemove(c.Request().Context(), pid, *request, userPID)
	if err != nil {
		return echo.NewHTTPError(status, "Failed to remove feature flag")
	}
	return c.JSON(http.StatusOK, response)
}

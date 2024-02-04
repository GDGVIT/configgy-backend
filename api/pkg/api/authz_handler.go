package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthzService interface {
	// GET /authz/permission/{pid})
	AuthzGet(ctx context.Context, authzPID string, userPID string) (AuthzPermission, int, error)

	// (GET /authz/credential/{pid})
	AuthzCredentialGet(ctx context.Context, credentialPID string, userPID string) (AuthzCredentialGetResponse, int, error)

	// (POST /authz/grant)
	AuthzCreate(ctx context.Context, req AuthzCreateRequest, userPID string) (GenericMessageResponse, int, error)

	// (GET /authz/group/{pid})
	AuthzGroupGet(ctx context.Context, groupPID string, userPID string) (AuthzGroupGetResponse, int, error)

	// (POST /authz/modify/{pid})
	AuthzEdit(ctx context.Context, req AuthzPermission, authzPID string, userPID string) (GenericMessageResponse, int, error)

	// (DELETE /authz/revoke/{pid})
	AuthzDelete(ctx context.Context, authzPID string, userPID string) (GenericMessageResponse, int, error)

	// (GET /authz/user/{pid})
	AuthzUserGet(ctx context.Context, userPID string) (AuthzUserGetResponse, int, error)

	// (GET /authz/vault/{pid})
	AuthzVaultGet(ctx context.Context, vaultPID string, userPID string) (AuthzVaultGetResponse, int, error)
}

// (GET /authz/permission/{pid})
func (svc *Service) AuthzGet(ctx echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	response, httpCode, err := svc.Services.AuthzSvc.AuthzGet(svc.ctx, pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to get authz:", err)
		return echo.NewHTTPError(httpCode, "Failed to get authz")
	}
	svc.logger.Info("Authz get request received")
	return ctx.JSON(http.StatusOK, response)
}

// (GET /authz/credential/{pid})
func (svc *Service) AuthzCredentialGet(ctx echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	response, httpCode, err := svc.Services.AuthzSvc.AuthzCredentialGet(svc.ctx, pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to get authz credential:", err)
		return echo.NewHTTPError(httpCode, "Failed to get authz credential")
	}
	svc.logger.Info("Authz credential get request received")
	return ctx.JSON(http.StatusOK, response)
}

// (POST /authz/grant)
func (svc *Service) AuthzCreate(ctx echo.Context) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}

	// Parse the request body into the AuthzCreateRequest struct
	request := &AuthzCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		svc.logger.Error("Failed to parse request body:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}

	// You can now perform your authz creation logic here
	response, httpCode, err := svc.Services.AuthzSvc.AuthzCreate(svc.ctx, *request, userPID)
	if err != nil {
		svc.logger.Error("Failed to create authz:", err)
		return echo.NewHTTPError(httpCode, "Failed to create authz")
	}
	svc.logger.Info("Authz create request received")

	// Return a response (e.g., a success message)
	return ctx.JSON(http.StatusOK, response)
}

// (GET /authz/group/{pid})
func (svc *Service) AuthzGroupGet(ctx echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	response, httpCode, err := svc.Services.AuthzSvc.AuthzGroupGet(svc.ctx, pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to get authz group:", err)
		return echo.NewHTTPError(httpCode, "Failed to get authz group")
	}
	svc.logger.Info("Authz group get request received")
	return ctx.JSON(http.StatusOK, response)
}

// (POST /authz/modify/{pid})
func (svc *Service) AuthzEdit(ctx echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}

	// Parse the request body into the AuthzPermission struct
	request := &AuthzPermission{}
	if err := ctx.Bind(request); err != nil {
		svc.logger.Error("Failed to parse request body:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}

	// You can now perform your authz modification logic here
	response, httpCode, err := svc.Services.AuthzSvc.AuthzEdit(svc.ctx, *request, pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to modify authz:", err)
		return echo.NewHTTPError(httpCode, "Failed to modify authz")
	}
	svc.logger.Info("Authz modify request received")

	// Return a response (e.g., a success message)
	return ctx.JSON(http.StatusOK, response)
}

// (DELETE /authz/revoke/{pid})
func (svc *Service) AuthzDelete(ctx echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}

	response, httpCode, err := svc.Services.AuthzSvc.AuthzDelete(svc.ctx, pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to delete authz:", err)
		return echo.NewHTTPError(httpCode, "Failed to delete authz")
	}
	svc.logger.Info("Authz delete request received")

	// Return a response (e.g., a success message)
	return ctx.JSON(http.StatusOK, response)
}

// (GET /authz/user/{pid})
func (svc *Service) AuthzUserGet(ctx echo.Context) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	response, httpCode, err := svc.Services.AuthzSvc.AuthzUserGet(svc.ctx, userPID)
	if err != nil {
		svc.logger.Error("Failed to get authz user:", err)
		return echo.NewHTTPError(httpCode, "Failed to get authz user")
	}
	svc.logger.Info("Authz user get request received")
	return ctx.JSON(http.StatusOK, response)
}

// (GET /authz/vault/{pid})
func (svc *Service) AuthzVaultGet(ctx echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	response, httpCode, err := svc.Services.AuthzSvc.AuthzVaultGet(svc.ctx, pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to get authz vault:", err)
		return echo.NewHTTPError(httpCode, "Failed to get authz vault")
	}
	svc.logger.Info("Authz vault get request received")
	return ctx.JSON(http.StatusOK, response)
}

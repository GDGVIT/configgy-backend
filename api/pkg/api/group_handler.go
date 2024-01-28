package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GroupService interface {
	// (GET /groups)
	GroupsGet(ctx context.Context, userPID string) ([]Group, int, error)

	// (POST /groups/add/{pid})
	GroupAdd(ctx context.Context, groupPID string, req GroupAddRequest, userPID string) (GenericMessageResponse, int, error)

	// (POST /groups/create)
	GroupCreate(ctx context.Context, req GroupCreateRequest, userPID string) (GenericMessageResponse, int, error)

	// (DELETE /groups/delete/{pid})
	GroupDelete(ctx context.Context, groupPID string, userPID string) (GenericMessageResponse, int, error)

	// (POST /groups/edit/{pid})
	GroupEdit(ctx context.Context, groupPID string, req GroupEditRequest, userPID string) (GenericMessageResponse, int, error)

	// (POST /groups/remove/{pid})
	GroupRemove(ctx context.Context, groupPID string, req GroupRemoveRequest, userPID string) (GenericMessageResponse, int, error)

	// (GET /groups/{pid})
	GroupGet(ctx context.Context, groupPID string, userPID string) (Group, int, error)
}

// (GET /groups)
func (svc *Service) GroupsGet(ctx echo.Context) error {
	return nil
}

// (POST /groups/add/{pid})
func (svc *Service) GroupAdd(ctx echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	// Parse the request body into the GroupAddRequest struct
	request := &GroupAddRequest{}
	if err := ctx.Bind(request); err != nil {
		svc.logger.Error("Failed to parse request body:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}
	response, httpCode, err := svc.Services.GroupSvc.GroupAdd(svc.ctx, pid, *request, userPID)
	if err != nil {
		svc.logger.Error("Failed to add group:", err)
		return echo.NewHTTPError(httpCode, "Failed to add group")
	}
	svc.logger.Info("Group add request received")
	return ctx.JSON(http.StatusOK, response)
}

// (POST /groups/create)
func (svc *Service) GroupCreate(ctx echo.Context) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	// Parse the request body into the GroupCreateRequest struct
	request := &GroupCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		svc.logger.Error("Failed to parse request body:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}
	response, httpCode, err := svc.Services.GroupSvc.GroupCreate(svc.ctx, *request, userPID)
	if err != nil {
		svc.logger.Error("Failed to create group:", err)
		return echo.NewHTTPError(httpCode, "Failed to create group")
	}
	svc.logger.Info("Group creation request received")
	return ctx.JSON(http.StatusOK, response)
}

// (DELETE /groups/delete/{pid})
func (svc *Service) GroupDelete(ctx echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	response, httpCode, err := svc.Services.GroupSvc.GroupDelete(svc.ctx, pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to delete group:", err)
		return echo.NewHTTPError(httpCode, "Failed to delete group")
	}
	svc.logger.Info("Group delete request received")
	return ctx.JSON(http.StatusOK, response)
}

// (POST /groups/edit/{pid})
func (svc *Service) GroupEdit(ctx echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}

	// Parse the request body into the GroupEditRequest struct
	request := &GroupEditRequest{}
	if err := ctx.Bind(request); err != nil {
		svc.logger.Error("Failed to parse request body:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}

	// You can now perform your group modification logic here
	response, httpCode, err := svc.Services.GroupSvc.GroupEdit(svc.ctx, pid, *request, userPID)
	if err != nil {
		svc.logger.Error("Failed to modify group:", err)
		return echo.NewHTTPError(httpCode, "Failed to modify group")
	}
	svc.logger.Info("Group modify request received")

	// Return a response (e.g., a success message)
	return ctx.JSON(http.StatusOK, response)
}

// (POST /groups/remove/{pid})
func (svc *Service) GroupRemove(ctx echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	// Parse the request body into the GroupRemoveRequest struct
	request := &GroupRemoveRequest{}
	if err := ctx.Bind(request); err != nil {
		svc.logger.Error("Failed to parse request body:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}
	response, httpCode, err := svc.Services.GroupSvc.GroupRemove(svc.ctx, pid, *request, userPID)
	if err != nil {
		svc.logger.Error("Failed to remove group:", err)
		return echo.NewHTTPError(httpCode, "Failed to remove group")
	}
	svc.logger.Info("Group remove request received")
	return ctx.JSON(http.StatusOK, response)
}

// (GET /groups/{pid})
func (svc *Service) GroupGet(ctx echo.Context, pid string) error {
	userPID, err := getUserPIDFromAuthToken(ctx)
	if err != nil {
		svc.logger.Error("Failed to get user PID from auth token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}
	response, httpCode, err := svc.Services.GroupSvc.GroupGet(svc.ctx, pid, userPID)
	if err != nil {
		svc.logger.Error("Failed to get group:", err)
		return echo.NewHTTPError(httpCode, "Failed to get group")
	}
	svc.logger.Info("Group get request received")
	return ctx.JSON(http.StatusOK, response)
}

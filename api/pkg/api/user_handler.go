package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	SignUp(c context.Context, req SignupRequest) (GenericMessageResponse, int, error)
	Login(c context.Context, req LoginRequest) (LoginResponse, int, error)
	VerifyUser(c context.Context, req VerifyUserParams) (GenericMessageResponse, int, error)
}

// SignUp - Signup
// (POST /user/signup)
func (svc *Service) Signup(c echo.Context) error {

	// Parse the request body into the SignupRequest struct
	request := &SignupRequest{}
	if err := c.Bind(request); err != nil {
		svc.logger.Error("Failed to parse request body:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}

	// You can now perform your signup logic here
	response, httpCode, err := svc.Services.UserSvc.SignUp(c.Request().Context(), *request)
	if err != nil {
		svc.logger.Error("Failed to signup:", err)
		return echo.NewHTTPError(httpCode, "Failed to signup")
	}
	svc.logger.Info("Signup request received")

	// Return a response (e.g., a success message)
	return c.JSON(http.StatusOK, response)
}

// Login - Login
// (POST /user/login)
func (svc *Service) Login(c echo.Context) error {
	svc.logger.Info("Login request received")

	// Parse the request body into the LoginRequest struct
	request := &LoginRequest{}
	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body format")
	}

	// You can now perform your signup logic here
	response, httpCode, err := svc.Services.UserSvc.Login(c.Request().Context(), *request)
	if err != nil {
		return c.JSON(httpCode, response)
	}

	// Return a response (e.g., a success message)
	return c.JSON(http.StatusOK, response)
}

// Verify - Verify
// (POST /user/verify)
func (svc *Service) VerifyUser(ctx echo.Context, params VerifyUserParams) error {
	svc.logger.Info("Verify request received")

	if params.Token == "" || params.UserPid == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request parameters")
	}
	// You can now perform your signup logic here
	response, httpCode, err := svc.Services.UserSvc.VerifyUser(ctx.Request().Context(), params)
	if err != nil {
		return ctx.JSON(httpCode, response)
	}

	// Return a response (e.g., a success message)
	return ctx.JSON(http.StatusOK, response)
}

package usersvc

import (
	"context"
	"errors"
	"net/http"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/authsvc"
	"github.com/GDGVIT/configgy-backend/constants"
	"github.com/GDGVIT/configgy-backend/pkg/crypto"
	"github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

func (svc *UserSvcImpl) Login(c context.Context, req api.LoginRequest) (api.LoginResponse, int, error) {
	var message string

	// print the email and password
	if req.Email == types.Email("") {
		message = "Email is required"
		return api.LoginResponse{
			Message: &message,
		}, http.StatusBadRequest, nil
	}

	if req.Password == "" {
		message = "Password is required"
		return api.LoginResponse{
			Message: &message,
		}, http.StatusBadRequest, nil
	}

	// Find user by email
	existingUser, err := svc.DB.GetUserByEmail(string(req.Email))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			message = "Internal Server Error"
			return api.LoginResponse{
				Message: &message,
			}, http.StatusInternalServerError, err
		}
	}

	if existingUser.Email == "" {
		message = "User does not exist"
		return api.LoginResponse{
			Message: &message,
		}, http.StatusBadRequest, nil
	}

	// Verify password
	if ok, err := crypto.VerifyPassword(string(existingUser.Password), req.Password); !ok {
		if err != nil {
			message = "Internal Server Error"
			return api.LoginResponse{
				Message: &message,
			}, http.StatusInternalServerError, err
		}

		message = "Incorrect password"
		return api.LoginResponse{
			Message: &message,
		}, http.StatusBadRequest, nil
	}

	// Create a new user
	// Create a new JWT token for the newly registered account

	var authData authsvc.TokenReq
	authData.UserID = existingUser.PID
	authData.Type = constants.TokenTypes.USER

	// Create a new token
	token, err := svc.authnSvc.GenerateToken(c, authData)
	if err != nil {
		message = "Internal Server Error"
		return api.LoginResponse{
			Message: &message,
		}, http.StatusInternalServerError, err
	}

	tokenData := api.AuthToken{
		AccessToken:           token.AccesssToken,
		AccessTokenExpiresIn:  token.AccessTokenExp,
		RefreshToken:          token.RefreshToken,
		RefreshTokenExpiresIn: token.RefreshTokenExp,
		UserId:                &existingUser.PID,
	}

	message = "Login successful"
	loginResponse := api.LoginResponse{
		Message: &message,
		Token:   &tokenData,
	}

	return loginResponse, http.StatusOK, nil
}

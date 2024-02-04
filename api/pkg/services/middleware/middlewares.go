package middleware

import (
	"context"
	"errors"

	"github.com/GDGVIT/configgy-backend/api/pkg/auth"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/authsvc"
	"github.com/GDGVIT/configgy-backend/constants"
)

func (g *MiddlewareImpl) CheckIfUser(c context.Context, token string) (*authsvc.AuthData, error) {
	var message string

	err := g.authnSvc.ValidateToken(token)
	if err != nil {
		message = "invalid token"
		return nil, errors.New(message)
	}

	// Check if user
	authData, err := auth.GetAuthDataFromToken(token)
	if err != nil {
		message = "invalid token"
		return nil, errors.New(message)
	}

	if authData.Type != constants.TokenTypes.USER {
		message = "invalid token"
		return nil, errors.New(message)
	}

	return authData, nil
}

func (g *MiddlewareImpl) CheckIfAdmin(c context.Context, token string) (*authsvc.AuthData, error) {
	var message string

	err := g.authnSvc.ValidateToken(token)
	if err != nil {
		message = "invalid token"
		return nil, errors.New(message)
	}

	// Check if user
	authData, err := auth.GetAuthDataFromToken(token)
	if err != nil {
		message = "invalid token"
		return nil, errors.New(message)
	}

	if authData.Type != constants.TokenTypes.ADMIN {
		message = "invalid token"
		return nil, errors.New(message)
	}

	return authData, nil
}

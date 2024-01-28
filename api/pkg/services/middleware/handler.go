package middleware

import (
	"context"

	"github.com/GDGVIT/configgy-backend/api/pkg/services/authsvc"
)

type MiddlewareImpl struct {
	authnSvc AuthSvc
}

type AuthSvc interface {
	ValidateToken(signedToken string) error
}

type Interface interface {
	CheckIfUser(c context.Context, token string) (*authsvc.AuthData, error)
	CheckIfAdmin(c context.Context, token string) (*authsvc.AuthData, error)
}

func Handler(authnSvc AuthSvc) Interface {
	return &MiddlewareImpl{
		authnSvc: authnSvc,
	}
}

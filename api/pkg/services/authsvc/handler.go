package authsvc

import (
	"context"

	"github.com/GDGVIT/configgy-backend/pkg/tables"
	"github.com/golang-jwt/jwt/v4"
)

type authSvcImpl struct {
	DB DB
}

func Handler(db DB) *authSvcImpl {
	return &authSvcImpl{DB: db}
}

type DB interface {
	GetUserByID(userId int) (*tables.Users, error)
	GetUserByPID(pid string) (*tables.Users, error)
}

// interface.
type Interface interface {
	GenerateToken(c context.Context, req TokenReq) (TokenRes, error)
	CreateToken(tokenAuthData AuthData) (*TokenDetails, error)
	ValidateToken(signedToken string) error
}

type AuthData struct {
	SessionPID string `json:"session_pid" binding:"required"`
	UserPID    string `json:"user_pid" binding:"required"`
	AdminPID   string `json:"admin_pid" binding:"required"`
	Sandbox    bool   `json:"sandbox" binding:"required"`
	Type       string `json:"type" binding:"required"`
	IsExternal bool   `json:"is_external" binding:"required"`
	jwt.RegisteredClaims
}

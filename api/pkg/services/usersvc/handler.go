package usersvc

import (
	"context"

	"github.com/GDGVIT/configgy-backend/api/pkg/services/authsvc"
	"github.com/GDGVIT/configgy-backend/pkg/logger"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

type UserSvcImpl struct {
	DB            UserDb
	logger        logger.Logger
	messageBroker MessageBroker
	authnSvc      AuthSvc
}

type AuthSvc interface {
	GenerateToken(c context.Context, req authsvc.TokenReq) (authsvc.TokenRes, error)
}

type UserDb interface {
	CreateUser(user *tables.Users) error
	GetUserByEmail(email string) (*tables.Users, error)
	GetUserByID(id int) (*tables.Users, error)
	GetUserByPID(pid string) (*tables.Users, error)

	CreateUserVerification(userVerification *tables.UserVerification) error
	GetUserVerificationByUserID(userID int) (*tables.UserVerification, error)
	UpdateUserVerification(userID int, isVerified bool) error
}

type MessageBroker interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
}

func Handler(userDb UserDb, logger logger.Logger, messageBroker MessageBroker, authsvcInterface AuthSvc) *UserSvcImpl {
	return &UserSvcImpl{
		DB:            userDb,
		logger:        logger,
		messageBroker: messageBroker,
		authnSvc:      authsvcInterface,
	}
}

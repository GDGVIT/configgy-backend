package usersvc

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/pkg/crypto"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
	"github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
)

func (svc *UserSvcImpl) SignUp(c context.Context, req api.SignupRequest) (api.GenericMessageResponse, int, error) {
	var message string

	// print the email and password
	if req.Email == types.Email("") {
		return api.GenericMessageResponse{
			Message: &message,
		}, http.StatusBadRequest, nil
	}

	if req.Password == "" {
		message = "Password is required"
		return api.GenericMessageResponse{
			Message: &message,
		}, http.StatusBadRequest, nil
	}

	// Find user by email
	existingUser, err := svc.DB.GetUserByEmail(string(req.Email))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		message = "Internal Server Error"
		return api.GenericMessageResponse{
			Message: &message,
		}, http.StatusInternalServerError, err
	}
	if existingUser != nil {
		message = "User already exists"
		return api.GenericMessageResponse{
			Message: &message,
		}, http.StatusBadRequest, nil
	}

	hashedPass, err := crypto.HashPassword(string(req.Password), nil)
	if err != nil {
		message = "Internal Server Error"
		return api.GenericMessageResponse{
			Message: &message,
		}, http.StatusInternalServerError, err
	}

	// Create a new user
	var user tables.Users
	user.Email = string(req.Email)
	user.Password = []byte(hashedPass)
	user.PublicKey = []byte(req.PublicKey)

	if err := svc.DB.CreateUser(&user); err != nil {
		message = "Failed to create user"
		return api.GenericMessageResponse{
			Message: &message,
		}, http.StatusInternalServerError, err
	}
	userVerificationToken, err := svc.generateVerificationToken()
	if err != nil {
		message = "Failed to generate verification token"
		return api.GenericMessageResponse{
			Message: &message,
		}, http.StatusInternalServerError, err
	}

	// Create a new user verification
	var userVerification tables.UserVerification
	userVerification.UserID = user.ID
	userVerification.Token = userVerificationToken

	if err := svc.DB.CreateUserVerification(&userVerification); err != nil {
		message = "Failed to create user verification"
		return api.GenericMessageResponse{
			Message: &message,
		}, http.StatusInternalServerError, err
	}

	// Create user personal vault
	vault := tables.Vault{
		PID:        tables.UUIDWithPrefix("vault"),
		Name:       user.PID,
		IsPersonal: true,
		PublicKey:  user.PublicKey,
	}

	err = svc.DB.CreateVault(vault, user.PID)
	if err != nil {
		message := "Failed to create user"
		return api.GenericMessageResponse{
			Success: false,
			Message: &message,
		}, http.StatusInternalServerError, err
	}

	message = "Signup successful. Please check your email for the verification link."
	var signupResponse api.GenericMessageResponse
	signupResponse.Message = &message

	var msg Message
	msg.From = os.Getenv("SMTP_USERNAME")
	msg.To = []string{string(req.Email)}
	msg.Subject = "Signup successful"
	msg.Body = "Click here to verify your email: http://" + os.Getenv("HOST") + "/v1/user/verify?token=" + userVerificationToken + "&user_pid=" + user.PID
	msg.Type = "text"

	exchange := "" // Use an empty exchange for direct exchange (default)
	routingKey := "mail"

	// Publish the message to the queue
	body, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Failed to marshal message to JSON: %v", err)
	}
	err = svc.messageBroker.Publish(c, exchange, routingKey, body)
	if err != nil {
		return signupResponse, http.StatusInternalServerError, err
	}

	return signupResponse, http.StatusOK, nil
}

type Message struct {
	From         string
	To           []string
	Subject      string
	Body         string
	TemplateName string
	Data         interface{}
	Type         string // template or text or html
}

func (svc *UserSvcImpl) generateVerificationToken() (string, error) {
	// generate a 6 digit random number
	randomString, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return randomString.String(), nil
}

package usersvc

import (
	"context"
	"net/http"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/pkg/errors"
)

func (svc *UserSvcImpl) VerifyUser(c context.Context, params api.VerifyUserParams) (api.GenericMessageResponse, int, error) {
	userPID := params.UserPid
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{}, http.StatusInternalServerError, err
	}
	if err != nil {
		return api.GenericMessageResponse{}, http.StatusBadRequest, err
	}
	userVerification, err := svc.DB.GetUserVerificationByUserID(user.ID)
	if err != nil {
		return api.GenericMessageResponse{}, http.StatusInternalServerError, err
	}
	if userVerification.Token != params.Token {
		return api.GenericMessageResponse{}, http.StatusBadRequest, errors.New("invalid token")
	}
	svc.DB.UpdateUserVerification(user.ID, true)
	message := "User successfully verified"
	return api.GenericMessageResponse{Message: &message}, http.StatusOK, nil
}

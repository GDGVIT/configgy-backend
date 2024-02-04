package authsvc

import (
	"context"

	"github.com/GDGVIT/configgy-backend/constants"
	"github.com/pkg/errors"
)

func (g *authSvcImpl) GenerateToken(c context.Context, req TokenReq) (TokenRes, error) {
	var res TokenRes
	var err error

	switch req.Type {

	case constants.TokenTypes.USER:
		{
			res, err := g.userTokenGeneration(c, req)
			if err != nil {
				return res, errors.Wrap(err, "[GenerateToken][userTokenGeneration]")
			}
			return res, err
		}
	case constants.TokenTypes.ADMIN:
		{
			res, err := g.adminTokenGeneration(c, req)
			if err != nil {
				return res, errors.Wrap(err, "[GenerateToken][adminTokenGeneration]")
			}
			return res, err
		}
	default:
		{
			err = errors.New("Invalid token type")
			return res, err
		}
	}
}

func (g *authSvcImpl) userTokenGeneration(c context.Context, req TokenReq) (TokenRes, error) {

	var res TokenRes
	var err error
	var authData AuthData

	authData.Type = req.Type

	//get sandbox

	userData, err := g.DB.GetUserByPID(req.UserID)
	if err != nil {
		return res, errors.Wrap(err, "[onboardingTokenGeneration][GetUserByPID]")
	}
	authData.UserPID = userData.PID

	tokenRes, err := g.CreateToken(authData)
	if err != nil {
		return res, errors.Wrap(err, "[onboardingTokenGeneration][CreateToken][onboarding]")
	}
	res.AccesssToken = tokenRes.AccessToken
	res.RefreshToken = tokenRes.RefreshToken
	res.AccessTokenExp = tokenRes.AtExpires
	res.RefreshTokenExp = tokenRes.RtExpires
	res.AccesssTokenPID = tokenRes.TokenUuid
	res.RefreshTokenPID = tokenRes.RefreshUuid
	res.Type = req.Type
	res.UserID = userData.PID

	return res, err
}

func (g *authSvcImpl) adminTokenGeneration(c context.Context, req TokenReq) (TokenRes, error) {
	var res TokenRes
	var err error
	var authData AuthData

	authData.Type = req.Type

	adminData, err := g.DB.GetUserByPID(req.UserID)
	if err != nil {
		return res, errors.Wrap(err, "[onboardingTokenGeneration][GetUserDetailsByPID]")
	}
	authData.AdminPID = adminData.PID

	tokenRes, err := g.CreateToken(authData)
	if err != nil {
		return res, errors.Wrap(err, "[onboardingTokenGeneration][CreateToken][onboarding]")
	}
	res.AccesssToken = tokenRes.AccessToken
	res.RefreshToken = tokenRes.RefreshToken
	res.AccessTokenExp = tokenRes.AtExpires
	res.RefreshTokenExp = tokenRes.RtExpires
	res.AccesssTokenPID = tokenRes.TokenUuid
	res.RefreshTokenPID = tokenRes.RefreshUuid
	res.Type = req.Type
	res.AdminID = adminData.PID

	return res, err
}

package auth

import (
	"context"
	"errors"

	repositoryauth "github.com/twrnakata/user-auth-api/internal/repository/auth"
	repositorymodel "github.com/twrnakata/user-auth-api/internal/repository/auth/model"
	servicemodel "github.com/twrnakata/user-auth-api/internal/service/auth/model"
	"golang.org/x/crypto/bcrypt"

	domainauth "github.com/twrnakata/user-auth-api/internal/domain/auth"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
)

type AuthLoginService struct {
	Repository   domainauth.AuthLoginRepository
	TokenService TokenService
}

type TokenService interface {
	CreateToken(userID string, name string) (string, error)
}

func (service *AuthLoginService) Login(executionContext context.Context, request *servicemodel.LoginUserRequestModel, response *servicemodel.LoginUserResponseModel) error {
	if service.Repository == nil {
		return apperror.ErrAuthLoginRepositoryNotConfigured
	}
	if service.TokenService == nil {
		return apperror.ErrTokenServiceNotConfigured
	}

	repositoryRequest := &repositorymodel.AuthLoginRepositoryRequestModel{
		Email: request.Email,
	}

	getLoginUserByEmailModel := &repositorymodel.GetLoginUserByEmailModel{}
	if err := service.Repository.GetLoginUserByEmail(executionContext, repositoryRequest, getLoginUserByEmailModel); err != nil {
		if errors.Is(err, repositoryauth.ErrNotFound) {
			// Same error as a wrong password: do not reveal that the email is missing.
			return domainauth.ErrInvalidCredentials
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(getLoginUserByEmailModel.PasswordHash), []byte(request.Password)); err != nil {
		return domainauth.ErrInvalidCredentials
	}

	token, err := service.TokenService.CreateToken(getLoginUserByEmailModel.ID, getLoginUserByEmailModel.Name)
	if err != nil {
		return err
	}

	response.Token = token
	response.ID = getLoginUserByEmailModel.ID
	response.Name = getLoginUserByEmailModel.Name
	return nil
}

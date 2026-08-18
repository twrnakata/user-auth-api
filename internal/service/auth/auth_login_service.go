package auth

import (
	"context"
	"errors"

	repositoryauth "backend-challenge-golang/internal/repository/auth"
	repositorymodel "backend-challenge-golang/internal/repository/auth/model"
	servicemodel "backend-challenge-golang/internal/service/auth/model"
	"golang.org/x/crypto/bcrypt"

	domainauth "backend-challenge-golang/internal/domain/auth"
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
		return errors.New("auth login repository not configured")
	}
	if service.TokenService == nil {
		return errors.New("token service not configured")
	}

	repositoryRequest := &repositorymodel.AuthLoginRepositoryRequestModel{
		Email: request.Email,
	}

	getLoginUserByEmailModel := &repositorymodel.GetLoginUserByEmailModel{}
	if err := service.Repository.GetLoginUserByEmail(executionContext, repositoryRequest, getLoginUserByEmailModel); err != nil {
		if errors.Is(err, repositoryauth.ErrNotFound) {
			return domainauth.ErrUserNotFound
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

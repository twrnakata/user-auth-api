package auth

import (
	"context"
	"errors"

	repositorymodel "backend-challenge-golang/internal/repository/auth/model"
	servicemodel "backend-challenge-golang/internal/service/auth/model"
	"golang.org/x/crypto/bcrypt"

	domainauth "backend-challenge-golang/internal/domain/auth"
)

type AuthLoginService struct {
	Repository domainauth.AuthLoginRepository
	BuildToken func(userID string, name string) string
}

func (service *AuthLoginService) Login(executionContext context.Context, request *servicemodel.LoginUserRequest, response *servicemodel.LoginUserResponse) error {
	if service.Repository == nil {
		return errors.New("auth login repository not configured")
	}

	repositoryRequest := &repositorymodel.AuthLoginRepositoryRequest{
		Email: request.Email,
	}

	repositoryResponse := &repositorymodel.AuthLoginRepositoryResponse{}
	if err := service.Repository.GetByEmail(executionContext, repositoryRequest, repositoryResponse); err != nil {
		return err
	}

	if repositoryResponse.Email == "" || repositoryResponse.PasswordHash == "" {
		return domainauth.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(repositoryResponse.PasswordHash), []byte(request.Password)); err != nil {
		return domainauth.ErrInvalidCredentials
	}

	token := ""
	if service.BuildToken != nil {
		token = service.BuildToken(repositoryResponse.ID, repositoryResponse.Name)
	}

	response.Token = token
	response.ID = repositoryResponse.ID
	response.Name = repositoryResponse.Name
	return nil
}

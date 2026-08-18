package auth

import (
	"context"
	"errors"

	repositorymodel "backend-challenge-golang/internal/repository/auth/model"
	servicemodel "backend-challenge-golang/internal/service/auth/model"
	"golang.org/x/crypto/bcrypt"

	domainauth "backend-challenge-golang/internal/domain/auth"
)

// RegisterUserService is the application-layer implementation of the domain port.
// It contains the register business rules (password hashing, calling repository).
type RegisterUserService struct {
	Repo domainauth.RegisterUserRepository
}

func (service *RegisterUserService) Register(executionContext context.Context, request *servicemodel.RegisterUserRequest, response *servicemodel.RegisterUserResponse) error {
	if service.Repo == nil {
		return errors.New("register repo not configured")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	createRequest := &repositorymodel.CreateRegisterUserRequest{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(hash),
	}

	createResp := &repositorymodel.CreateRegisterUserResponse{}
	err = service.Repo.Create(executionContext, createRequest, createResp)
	if err != nil {
		return err
	}

	response.ID = createResp.ID
	response.Name = createResp.Name
	response.Email = createResp.Email
	response.CreatedAt = createResp.CreatedAt
	return nil
}

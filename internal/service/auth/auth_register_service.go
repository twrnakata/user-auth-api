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

// RegisterUserService is the application-layer implementation of the domain port.
// It contains the register business rules (password hashing, calling repository).
type RegisterUserService struct {
	Repository domainauth.RegisterUserRepository
}

func (service *RegisterUserService) Register(executionContext context.Context, request *servicemodel.RegisterUserRequest, response *servicemodel.RegisterUserResponse) error {
	if service.Repository == nil {
		return errors.New("register repository not configured")
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
	err = service.Repository.CreateRegisterUser(executionContext, createRequest, createResp)
	if err != nil {
		if errors.Is(err, repositoryauth.ErrDuplicateKey) {
			return domainauth.ErrEmailAlreadyExists
		}
		return err
	}

	response.ID = createResp.ID
	response.Name = createResp.Name
	response.Email = createResp.Email
	response.CreatedAt = createResp.CreatedAt
	return nil
}

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

func (service *RegisterUserService) Register(executionContext context.Context, request *servicemodel.RegisterUserRequestModel, response *servicemodel.RegisterUserResponseModel) error {
	if service.Repository == nil {
		return errors.New("register repository not configured")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	createRegisterUserRequestModel := &repositorymodel.CreateRegisterUserRequestModel{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(hash),
	}

	createRegisterUserModel := &repositorymodel.CreateRegisterUserModel{}
	err = service.Repository.CreateRegisterUser(executionContext, createRegisterUserRequestModel, createRegisterUserModel)
	if err != nil {
		if errors.Is(err, repositoryauth.ErrDuplicateKey) {
			return domainauth.ErrEmailAlreadyExists
		}
		return err
	}

	response.ID = createRegisterUserModel.ID
	response.Name = createRegisterUserModel.Name
	response.Email = createRegisterUserModel.Email
	response.CreatedAt = createRegisterUserModel.CreatedAt
	return nil
}

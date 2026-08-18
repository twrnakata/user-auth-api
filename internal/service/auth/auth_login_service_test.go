package auth

import (
	"context"
	"errors"
	"testing"

	domainauth "backend-challenge-golang/internal/domain/auth"
	repositoryauth "backend-challenge-golang/internal/repository/auth"
	repositorymodel "backend-challenge-golang/internal/repository/auth/model"
	servicemodel "backend-challenge-golang/internal/service/auth/model"
	jwtpkg "backend-challenge-golang/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type fakeTokenService struct {
	token string
	err   error
}

func (service *fakeTokenService) CreateToken(userID string, name string) (string, error) {
	return service.token, service.err
}

type fakeAuthLoginRepository struct {
	called   bool
	gotEmail string
	response *repositorymodel.GetLoginUserByEmailModel
	err      error
}

func (repository *fakeAuthLoginRepository) GetLoginUserByEmail(executionContext context.Context, request *repositorymodel.AuthLoginRepositoryRequestModel, response *repositorymodel.GetLoginUserByEmailModel) error {
	repository.called = true

	if request != nil {
		repository.gotEmail = request.Email
	}

	if response != nil && repository.response != nil {
		*response = *repository.response
	}

	return repository.err
}

func TestAuthLoginService_Login_FillsResponseFromRepository(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("expected no error generating password hash, got: %v", err)
	}

	repository := &fakeAuthLoginRepository{
		response: &repositorymodel.GetLoginUserByEmailModel{
			ID:           "u-1",
			Name:         "Alice",
			Email:        "alice@example.com",
			PasswordHash: string(passwordHash),
		},
	}

	service := &AuthLoginService{
		Repository: repository,
		TokenService: &fakeTokenService{
			token: "mock-token",
		},
	}

	request := &servicemodel.LoginUserRequestModel{
		Email:    "alice@example.com",
		Password: "secret",
	}
	response := &servicemodel.LoginUserResponseModel{}

	err = service.Login(context.Background(), request, response)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !repository.called {
		t.Fatalf("expected repository to be called")
	}

	if repository.gotEmail != "alice@example.com" {
		t.Fatalf("unexpected email passed to repository: %s", repository.gotEmail)
	}

	if response.Token != "mock-token" || response.ID != "u-1" || response.Name != "Alice" {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestAuthLoginService_Login_NotFound_ReturnsErrUserNotFound(t *testing.T) {
	repository := &fakeAuthLoginRepository{
		err: repositoryauth.ErrNotFound,
	}

	service := &AuthLoginService{
		Repository:   repository,
		TokenService: &fakeTokenService{token: "unused"},
	}

	err := service.Login(context.Background(), &servicemodel.LoginUserRequestModel{
		Email:    "unknown@example.com",
		Password: "secret",
	}, &servicemodel.LoginUserResponseModel{})
	if !errors.Is(err, domainauth.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestAuthLoginService_Login_TokenServiceNil_DoesNotQueryRepository(t *testing.T) {
	repository := &fakeAuthLoginRepository{}
	service := &AuthLoginService{
		Repository: repository,
	}

	err := service.Login(context.Background(), &servicemodel.LoginUserRequestModel{
		Email:    "alice@example.com",
		Password: "secret",
	}, &servicemodel.LoginUserResponseModel{})
	if err == nil || err.Error() != "token service not configured" {
		t.Fatalf("expected token service not configured, got: %v", err)
	}
	if repository.called {
		t.Fatalf("expected repository not to be called when token service is nil")
	}
}

func TestAuthLoginService_Login_CreatesRealJWT(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("expected no error generating password hash, got: %v", err)
	}

	jwtService, err := jwtpkg.NewJWTService("test-secret", jwtpkg.DefaultExpireDuration)
	if err != nil {
		t.Fatalf("expected no error creating jwt service, got: %v", err)
	}

	service := &AuthLoginService{
		Repository: &fakeAuthLoginRepository{
			response: &repositorymodel.GetLoginUserByEmailModel{
				ID:           "u-1",
				Name:         "Alice",
				Email:        "alice@example.com",
				PasswordHash: string(passwordHash),
			},
		},
		TokenService: jwtService,
	}

	response := &servicemodel.LoginUserResponseModel{}
	err = service.Login(context.Background(), &servicemodel.LoginUserRequestModel{
		Email:    "alice@example.com",
		Password: "secret",
	}, response)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	claims, err := jwtService.ParseToken(response.Token)
	if err != nil {
		t.Fatalf("expected token to parse, got: %v", err)
	}
	if claims.UserID != "u-1" || claims.Name != "Alice" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}

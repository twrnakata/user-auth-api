package apperror

import "errors"

var (
	ErrRegisterServiceNotInitialized   = errors.New("register service not initialized")
	ErrLoginServiceNotInitialized      = errors.New("login service not initialized")
	ErrListUserServiceNotInitialized   = errors.New("list user service not initialized")
	ErrGetUserServiceNotInitialized    = errors.New("get user service not initialized")
	ErrUpdateUserServiceNotInitialized = errors.New("update user service not initialized")
	ErrDeleteUserServiceNotInitialized = errors.New("delete user service not initialized")
	ErrJWTServiceNotInitialized        = errors.New("jwt service not initialized")
	ErrCountUserServiceNil             = errors.New("count user service is nil")

	ErrInvalidJSONBody                  = errors.New("invalid json body")
	ErrIDRequired                       = errors.New("id is required")
	ErrNameOrEmailRequired              = errors.New("name or email is required")
	ErrEmailAndPasswordRequired         = errors.New("email and password are required")
	ErrNameEmailPasswordRequired        = errors.New("name, email, password are required")
	ErrInvalidEmail                     = errors.New("invalid email format")
	ErrInvalidUserID                    = errors.New("invalid user id")
	ErrMissingAuthorizationHeader       = errors.New("missing authorization header")
	ErrInvalidAuthorizationHeaderFormat = errors.New("invalid authorization header format")
	ErrMissingBearerToken               = errors.New("missing bearer token")
	ErrInvalidOrExpiredToken            = errors.New("invalid or expired token")

	ErrUserCollectionNil                   = errors.New("user collection is nil")
	ErrRegisterRepositoryNotConfigured     = errors.New("register repository not configured")
	ErrAuthRegisterRepositoryNotConfigured = errors.New("auth register repository not configured")
	ErrAuthLoginRepositoryNotConfigured    = errors.New("auth login repository not configured")
	ErrTokenServiceNotConfigured           = errors.New("token service not configured")
	ErrListUserRepositoryNotConfigured     = errors.New("list user repository not configured")
	ErrGetUserRepositoryNotConfigured      = errors.New("get user repository not configured")
	ErrUpdateUserRepositoryNotConfigured   = errors.New("update user repository not configured")
	ErrDeleteUserRepositoryNotConfigured   = errors.New("delete user repository not configured")
	ErrCountUserRepositoryNotConfigured    = errors.New("count user repository not configured")

	ErrUserResponseNil  = errors.New("user response is nil")
	ErrUsersResponseNil = errors.New("users response is nil")
	ErrCountResponseNil = errors.New("count response is nil")

	ErrMongoURIRequired        = errors.New("MONGO_URI is required")
	ErrMongoDatabaseRequired   = errors.New("MONGO_DATABASE is required")
	ErrJWTSecretRequired       = errors.New("JWT_SECRET is required")
	ErrMongoURIEmpty           = errors.New("mongo uri is empty")
	ErrJWTSecretKeyRequired    = errors.New("jwt secret key is required")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)

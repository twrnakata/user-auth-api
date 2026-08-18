package model

type AuthLoginRepositoryRequest struct {
	Email string
}

type AuthLoginRepositoryResponse struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

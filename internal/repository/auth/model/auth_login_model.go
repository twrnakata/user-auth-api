package model

import "go.mongodb.org/mongo-driver/v2/bson"

type AuthLoginRepositoryRequest struct {
	Email string
}

type AuthLoginRepositoryResponse struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

type GetLoginUserDocument struct {
	ID           bson.ObjectID `bson:"_id"`
	Name         string        `bson:"name"`
	Email        string        `bson:"email"`
	PasswordHash string        `bson:"passwordHash"`
}

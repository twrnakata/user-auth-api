package model

import "go.mongodb.org/mongo-driver/v2/bson"

type AuthLoginRepositoryRequestModel struct {
	Email string
}

type GetLoginUserByEmailModel struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

type GetLoginUserDocumentModel struct {
	ID           bson.ObjectID `bson:"_id"`
	Name         string        `bson:"name"`
	Email        string        `bson:"email"`
	PasswordHash string        `bson:"passwordHash"`
}

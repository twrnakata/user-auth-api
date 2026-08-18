package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UpdateUserDocumentModel struct {
	ID        bson.ObjectID `bson:"_id"`
	Name      string        `bson:"name"`
	Email     string        `bson:"email"`
	CreatedAt time.Time     `bson:"createdAt"`
}

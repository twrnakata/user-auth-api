package model

import "time"

type GetUserRequestModel struct {
	ID string
}

type GetUserResponseModel struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

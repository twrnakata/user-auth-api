package model

type UserGetRequestModel struct {
	ID string
}

type UserGetResponseModel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

package model

type UserDeleteRequestModel struct {
	ID string
}

type UserDeleteResponseModel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

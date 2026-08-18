package model

type UserUpdateRequestModel struct {
	ID    string
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserUpdateResponseModel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

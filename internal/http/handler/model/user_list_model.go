package model

type UserListItemModel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type UserListResponseModel struct {
	Users []UserListItemModel `json:"users"`
}

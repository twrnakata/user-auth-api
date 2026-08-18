package user

import "context"

type GetUserRepository interface {
	GetUserByID(executionContext context.Context, userID string, user *User) error
}

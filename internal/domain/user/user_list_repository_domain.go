package user

import "context"

type ListUserRepository interface {
	ListUsers(executionContext context.Context, users *[]User) error
}

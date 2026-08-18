package user

import "context"

type CountUserRepository interface {
	CountUsers(executionContext context.Context, count *int64) error
}

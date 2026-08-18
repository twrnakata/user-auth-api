package user

import (
	"context"
	"errors"
	"testing"
)

type fakeCountUserRepository struct {
	called bool
	count  int64
	err    error
}

func (repository *fakeCountUserRepository) CountUsers(executionContext context.Context, count *int64) error {
	repository.called = true
	if count != nil {
		*count = repository.count
	}
	return repository.err
}

func TestCountUserService_Count_FillsCountFromRepository(t *testing.T) {
	repository := &fakeCountUserRepository{
		count: 3,
	}

	service := &CountUserService{
		Repository: repository,
	}

	var count int64
	err := service.Count(context.Background(), &count)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !repository.called {
		t.Fatalf("expected repository CountUsers to be called")
	}
	if count != 3 {
		t.Fatalf("unexpected count: %d", count)
	}
}

func TestCountUserService_Count_NilCount_ReturnsError(t *testing.T) {
	repository := &fakeCountUserRepository{}

	service := &CountUserService{
		Repository: repository,
	}

	err := service.Count(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if repository.called {
		t.Fatalf("expected repository CountUsers not to be called")
	}
}

func TestCountUserService_Count_RepositoryError_ReturnsSameError(t *testing.T) {
	repositoryError := errors.New("network timeout")
	repository := &fakeCountUserRepository{
		err: repositoryError,
	}

	service := &CountUserService{
		Repository: repository,
	}

	var count int64
	err := service.Count(context.Background(), &count)
	if !errors.Is(err, repositoryError) {
		t.Fatalf("expected repository error, got: %v", err)
	}
}

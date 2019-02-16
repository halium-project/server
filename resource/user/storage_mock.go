package user

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (t *StorageMock) Set(_ context.Context, userID string, rev string, value *User) (string, error) {
	args := t.Called(userID, rev, value)

	return args.String(0), args.Error(1)
}

func (t *StorageMock) Get(_ context.Context, userID string) (string, *User, error) {
	args := t.Called(userID)

	if args.Get(1) == nil {
		return "", nil, args.Error(2)
	}

	return args.String(0), args.Get(1).(*User), args.Error(2)
}

func (t *StorageMock) GetAll(ctx context.Context) (map[string]User, error) {
	args := t.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]User), nil
}

func (t *StorageMock) FindOneByUsername(ctx context.Context, username string) (string, string, *User, error) {
	args := t.Called(username)

	if args.Get(2) == nil {
		return "", "", nil, args.Error(3)
	}

	return args.String(0), args.String(1), args.Get(2).(*User), args.Error(3)
}

func (t *StorageMock) FindTotalUserCount(_ context.Context) (int, error) {
	args := t.Called()

	return args.Int(0), args.Error(1)
}

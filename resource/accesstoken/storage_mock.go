package accesstoken

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (t *StorageMock) Set(_ context.Context, code string, rev string, value *AccessToken) (string, error) {
	value.CreatedAt = ValidAccessToken.CreatedAt

	args := t.Called(code, rev, value)

	return args.String(0), args.Error(1)
}

func (t *StorageMock) Get(_ context.Context, code string) (string, *AccessToken, error) {
	args := t.Called(code)

	if args.Get(1) == nil {
		return "", nil, args.Error(2)
	}

	return args.String(0), args.Get(1).(*AccessToken), args.Error(2)
}

func (t *StorageMock) FindOneByRefreshToken(ctx context.Context, refreshToken string) (string, string, *AccessToken, error) {
	args := t.Called(refreshToken)

	if args.Get(2) == nil {
		return "", "", nil, args.Error(3)
	}

	return args.String(0), args.String(1), args.Get(2).(*AccessToken), args.Error(3)
}

func (t *StorageMock) Delete(_ context.Context, code string, rev string) error {
	return t.Called(code, rev).Error(0)
}

package authorizationcode

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (t *StorageMock) Set(_ context.Context, code string, rev string, value *AuthorizationCode) (string, error) {
	value.CreatedAt = ValidAuthorizationCode.CreatedAt

	args := t.Called(code, rev, value)

	return args.String(0), args.Error(1)
}

func (t *StorageMock) Get(_ context.Context, code string) (string, *AuthorizationCode, error) {
	args := t.Called(code)

	if args.Get(1) == nil {
		return "", nil, args.Error(2)
	}

	return args.String(0), args.Get(1).(*AuthorizationCode), args.Error(2)
}

func (t *StorageMock) Delete(_ context.Context, code string, rev string) error {
	return t.Called(code, rev).Error(0)
}

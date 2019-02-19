package client

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (t *StorageMock) Set(_ context.Context, id string, rev string, value *Client) (string, error) {
	args := t.Called(id, rev, value)

	return args.String(0), args.Error(1)
}

func (t *StorageMock) Get(_ context.Context, id string) (string, *Client, error) {
	args := t.Called(id)

	if args.Get(1) == nil {
		return "", nil, args.Error(2)
	}

	return args.String(0), args.Get(1).(*Client), args.Error(2)
}

func (t *StorageMock) Delete(ctx context.Context, id string) error {
	return t.Called(id).Error(0)
}

func (t *StorageMock) GetAll(ctx context.Context) (map[string]Client, error) {
	args := t.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]Client), nil
}

func (t *StorageMock) FindOneByName(ctx context.Context, name string) (string, string, *Client, error) {
	args := t.Called(name)

	if args.Get(2) == nil {
		return "", "", nil, args.Error(3)
	}

	return args.String(0), args.String(1), args.Get(2).(*Client), args.Error(3)
}

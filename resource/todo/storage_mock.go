package todo

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (t *StorageMock) Set(_ context.Context, id string, rev string, value *Todo) (string, error) {
	args := t.Called(id, rev, value)

	return args.String(0), args.Error(1)
}

func (t *StorageMock) Get(_ context.Context, id string) (string, *Todo, error) {
	args := t.Called(id)

	if args.Get(1) == nil {
		return "", nil, args.Error(2)
	}

	return args.String(0), args.Get(1).(*Todo), args.Error(2)
}

func (t *StorageMock) Delete(ctx context.Context, id string) error {
	return t.Called(id).Error(0)
}

func (t *StorageMock) GetAll(ctx context.Context) (map[string]Todo, error) {
	args := t.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]Todo), nil
}

func (t *StorageMock) FindOneByTitle(ctx context.Context, name string) (string, string, *Todo, error) {
	args := t.Called(name)

	if args.Get(2) == nil {
		return "", "", nil, args.Error(3)
	}

	return args.String(0), args.String(1), args.Get(2).(*Todo), args.Error(3)
}

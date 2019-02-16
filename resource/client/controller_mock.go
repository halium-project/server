package client

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ControllerMock struct {
	mock.Mock
}

func (t *ControllerMock) Create(ctx context.Context, cmd *CreateCmd) error {
	return t.Called(cmd).Error(0)
}

func (t *ControllerMock) Get(ctx context.Context, cmd *GetCmd) (*Client, error) {
	args := t.Called(cmd)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Client), args.Error(1)
}

func (t *ControllerMock) GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]Client, error) {
	args := t.Called(cmd)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]Client), args.Error(1)
}

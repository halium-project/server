package accesstoken

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ControllerMock struct {
	mock.Mock
}

func (t *ControllerMock) Create(ctx context.Context, cmd *CreateCmd) error {
	args := t.Called(cmd)

	return args.Error(0)
}

func (t *ControllerMock) Get(ctx context.Context, cmd *GetCmd) (*AccessToken, error) {
	args := t.Called(cmd)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*AccessToken), args.Error(1)
}

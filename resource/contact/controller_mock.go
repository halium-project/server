package contact

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ControllerMock struct {
	mock.Mock
}

func (t *ControllerMock) Create(ctx context.Context, cmd *CreateCmd) (string, error) {
	args := t.Called(cmd)

	return args.String(0), args.Error(1)
}

func (t *ControllerMock) Get(ctx context.Context, cmd *GetCmd) (*Contact, error) {
	args := t.Called(cmd)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Contact), args.Error(1)
}

func (t *ControllerMock) GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]Contact, error) {
	args := t.Called(cmd)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]Contact), args.Error(1)
}

func (t *ControllerMock) Delete(ctx context.Context, cmd *DeleteCmd) error {
	return t.Called(cmd).Error(0)
}

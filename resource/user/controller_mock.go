package user

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

func (t *ControllerMock) Get(ctx context.Context, cmd *GetCmd) (*User, error) {
	args := t.Called(cmd)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

func (t *ControllerMock) Delete(ctx context.Context, cmd *DeleteCmd) error {
	return t.Called(cmd).Error(0)
}

func (t *ControllerMock) GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]User, error) {
	args := t.Called(cmd)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]User), args.Error(1)
}

func (t *ControllerMock) Validate(ctx context.Context, cmd *ValidateCmd) (string, *User, error) {
	args := t.Called(cmd)

	if args.Get(1) == nil {
		return "", nil, args.Error(2)
	}

	return args.String(0), args.Get(1).(*User), args.Error(2)
}

func (t *ControllerMock) GetTotalUserCount(ctx context.Context) (int, error) {
	args := t.Called()

	return args.Int(0), args.Error(1)
}

func (t *ControllerMock) Update(ctx context.Context, cmd *UpdateCmd) error {
	return t.Called(cmd).Error(0)
}

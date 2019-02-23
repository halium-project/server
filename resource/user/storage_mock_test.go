package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_User_StorageMock_Set(t *testing.T) {
	mock := new(StorageMock)

	mock.On("Set", "some-id", "some-rev", &ValidUser).Return("some-rev-2", nil)

	rev, err := mock.Set(context.Background(), "some-id", "some-rev", &ValidUser)

	assert.NoError(t, err)
	assert.Equal(t, "some-rev-2", rev)

	mock.AssertExpectations(t)
}

func Test_User_StorageMock_Get(t *testing.T) {
	mock := new(StorageMock)

	mock.On("Get", "some-id").Return("some-rev", &ValidUser, nil)

	rev, res, err := mock.Get(context.Background(), "some-id")
	assert.NoError(t, err)
	assert.Equal(t, "some-rev", rev)
	assert.EqualValues(t, &ValidUser, res)

	mock.AssertExpectations(t)
}

func Test_User_StorageMock_Get_with_error(t *testing.T) {
	mock := new(StorageMock)

	mock.On("Get", "some-id").Return("", nil, errors.New("some-error"))

	rev, res, err := mock.Get(context.Background(), "some-id")
	assert.EqualError(t, err, "some-error")
	assert.Empty(t, rev)
	assert.Nil(t, res)

	mock.AssertExpectations(t)
}

func Test_User_StorageMock_Delete(t *testing.T) {
	mock := new(StorageMock)

	mock.On("Delete", "some-id").Return(nil)

	err := mock.Delete(context.Background(), "some-id")

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

func Test_User_StorageMock_GetAll(t *testing.T) {
	mock := new(StorageMock)

	mock.On("GetAll").Return(map[string]User{
		"some-id": ValidUser,
	}, nil).Once()

	res, err := mock.GetAll(context.Background())

	assert.NoError(t, err)
	assert.EqualValues(t, map[string]User{
		"some-id": ValidUser,
	}, res)

	mock.AssertExpectations(t)
}

func Test_User_StorageMock_GetAll_with_an_error(t *testing.T) {
	mock := new(StorageMock)

	mock.On("GetAll").Return(nil, fmt.Errorf("some-error")).Once()

	res, err := mock.GetAll(context.Background())

	assert.Empty(t, res)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_User_StorageMock_FindOneByUsername(t *testing.T) {
	mock := new(StorageMock)

	mock.On("FindOneByUsername", "some-username").Return("some-id", "some-rev", &ValidUser, nil)

	id, rev, res, err := mock.FindOneByUsername(context.Background(), "some-username")

	assert.NoError(t, err)
	assert.Equal(t, "some-id", id)
	assert.Equal(t, "some-rev", rev)
	assert.EqualValues(t, &ValidUser, res)

	mock.AssertExpectations(t)
}

func Test_User_StorageMock_FindOneByUsername_with_error(t *testing.T) {
	mock := new(StorageMock)

	mock.On("FindOneByUsername", "some-username").Return("", "", nil, errors.New("some-error"))

	id, rev, res, err := mock.FindOneByUsername(context.Background(), "some-username")

	assert.Empty(t, rev)
	assert.Empty(t, id)
	assert.Nil(t, res)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

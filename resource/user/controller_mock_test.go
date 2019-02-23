package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_User_ControllerMock_Create(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Create", &CreateCmd{
		Username: "some-username",
		Password: "some-password",
	}).Return("some-user-id", nil).Once()

	userID, err := mock.Create(context.Background(), &CreateCmd{
		Username: "some-username",
		Password: "some-password",
	})

	assert.NoError(t, err)
	assert.Equal(t, "some-user-id", userID)

	mock.AssertExpectations(t)
}

func Test_User_ControllerMock_Create_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Create", &CreateCmd{
		Username: "some-username",
		Password: "some-password",
	}).Return("", fmt.Errorf("some-error")).Once()

	userID, err := mock.Create(context.Background(), &CreateCmd{
		Username: "some-username",
		Password: "some-password",
	})

	assert.Empty(t, userID)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_User_ControllerMock_Get(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		UserID: "some-user-id",
	}).Return(&ValidUser, nil).Once()

	res, err := mock.Get(context.Background(), &GetCmd{
		UserID: "some-user-id",
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidUser, res)

	mock.AssertExpectations(t)
}

func Test_User_ControllerMock_Get_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		UserID: "some-user-id",
	}).Return(nil, fmt.Errorf("some-error")).Once()

	session, err := mock.Get(context.Background(), &GetCmd{
		UserID: "some-user-id",
	})

	assert.Nil(t, session)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_User_ControllerMock_GetAll(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("GetAll", &GetAllCmd{}).Return(map[string]User{
		"some-id": ValidUser,
	}, nil).Once()

	res, err := mock.GetAll(context.Background(), &GetAllCmd{})

	assert.NoError(t, err)
	assert.EqualValues(t, map[string]User{
		"some-id": ValidUser,
	}, res)

	mock.AssertExpectations(t)
}

func Test_User_ControllerMock_GetAll_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("GetAll", &GetAllCmd{}).Return(nil, fmt.Errorf("some-error")).Once()

	user, err := mock.GetAll(context.Background(), &GetAllCmd{})

	assert.Nil(t, user)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_User_ControllerMock_Validate(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Validate", &ValidateCmd{
		Username: "some-email",
		Password: "some-clear-password",
	}).Return("some-user-id", &ValidUser, nil).Once()

	userID, user, err := mock.Validate(context.Background(), &ValidateCmd{
		Username: "some-email",
		Password: "some-clear-password",
	})

	assert.NoError(t, err)
	assert.Equal(t, "some-user-id", userID)
	assert.EqualValues(t, &ValidUser, user)

	mock.AssertExpectations(t)
}

func Test_User_ControllerMock_Validate_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Validate", &ValidateCmd{
		Username: "some-email",
		Password: "some-clear-password",
	}).Return("", nil, fmt.Errorf("some-error")).Once()

	userID, user, err := mock.Validate(context.Background(), &ValidateCmd{
		Username: "some-email",
		Password: "some-clear-password",
	})

	assert.Empty(t, userID)
	assert.Nil(t, user)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_User_ControllerMock_GetTotalUserCount(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("GetTotalUserCount").Return(42, nil).Once()

	res, err := mock.GetTotalUserCount(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 42, res)

	mock.AssertExpectations(t)
}

func Test_User_ControllerMock_Update(t *testing.T) {
	mock := new(ControllerMock)

	cmd := UpdateCmd{
		UserID:   "some-user-id",
		Username: "some-username",
		Role:     "some-role",
	}

	mock.On("Update", &cmd).Return(nil).Once()

	err := mock.Update(context.Background(), &cmd)

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

func Test_Client_ControllerMock_Delete(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Delete", &DeleteCmd{UserID: "some-id"}).Return(nil).Once()

	err := mock.Delete(context.Background(), &DeleteCmd{UserID: "some-id"})

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

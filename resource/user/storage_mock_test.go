package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/server/db"
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

func Test_User_StorageMock_GetAll(t *testing.T) {
	dbDriver := new(db.DriverMock)
	service := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_username",
		Limit:     200,
	}).Return([]db.ViewRow{
		{ID: "some-id"},
		{ID: "some-id-2"},
	}, nil).Once()

	dbDriver.On("GetMany", []string{"some-id", "some-id-2"}).Return(map[string]User{
		"some-id":   ValidUser,
		"some-id-2": ValidUser,
	}, nil).Once()

	res, err := service.GetAll(context.Background())

	assert.NoError(t, err)
	assert.EqualValues(t, map[string]User{
		"some-id":   ValidUser,
		"some-id-2": ValidUser,
	}, res)

	dbDriver.AssertExpectations(t)
}

func Test_User_StorageMock_GetAll_empty(t *testing.T) {
	dbDriver := new(db.DriverMock)
	service := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_username",
		Limit:     200,
	}).Return([]db.ViewRow{}, nil).Once()

	res, err := service.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, res)

	dbDriver.AssertExpectations(t)
}

func Test_User_StorageMock_GetAll_with_view_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	service := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_username",
		Limit:     200,
	}).Return(nil, fmt.Errorf("some-error")).Once()

	res, err := service.GetAll(context.Background())

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to query the view",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	dbDriver.AssertExpectations(t)
}

func Test_User_StorageMock_GetAll_with_GetMany_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	service := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_username",
		Limit:     200,
	}).Return([]db.ViewRow{
		{ID: "some-id"},
		{ID: "some-id-2"},
	}, nil).Once()

	dbDriver.On("GetMany", []string{"some-id", "some-id-2"}).Return(nil, fmt.Errorf("some-error")).Once()

	res, err := service.GetAll(context.Background())

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get the documents",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	dbDriver.AssertExpectations(t)
}

func Test_Client_StorageMock_FindOneByUsername(t *testing.T) {
	mock := new(StorageMock)

	mock.On("FindOneByUsername", "some-username").Return("some-id", "some-rev", &ValidUser, nil)

	id, rev, res, err := mock.FindOneByUsername(context.Background(), "some-username")

	assert.NoError(t, err)
	assert.Equal(t, "some-id", id)
	assert.Equal(t, "some-rev", rev)
	assert.EqualValues(t, &ValidUser, res)

	mock.AssertExpectations(t)
}

func Test_Client_StorageMock_FindOneByUsername_with_error(t *testing.T) {
	mock := new(StorageMock)

	mock.On("FindOneByUsername", "some-username").Return("", "", nil, errors.New("some-error"))

	id, rev, res, err := mock.FindOneByUsername(context.Background(), "some-username")

	assert.Empty(t, rev)
	assert.Empty(t, id)
	assert.Nil(t, res)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/server/db"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_User_Storage_Set(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Set", "some-user-id", "", &ValidUser).Return("some-rev", nil).Once()

	rev, err := user.Set(context.Background(), "some-user-id", "", &ValidUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, rev)

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_Set_with_driver_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Set", "some-user-id", "some-rev", &ValidUser).Return("", fmt.Errorf("some-error")).Once()

	rev, err := user.Set(context.Background(), "some-user-id", "some-rev", &ValidUser)

	assert.Empty(t, rev)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message": "failed to set the document into the storage",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_Get(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Get", "some-user-id").Return("some-rev", &ValidUser, nil).Once()

	rev, res, err := user.Get(context.Background(), "some-user-id")

	assert.NoError(t, err)
	assert.Equal(t, "some-rev", rev)
	assert.EqualValues(t, ValidUser, *res)

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_Get_not_found(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Get", "some-user-id").Return("", nil, nil).Once()

	rev, res, err := user.Get(context.Background(), "some-user-id")

	assert.NoError(t, err)
	assert.Empty(t, "", rev)
	assert.Nil(t, res)

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_Get_with_driver_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Get", "some-user-id").Return("", nil, errors.New("some-error")).Once()

	rev, res, err := user.Get(context.Background(), "some-user-id")

	assert.Empty(t, "", rev)
	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message": "failed to get the document from the storage",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_Delete(t *testing.T) {
	dbDriver := new(db.DriverMock)
	storage := NewStorage(dbDriver)

	dbDriver.On("Get", "some-id").Return("some-rev", &ValidUser, nil).Once()
	dbDriver.On("Delete", "some-id", "some-rev").Return(nil).Once()

	err := storage.Delete(context.Background(), "some-id")

	assert.NoError(t, err)

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_Delete_with_a_get_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	storage := NewStorage(dbDriver)

	dbDriver.On("Get", "some-id").Return("", nil, errors.New("some-error")).Once()

	err := storage.Delete(context.Background(), "some-id")

	assert.JSONEq(t, `{
		"kind": "internalError",
		"message": "failed to get the document from the storage",
		"reason": {
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_Delete_with_no_document_found(t *testing.T) {
	dbDriver := new(db.DriverMock)
	storage := NewStorage(dbDriver)

	dbDriver.On("Get", "some-id").Return("", nil, nil).Once()

	err := storage.Delete(context.Background(), "some-id")

	assert.NoError(t, err)

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_Delete_with_a_delete_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	storage := NewStorage(dbDriver)

	dbDriver.On("Get", "some-id").Return("some-rev", &ValidUser, nil).Once()
	dbDriver.On("Delete", "some-id", "some-rev").Return(errors.New("some-error")).Once()

	err := storage.Delete(context.Background(), "some-id")

	assert.JSONEq(t, `{
		"kind": "internalError",
		"message": "failed to delete the document from the storage",
		"reason": {
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_FindOneByUsername(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_username",
		Limit:     1,
		Equals:    []interface{}{"some-username"},
	}).Return([]db.ViewRow{
		{ID: "some-id"},
	}, nil).Once()

	dbDriver.On("Get", "some-id").Return("some-rev", &ValidUser, nil).Once()

	id, rev, res, err := user.FindOneByUsername(context.Background(), "some-username")

	assert.Equal(t, "some-id", id)
	assert.Equal(t, "some-rev", rev)
	assert.NoError(t, err)
	assert.EqualValues(t, &ValidUser, res)

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_FindOneByUsername_with_no_user_found(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_username",
		Limit:     1,
		Equals:    []interface{}{"some-username"},
	}).Return(nil, nil).Once()

	id, rev, res, err := user.FindOneByUsername(context.Background(), "some-username")

	assert.Empty(t, id)
	assert.Empty(t, rev)
	assert.NoError(t, err)
	assert.Nil(t, res)

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_FindOneByUsername_with_query_view_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_username",
		Limit:     1,
		Equals:    []interface{}{"some-username"},
	}).Return(nil, fmt.Errorf("some-error")).Once()

	id, rev, res, err := user.FindOneByUsername(context.Background(), "some-username")

	assert.Empty(t, id)
	assert.Empty(t, rev)
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

func Test_User_Storage_FindOneByUsername_with_get_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_username",
		Limit:     1,
		Equals:    []interface{}{"some-username"},
	}).Return([]db.ViewRow{
		{ID: "some-id"},
	}, nil).Once()

	dbDriver.On("Get", "some-id").Return("", nil, fmt.Errorf("some-error")).Once()

	id, rev, res, err := user.FindOneByUsername(context.Background(), "some-username")

	assert.Empty(t, id)
	assert.Empty(t, rev)
	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get the document",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_GetAll(t *testing.T) {
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

func Test_User_Storage_GetAll_empty(t *testing.T) {
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

func Test_User_Storage_GetAll_with_view_error(t *testing.T) {
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

func Test_User_Storage_GetAll_with_GetMany_error(t *testing.T) {
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

func Test_User_Storage_FindTotalUserCount(t *testing.T) {
	dbDriver := new(db.DriverMock)
	service := NewStorage(dbDriver)

	dbDriver.On("GetTotalRow").Return(42, nil).Once()

	nbUsers, err := service.FindTotalUserCount(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 41, nbUsers)

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_FindTotalUserCount_with_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	service := NewStorage(dbDriver)

	dbDriver.On("GetTotalRow").Return(0, fmt.Errorf("some-error")).Once()

	nbUsers, err := service.FindTotalUserCount(context.Background())

	assert.Equal(t, 0, nbUsers)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to query the driver",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	dbDriver.AssertExpectations(t)
}

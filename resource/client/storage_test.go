package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/server/db"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Client_Storage_Set(t *testing.T) {
	dbDriver := new(db.DriverMock)
	storage := NewStorage(dbDriver)

	dbDriver.On("Set", "some-id", "", &ValidClient).Return("some-rev", nil).Once()

	rev, err := storage.Set(context.Background(), "some-id", "", &ValidClient)
	assert.NoError(t, err)
	assert.NotEmpty(t, rev)

	dbDriver.AssertExpectations(t)
}

func Test_Client_Storage_Set_with_driver_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	storage := NewStorage(dbDriver)

	dbDriver.On("Set", "some-id", "some-rev", &ValidClient).Return("", errors.New("some-error")).Once()

	rev, err := storage.Set(context.Background(), "some-id", "some-rev", &ValidClient)
	assert.Empty(t, rev)
	assert.JSONEq(t, `{
		"kind": "internalError",
		"message": "failed to set the document into the storage",
		"reason": {
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	dbDriver.AssertExpectations(t)
}

func Test_Client_Storage_Get(t *testing.T) {
	dbDriver := new(db.DriverMock)
	storage := NewStorage(dbDriver)

	dbDriver.On("Get", "some-id").Return("some-rev", &ValidClient, nil).Once()

	rev, res, err := storage.Get(context.Background(), "some-id")

	assert.NoError(t, err)
	assert.Equal(t, "some-rev", rev)
	assert.EqualValues(t, ValidClient, *res)

	dbDriver.AssertExpectations(t)
}

func Test_Client_Storage_Get_not_found(t *testing.T) {
	dbDriver := new(db.DriverMock)

	dbDriver.On("Get", "some-id").Return("", nil, nil).Once()

	storage := NewStorage(dbDriver)

	rev, res, err := storage.Get(context.Background(), "some-id")

	assert.NoError(t, err)
	assert.Empty(t, "", rev)
	assert.Nil(t, res)

	dbDriver.AssertExpectations(t)
}

func Test_Client_Storage_Get_with_driver_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	storage := NewStorage(dbDriver)

	dbDriver.On("Get", "some-id").Return("", nil, errors.New("some-error")).Once()

	rev, res, err := storage.Get(context.Background(), "some-id")

	assert.Empty(t, "", rev)
	assert.Nil(t, res)
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

func Test_Client_Storage_Delete(t *testing.T) {
	dbDriver := new(db.DriverMock)
	storage := NewStorage(dbDriver)

	dbDriver.On("Get", "some-id").Return("some-rev", &ValidClient, nil).Once()
	dbDriver.On("Delete", "some-id", "some-rev").Return(nil).Once()

	err := storage.Delete(context.Background(), "some-id")

	assert.NoError(t, err)

	dbDriver.AssertExpectations(t)
}

func Test_Client_Storage_Delete_with_a_get_error(t *testing.T) {
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

func Test_Client_Storage_Delete_with_no_document_found(t *testing.T) {
	dbDriver := new(db.DriverMock)
	storage := NewStorage(dbDriver)

	dbDriver.On("Get", "some-id").Return("", nil, nil).Once()

	err := storage.Delete(context.Background(), "some-id")

	assert.NoError(t, err)

	dbDriver.AssertExpectations(t)
}

func Test_Client_Storage_Delete_with_a_delete_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	storage := NewStorage(dbDriver)

	dbDriver.On("Get", "some-id").Return("some-rev", &ValidClient, nil).Once()
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

func Test_Client_Storage_GetAll(t *testing.T) {
	dbDriver := new(db.DriverMock)
	service := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_name",
		Limit:     200,
	}).Return([]db.ViewRow{
		{ID: "some-id"},
		{ID: "some-id-2"},
	}, nil).Once()

	dbDriver.On("GetMany", []string{"some-id", "some-id-2"}).Return(map[string]Client{
		"some-id":   ValidClient,
		"some-id-2": ValidClient,
	}, nil).Once()

	res, err := service.GetAll(context.Background())

	assert.NoError(t, err)
	assert.EqualValues(t, map[string]Client{
		"some-id":   ValidClient,
		"some-id-2": ValidClient,
	}, res)

	dbDriver.AssertExpectations(t)
}

func Test_Client_Storage_GetAll_empty(t *testing.T) {
	dbDriver := new(db.DriverMock)
	service := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_name",
		Limit:     200,
	}).Return([]db.ViewRow{}, nil).Once()

	res, err := service.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, res)

	dbDriver.AssertExpectations(t)
}

func Test_Client_Storage_GetAll_with_view_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	service := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_name",
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

func Test_Client_Storage_GetAll_with_GetMany_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	service := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_name",
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

func Test_User_Storage_FindOneByName(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_name",
		Limit:     1,
		Equals:    []interface{}{"some-username"},
	}).Return([]db.ViewRow{
		{ID: "some-id"},
	}, nil).Once()

	dbDriver.On("Get", "some-id").Return("some-rev", &ValidClient, nil).Once()

	id, rev, res, err := user.FindOneByName(context.Background(), "some-username")

	assert.Equal(t, "some-id", id)
	assert.Equal(t, "some-rev", rev)
	assert.NoError(t, err)
	assert.EqualValues(t, &ValidClient, res)

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_FindOneByName_with_no_user_found(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_name",
		Limit:     1,
		Equals:    []interface{}{"some-username"},
	}).Return(nil, nil).Once()

	id, rev, res, err := user.FindOneByName(context.Background(), "some-username")

	assert.Empty(t, id)
	assert.Empty(t, rev)
	assert.NoError(t, err)
	assert.Nil(t, res)

	dbDriver.AssertExpectations(t)
}

func Test_User_Storage_FindOneByName_with_query_view_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_name",
		Limit:     1,
		Equals:    []interface{}{"some-username"},
	}).Return(nil, fmt.Errorf("some-error")).Once()

	id, rev, res, err := user.FindOneByName(context.Background(), "some-username")

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

func Test_User_Storage_FindOneByName_with_get_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_name",
		Limit:     1,
		Equals:    []interface{}{"some-username"},
	}).Return([]db.ViewRow{
		{ID: "some-id"},
	}, nil).Once()

	dbDriver.On("Get", "some-id").Return("", nil, fmt.Errorf("some-error")).Once()

	id, rev, res, err := user.FindOneByName(context.Background(), "some-username")

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

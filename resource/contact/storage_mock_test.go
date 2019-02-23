package contact

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/server/db"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Contact_StorageMock_Set(t *testing.T) {
	mock := new(StorageMock)

	mock.On("Set", "some-id", "some-rev", &ValidContact).Return("some-rev-2", nil)

	rev, err := mock.Set(context.Background(), "some-id", "some-rev", &ValidContact)

	assert.NoError(t, err)
	assert.Equal(t, "some-rev-2", rev)

	mock.AssertExpectations(t)
}

func Test_Contact_StorageMock_Get(t *testing.T) {
	mock := new(StorageMock)

	mock.On("Get", "some-id").Return("some-rev", &ValidContact, nil)

	rev, res, err := mock.Get(context.Background(), "some-id")

	assert.NoError(t, err)
	assert.Equal(t, "some-rev", rev)
	assert.EqualValues(t, &ValidContact, res)

	mock.AssertExpectations(t)
}

func Test_Contact_StorageMock_Delete(t *testing.T) {
	mock := new(StorageMock)

	mock.On("Delete", "some-id").Return(nil)

	err := mock.Delete(context.Background(), "some-id")

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

func Test_Contact_StorageMock_Get_with_error(t *testing.T) {
	mock := new(StorageMock)

	mock.On("Get", "some-id").Return("", nil, errors.New("some-error"))

	rev, res, err := mock.Get(context.Background(), "some-id")
	assert.EqualError(t, err, "some-error")
	assert.Empty(t, rev)
	assert.Nil(t, res)

	mock.AssertExpectations(t)
}

func Test_Contact_StorageMock_GetAll(t *testing.T) {
	dbDriver := new(db.DriverMock)
	service := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_name",
		Limit:     200,
	}).Return([]db.ViewRow{
		{ID: "some-id"},
		{ID: "some-id-2"},
	}, nil).Once()

	dbDriver.On("GetMany", []string{"some-id", "some-id-2"}).Return(map[string]Contact{
		"some-id":   ValidContact,
		"some-id-2": ValidContact,
	}, nil).Once()

	res, err := service.GetAll(context.Background())

	assert.NoError(t, err)
	assert.EqualValues(t, map[string]Contact{
		"some-id":   ValidContact,
		"some-id-2": ValidContact,
	}, res)

	dbDriver.AssertExpectations(t)
}

func Test_Contact_StorageMock_GetAll_empty(t *testing.T) {
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

func Test_Contact_StorageMock_GetAll_with_view_error(t *testing.T) {
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

func Test_Contact_StorageMock_GetAll_with_GetMany_error(t *testing.T) {
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

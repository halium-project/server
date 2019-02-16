package accesstoken

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/server/db"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_AccessToken_Storage_Set(t *testing.T) {
	dbDriver := new(db.DriverMock)
	accessToken := NewStorage(dbDriver)

	dbDriver.On("Set", "some-authorization-code", "", &ValidAccessToken).Return("some-rev", nil).Once()

	rev, err := accessToken.Set(context.Background(), "some-authorization-code", "", &ValidAccessToken)

	assert.NoError(t, err)
	assert.NotEmpty(t, rev)

	dbDriver.AssertExpectations(t)
}

func Test_AccessToken_Storage_Set_with_driver_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	accessToken := NewStorage(dbDriver)

	dbDriver.On("Set", "some-authorization-code", "some-rev", &ValidAccessToken).Return("", fmt.Errorf("some-error")).Once()

	rev, err := accessToken.Set(context.Background(), "some-authorization-code", "some-rev", &ValidAccessToken)

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

func Test_AccessToken_Storage_Get(t *testing.T) {
	dbDriver := new(db.DriverMock)
	accessToken := NewStorage(dbDriver)

	dbDriver.On("Get", "some-authorization-code").Return("some-rev", &ValidAccessToken, nil).Once()

	rev, res, err := accessToken.Get(context.Background(), "some-authorization-code")

	assert.NoError(t, err)
	assert.Equal(t, "some-rev", rev)
	assert.EqualValues(t, ValidAccessToken, *res)

	dbDriver.AssertExpectations(t)
}

func Test_AccessToken_Storage_Get_not_found(t *testing.T) {
	dbDriver := new(db.DriverMock)
	accessToken := NewStorage(dbDriver)

	dbDriver.On("Get", "some-authorization-code").Return("", nil, nil).Once()

	rev, res, err := accessToken.Get(context.Background(), "some-authorization-code")

	assert.NoError(t, err)
	assert.Empty(t, "", rev)
	assert.Nil(t, res)

	dbDriver.AssertExpectations(t)
}

func Test_AccessToken_Storage_Get_with_driver_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	accessToken := NewStorage(dbDriver)

	dbDriver.On("Get", "some-authorization-code").Return("", nil, errors.New("some-error")).Once()

	rev, res, err := accessToken.Get(context.Background(), "some-authorization-code")

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

func Test_AccessToken_Storage_Delete(t *testing.T) {
	dbDriver := new(db.DriverMock)
	accessToken := NewStorage(dbDriver)

	dbDriver.On("Delete", "some-authorization-code", "some-rev").Return(nil).Once()

	err := accessToken.Delete(context.Background(), "some-authorization-code", "some-rev")

	assert.NoError(t, err)

	dbDriver.AssertExpectations(t)
}

func Test_AccessToken_Storage_Delete_with_driver_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	accessToken := NewStorage(dbDriver)

	dbDriver.On("Delete", "some-authorization-code", "some-rev").Return(errors.New("some-error")).Once()

	err := accessToken.Delete(context.Background(), "some-authorization-code", "some-rev")

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message": "failed to delete the document from the storage",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	dbDriver.AssertExpectations(t)
}

func Test_AccessToken_Storage_FindOneByRefreshToken(t *testing.T) {
	dbDriver := new(db.DriverMock)
	accessToken := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_refresh_token",
		Limit:     1,
		Equals:    []interface{}{"some-refreshToken"},
	}).Return([]db.ViewRow{
		{ID: "some-id"},
	}, nil).Once()

	dbDriver.On("Get", "some-id").Return("some-rev", &ValidAccessToken, nil).Once()

	id, rev, res, err := accessToken.FindOneByRefreshToken(context.Background(), "some-refreshToken")

	assert.Equal(t, "some-id", id)
	assert.Equal(t, "some-rev", rev)
	assert.NoError(t, err)
	assert.EqualValues(t, &ValidAccessToken, res)

	dbDriver.AssertExpectations(t)
}

func Test_AccessToken_Storage_FindOneByRefreshToken_with_no_accessToken_found(t *testing.T) {
	dbDriver := new(db.DriverMock)
	accessToken := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_refresh_token",
		Limit:     1,
		Equals:    []interface{}{"some-refreshToken"},
	}).Return(nil, nil).Once()

	id, rev, res, err := accessToken.FindOneByRefreshToken(context.Background(), "some-refreshToken")

	assert.Empty(t, id)
	assert.Empty(t, rev)
	assert.NoError(t, err)
	assert.Nil(t, res)

	dbDriver.AssertExpectations(t)
}

func Test_AccessToken_Storage_FindOneByRefreshToken_with_query_view_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	accessToken := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_refresh_token",
		Limit:     1,
		Equals:    []interface{}{"some-refreshToken"},
	}).Return(nil, fmt.Errorf("some-error")).Once()

	id, rev, res, err := accessToken.FindOneByRefreshToken(context.Background(), "some-refreshToken")

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

func Test_AccessToken_Storage_FindOneByRefreshToken_with_get_error(t *testing.T) {
	dbDriver := new(db.DriverMock)
	accessToken := NewStorage(dbDriver)

	dbDriver.On("ExecuteViewQuery", &db.Query{
		IndexName: "by_refresh_token",
		Limit:     1,
		Equals:    []interface{}{"some-refreshToken"},
	}).Return([]db.ViewRow{
		{ID: "some-id"},
	}, nil).Once()

	dbDriver.On("Get", "some-id").Return("", nil, fmt.Errorf("some-error")).Once()

	id, rev, res, err := accessToken.FindOneByRefreshToken(context.Background(), "some-refreshToken")

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

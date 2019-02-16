package authorizationcode

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/server/db"
	"github.com/halium-project/server/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_AuthorizationCode_Storage_Set(t *testing.T) {
	util.TestIs(t, util.Unit)

	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Set", "some-authorization-code", "", &ValidAuthorizationCode).Return("some-rev", nil).Once()

	rev, err := user.Set(context.Background(), "some-authorization-code", "", &ValidAuthorizationCode)

	assert.NoError(t, err)
	assert.NotEmpty(t, rev)

	dbDriver.AssertExpectations(t)
}

func Test_AuthorizationCode_Storage_Set_with_driver_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Set", "some-authorization-code", "some-rev", &ValidAuthorizationCode).Return("", fmt.Errorf("some-error")).Once()

	rev, err := user.Set(context.Background(), "some-authorization-code", "some-rev", &ValidAuthorizationCode)

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

func Test_AuthorizationCode_Storage_Get(t *testing.T) {
	util.TestIs(t, util.Unit)

	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Get", "some-authorization-code").Return("some-rev", &ValidAuthorizationCode, nil).Once()

	rev, res, err := user.Get(context.Background(), "some-authorization-code")

	assert.NoError(t, err)
	assert.Equal(t, "some-rev", rev)
	assert.EqualValues(t, ValidAuthorizationCode, *res)

	dbDriver.AssertExpectations(t)
}

func Test_AuthorizationCode_Storage_Get_not_found(t *testing.T) {
	util.TestIs(t, util.Unit)

	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Get", "some-authorization-code").Return("", nil, nil).Once()

	rev, res, err := user.Get(context.Background(), "some-authorization-code")

	assert.NoError(t, err)
	assert.Empty(t, "", rev)
	assert.Nil(t, res)

	dbDriver.AssertExpectations(t)
}

func Test_AuthorizationCode_Storage_Get_with_driver_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Get", "some-authorization-code").Return("", nil, errors.New("some-error")).Once()

	rev, res, err := user.Get(context.Background(), "some-authorization-code")

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

func Test_AuthorizationCode_Storage_Delete(t *testing.T) {
	util.TestIs(t, util.Unit)

	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Delete", "some-authorization-code", "some-rev").Return(nil).Once()

	err := user.Delete(context.Background(), "some-authorization-code", "some-rev")

	assert.NoError(t, err)

	dbDriver.AssertExpectations(t)
}

func Test_AuthorizationCode_Storage_Delete_with_driver_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	dbDriver := new(db.DriverMock)
	user := NewStorage(dbDriver)

	dbDriver.On("Delete", "some-authorization-code", "some-rev").Return(errors.New("some-error")).Once()

	err := user.Delete(context.Background(), "some-authorization-code", "some-rev")

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

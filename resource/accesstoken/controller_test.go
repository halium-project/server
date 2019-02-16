package accesstoken

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/server/util"
	"github.com/halium-project/server/util/password"
	"github.com/halium-project/server/util/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_AccessToken_Controller_Create(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Set", "some-access-token", "", &ValidAccessToken).Return("some-rev", nil).Once()

	err := controller.Create(context.Background(), &CreateCmd{
		ClientID:     ValidAccessToken.ClientID,
		AccessToken:  ValidAccessToken.AccessToken,
		RefreshToken: ValidAccessToken.RefreshToken,
		ExpiresIn:    ValidAccessToken.ExpiresIn,
		Scopes:       ValidAccessToken.Scopes,
	})

	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_Create_with_validation_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	err := controller.Create(context.Background(), &CreateCmd{
		ClientID:     "invalid-format-id",
		AccessToken:  ValidAccessToken.AccessToken,
		RefreshToken: ValidAccessToken.RefreshToken,
		ExpiresIn:    ValidAccessToken.ExpiresIn,
		Scopes:       ValidAccessToken.Scopes,
	})

	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"clientID":"INVALID_FORMAT"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_Create_storage_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Set", "some-access-token", "", &ValidAccessToken).Return("", fmt.Errorf("some-error")).Once()

	err := controller.Create(context.Background(), &CreateCmd{
		ClientID:     ValidAccessToken.ClientID,
		AccessToken:  ValidAccessToken.AccessToken,
		RefreshToken: ValidAccessToken.RefreshToken,
		ExpiresIn:    ValidAccessToken.ExpiresIn,
		Scopes:       ValidAccessToken.Scopes,
	})

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to save the accessToken",
		"reason": {
			"kind":"internalError",
			"message":"some-error"
		}

	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_Get(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "some-access-token").Return("some-rev", &ValidAccessToken, nil).Once()

	res, err := controller.Get(context.Background(), &GetCmd{
		AccessToken: "some-access-token",
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidAccessToken, res)

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_Get_with_validationError(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	res, err := controller.Get(context.Background(), &GetCmd{
		AccessToken: "1", // too short
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors":{
			"accessToken":"TOO_SHORT"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_Get_driver_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "some-access-token").Return("", nil, fmt.Errorf("some-error")).Once()

	res, err := controller.Get(context.Background(), &GetCmd{
		AccessToken: "some-access-token",
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get the accessToken",
		"reason":{
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_Delete(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "some-access-token").Return("some-rev", &ValidAccessToken, nil).Once()
	storageMock.On("Delete", "some-access-token", "some-rev").Return(nil).Once()

	err := controller.Delete(context.Background(), &DeleteCmd{
		AccessToken: "some-access-token",
	})

	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_Delete_with_validationError(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	err := controller.Delete(context.Background(), &DeleteCmd{
		AccessToken: "1", // too short
	})

	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors":{
			"accessToken":"TOO_SHORT"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_Delete_driver_get_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "some-access-token").Return("some-rev", &ValidAccessToken, fmt.Errorf("some-error")).Once()

	err := controller.Delete(context.Background(), &DeleteCmd{
		AccessToken: "some-access-token",
	})

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get the accessToken",
		"reason":{
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_Delete_driver_delete_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "some-access-token").Return("some-rev", &ValidAccessToken, nil).Once()
	storageMock.On("Delete", "some-access-token", "some-rev").Return(fmt.Errorf("some-error")).Once()

	err := controller.Delete(context.Background(), &DeleteCmd{
		AccessToken: "some-access-token",
	})

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to delete the accessToken",
		"reason":{
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_FindOnebyRefreshToken(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByRefreshToken", "some-refresh-token").Return("some-id", "some-rev", &ValidAccessToken, nil).Once()

	id, res, err := controller.FindOneByRefreshToken(context.Background(), &FindOneByRefreshTokenCmd{
		RefreshToken: "some-refresh-token",
	})

	assert.NoError(t, err)
	assert.Equal(t, "some-id", id)
	assert.EqualValues(t, &ValidAccessToken, res)

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_FindOnebyRefreshToken_with_validation_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	id, res, err := controller.FindOneByRefreshToken(context.Background(), &FindOneByRefreshTokenCmd{
		RefreshToken: "a", // input too short
	})

	assert.Empty(t, id)
	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"refreshToken": "TOO_SHORT"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AccessToken_Controller_FindOnebyRefreshToken_with_storage_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByRefreshToken", "some-refresh-token").Return("", "", nil, fmt.Errorf("some-error")).Once()

	id, res, err := controller.FindOneByRefreshToken(context.Background(), &FindOneByRefreshTokenCmd{
		RefreshToken: "some-refresh-token",
	})

	assert.Empty(t, id)
	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get the accessToken",
		"reason":{
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

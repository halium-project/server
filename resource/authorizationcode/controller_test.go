package authorizationcode

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/server/util/password"
	"github.com/halium-project/server/util/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_AuthorizationCode_Controller_Create(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Set", "some-authorization-code", "", &ValidAuthorizationCode).Return("some-rev", nil).Once()

	err := controller.Create(context.Background(), &CreateCmd{
		ClientID:            ValidAuthorizationCode.ClientID,
		Code:                "some-authorization-code",
		ExpiresIn:           ValidAuthorizationCode.ExpiresIn,
		Scopes:              ValidAuthorizationCode.Scopes,
		RedirectURI:         ValidAuthorizationCode.RedirectURI,
		State:               ValidAuthorizationCode.State,
		CodeChallenge:       ValidAuthorizationCode.CodeChallenge,
		CodeChallengeMethod: ValidAuthorizationCode.CodeChallengeMethod,
	})

	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AuthorizationCode_Controller_Create_with_validation_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	err := controller.Create(context.Background(), &CreateCmd{
		ClientID:            ValidAuthorizationCode.ClientID,
		Code:                "some-authorization-code",
		ExpiresIn:           -1, // should be positif
		Scopes:              ValidAuthorizationCode.Scopes,
		RedirectURI:         ValidAuthorizationCode.RedirectURI,
		State:               ValidAuthorizationCode.State,
		CodeChallenge:       ValidAuthorizationCode.CodeChallenge,
		CodeChallengeMethod: ValidAuthorizationCode.CodeChallengeMethod,
	})

	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"expiresIn":"UNEXPECTED_VALUE"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AuthorizationCode_Controller_Create_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Set", "some-authorization-code", "", &ValidAuthorizationCode).Return("", fmt.Errorf("some-error")).Once()

	err := controller.Create(context.Background(), &CreateCmd{
		ClientID:            ValidAuthorizationCode.ClientID,
		Code:                "some-authorization-code",
		ExpiresIn:           ValidAuthorizationCode.ExpiresIn,
		Scopes:              ValidAuthorizationCode.Scopes,
		RedirectURI:         ValidAuthorizationCode.RedirectURI,
		State:               ValidAuthorizationCode.State,
		CodeChallenge:       ValidAuthorizationCode.CodeChallenge,
		CodeChallengeMethod: ValidAuthorizationCode.CodeChallengeMethod,
	})

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to save the user",
		"reason": {
			"kind":"internalError",
			"message":"some-error"
		}

	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AuthorizationCode_Controller_Get(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "some-authorization-code").Return("some-rev", &ValidAuthorizationCode, nil).Once()

	res, err := controller.Get(context.Background(), &GetCmd{
		Code: "some-authorization-code",
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidAuthorizationCode, res)

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AuthorizationCode_Controller_Get_with_validationError(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	res, err := controller.Get(context.Background(), &GetCmd{
		Code: "1", // too short
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors":{
			"code":"TOO_SHORT"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AuthorizationCode_Controller_Get_driver_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "some-authorization-code").Return("", nil, fmt.Errorf("some-error")).Once()

	res, err := controller.Get(context.Background(), &GetCmd{
		Code: "some-authorization-code",
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get the user",
		"reason":{
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AuthorizationCode_Controller_Delete(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "some-authorization-code").Return("some-rev", &ValidAuthorizationCode, nil).Once()
	storageMock.On("Delete", "some-authorization-code", "some-rev").Return(nil).Once()

	err := controller.Delete(context.Background(), &DeleteCmd{
		Code: "some-authorization-code",
	})

	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AuthorizationCode_Controller_Delete_with_validationError(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	err := controller.Delete(context.Background(), &DeleteCmd{
		Code: "1", // too short
	})

	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors":{
			"code":"TOO_SHORT"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AuthorizationCode_Controller_Delete_driver_get_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "some-authorization-code").Return("some-rev", &ValidAuthorizationCode, fmt.Errorf("some-error")).Once()

	err := controller.Delete(context.Background(), &DeleteCmd{
		Code: "some-authorization-code",
	})

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get the user",
		"reason":{
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_AuthorizationCode_Controller_Delete_driver_delete_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "some-authorization-code").Return("some-rev", &ValidAuthorizationCode, nil).Once()
	storageMock.On("Delete", "some-authorization-code", "some-rev").Return(fmt.Errorf("some-error")).Once()

	err := controller.Delete(context.Background(), &DeleteCmd{
		Code: "some-authorization-code",
	})

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to delete the user",
		"reason":{
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

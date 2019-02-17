package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/go-server-utils/password"
	"github.com/halium-project/go-server-utils/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// nolint
// Hardcoded secret used only for tests
const validSecret = "2558539b-e119-408d-a31b-a1d6cf1b60aa"

func Test_Client_Controller_Create(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByName", ValidClient.Name).Return("", "", nil, nil).Once()
	uuidMock.On("New").Return(validSecret).Once()
	passwordMock.On("Hash", validSecret).Return("some-hashed-secret", nil).Once()
	storageMock.On("Set", ValidID, "", &ValidClient).Return("some-rev", nil).Once()

	id, secret, err := handler.Create(context.Background(), &CreateCmd{
		Name:          ValidClient.Name,
		RedirectURIs:  ValidClient.RedirectURIs,
		GrantTypes:    ValidClient.GrantTypes,
		ResponseTypes: ValidClient.ResponseTypes,
		Scopes:        ValidClient.Scopes,
		Public:        ValidClient.Public,
	})

	assert.NoError(t, err)
	assert.Equal(t, ValidID, id)
	assert.Equal(t, validSecret, secret)

	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Client_Controller_Create_with_validationError(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, passwordMock, storageMock)

	id, secret, err := handler.Create(context.Background(), &CreateCmd{
		Name:          "i", // name too short
		RedirectURIs:  ValidClient.RedirectURIs,
		GrantTypes:    ValidClient.GrantTypes,
		ResponseTypes: ValidClient.ResponseTypes,
		Scopes:        ValidClient.Scopes,
		Public:        ValidClient.Public,
	})

	assert.Empty(t, id)
	assert.Empty(t, secret)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"name": "TOO_SHORT"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Client_Controller_Create_with_storage_get_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByName", ValidClient.Name).Return("", "", nil, fmt.Errorf("some-error")).Once()

	id, secret, err := handler.Create(context.Background(), &CreateCmd{
		Name:          ValidClient.Name,
		RedirectURIs:  ValidClient.RedirectURIs,
		GrantTypes:    ValidClient.GrantTypes,
		ResponseTypes: ValidClient.ResponseTypes,
		Scopes:        ValidClient.Scopes,
		Public:        ValidClient.Public,
	})

	assert.Empty(t, id)
	assert.Empty(t, secret)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to check if the name is already taken",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Client_Controller_Create_with_name_already_taken(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByName", ValidClient.Name).Return("some-id", "some-rev", &ValidClient, nil).Once()

	id, secret, err := handler.Create(context.Background(), &CreateCmd{
		Name:          ValidClient.Name,
		RedirectURIs:  ValidClient.RedirectURIs,
		GrantTypes:    ValidClient.GrantTypes,
		ResponseTypes: ValidClient.ResponseTypes,
		Scopes:        ValidClient.Scopes,
		Public:        ValidClient.Public,
	})

	assert.Empty(t, id)
	assert.Empty(t, secret)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"name": "ALREADY_USED"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Client_Controller_Create_with_password_hash_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByName", ValidClient.Name).Return("", "", nil, nil).Once()
	uuidMock.On("New").Return(validSecret).Once() // one time for the id / one time for the secret
	passwordMock.On("Hash", validSecret).Return("", fmt.Errorf("some-error")).Once()

	id, secret, err := handler.Create(context.Background(), &CreateCmd{
		Name:          ValidClient.Name,
		RedirectURIs:  ValidClient.RedirectURIs,
		GrantTypes:    ValidClient.GrantTypes,
		ResponseTypes: ValidClient.ResponseTypes,
		Scopes:        ValidClient.Scopes,
		Public:        ValidClient.Public,
	})

	assert.Empty(t, id)
	assert.Empty(t, secret)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to hash password",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Client_Controller_Create_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByName", ValidClient.Name).Return("", "", nil, nil).Once()
	uuidMock.On("New").Return(validSecret).Once() // one time for the id / one time for the secret
	passwordMock.On("Hash", validSecret).Return("some-hashed-secret", nil).Once()
	storageMock.On("Set", ValidID, "", &ValidClient).Return("", fmt.Errorf("some-error")).Once()

	id, secret, err := handler.Create(context.Background(), &CreateCmd{
		Name:          ValidClient.Name,
		RedirectURIs:  ValidClient.RedirectURIs,
		GrantTypes:    ValidClient.GrantTypes,
		ResponseTypes: ValidClient.ResponseTypes,
		Scopes:        ValidClient.Scopes,
		Public:        ValidClient.Public,
	})

	assert.Empty(t, id)
	assert.Empty(t, secret)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to save a client",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Client_Controller_Get(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", validSecret).Return("some-rev", &ValidClient, nil).Once()

	res, err := handler.Get(context.Background(), &GetCmd{
		ClientID: validSecret,
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidClient, res)

	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Client_Controller_Get_with_validationError(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, passwordMock, storageMock)

	res, err := handler.Get(context.Background(), &GetCmd{
		ClientID: "i", // too short
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind": "validationError",
		"errors": {
			"clientId": "TOO_SHORT"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Client_Controller_Get_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", validSecret).Return("", nil, fmt.Errorf("some-error")).Once()

	res, err := handler.Get(context.Background(), &GetCmd{
		ClientID: validSecret,
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get a client",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Client_Controller_Get_with_resource_notFound(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", validSecret).Return("", nil, nil).Once()

	res, err := handler.Get(context.Background(), &GetCmd{
		ClientID: validSecret,
	})

	assert.NoError(t, err)
	assert.Nil(t, res)

	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Client_Controller_GetAll(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("GetAll").Return(map[string]Client{
		"some-rev":   ValidClient,
		"some-rev-2": ValidClient,
	}, nil).Once()

	res, err := controller.GetAll(context.Background(), &GetAllCmd{})

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.EqualValues(t, map[string]Client{
		"some-rev":   ValidClient,
		"some-rev-2": ValidClient,
	}, res)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_Client_Controller_GetAll_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("GetAll").Return(nil, errors.New("some-error")).Once()

	res, err := controller.GetAll(context.Background(), &GetAllCmd{})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get all clients",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

package contact

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/go-server-utils/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Contact_Controller_Create(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("FindOneByName", ValidContact.Name).Return("", "", nil, nil).Once()
	uuidMock.On("New").Return(ValidContactID).Once()
	storageMock.On("Set", ValidContactID, "", &ValidContact).Return("some-rev", nil).Once()

	id, err := handler.Create(context.Background(), &CreateCmd{
		Name: ValidContact.Name,
	})

	assert.NoError(t, err)
	assert.Equal(t, ValidContactID, id)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_Create_with_validationError(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	id, err := handler.Create(context.Background(), &CreateCmd{
		Name: "i", // name too short
	})

	assert.Empty(t, id)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"name": "TOO_SHORT"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_Create_with_storage_get_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("FindOneByName", ValidContact.Name).Return("", "", nil, fmt.Errorf("some-error")).Once()

	id, err := handler.Create(context.Background(), &CreateCmd{
		Name: ValidContact.Name,
	})

	assert.Empty(t, id)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to check if the name is already taken",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_Create_with_name_already_taken(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("FindOneByName", ValidContact.Name).Return("some-id", "some-rev", &ValidContact, nil).Once()

	id, err := handler.Create(context.Background(), &CreateCmd{
		Name: ValidContact.Name,
	})

	assert.Empty(t, id)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"name": "ALREADY_USED"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_Create_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("FindOneByName", ValidContact.Name).Return("", "", nil, nil).Once()
	uuidMock.On("New").Return(ValidContactID).Once() // one time for the id / one time for the secret
	storageMock.On("Set", ValidContactID, "", &ValidContact).Return("", fmt.Errorf("some-error")).Once()

	id, err := handler.Create(context.Background(), &CreateCmd{
		Name: ValidContact.Name,
	})

	assert.Empty(t, id)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to save a contact",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_Get(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("Get", ValidContactID).Return("some-rev", &ValidContact, nil).Once()

	res, err := handler.Get(context.Background(), &GetCmd{
		ContactID: ValidContactID,
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidContact, res)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_Get_with_validationError(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	res, err := handler.Get(context.Background(), &GetCmd{
		ContactID: "i", // too short
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind": "validationError",
		"errors": {
			"contactID": "TOO_SHORT"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_Get_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("Get", ValidContactID).Return("", nil, fmt.Errorf("some-error")).Once()

	res, err := handler.Get(context.Background(), &GetCmd{
		ContactID: ValidContactID,
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get a contact",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_Get_with_resource_notFound(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("Get", ValidContactID).Return("", nil, nil).Once()

	res, err := handler.Get(context.Background(), &GetCmd{
		ContactID: ValidContactID,
	})

	assert.NoError(t, err)
	assert.Nil(t, res)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_GetAll(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, storageMock)

	storageMock.On("GetAll").Return(map[string]Contact{
		"some-rev":   ValidContact,
		"some-rev-2": ValidContact,
	}, nil).Once()

	res, err := controller.GetAll(context.Background(), &GetAllCmd{})

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.EqualValues(t, map[string]Contact{
		"some-rev":   ValidContact,
		"some-rev-2": ValidContact,
	}, res)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_GetAll_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, storageMock)

	storageMock.On("GetAll").Return(nil, errors.New("some-error")).Once()

	res, err := controller.GetAll(context.Background(), &GetAllCmd{})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get all contacts",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_Delete(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, storageMock)

	storageMock.On("Delete", "some-id").Return(nil).Once()

	err := controller.Delete(context.Background(), &DeleteCmd{ContactID: "some-id"})

	assert.NoError(t, err)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_Delete_with_a_validation_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, storageMock)

	err := controller.Delete(context.Background(), &DeleteCmd{ContactID: ""})

	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"contactID": "MISSING_FIELD"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Contact_Controller_Delete_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, storageMock)

	storageMock.On("Delete", "some-id").Return(errors.New("some-error")).Once()

	err := controller.Delete(context.Background(), &DeleteCmd{ContactID: "some-id"})

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to delete the contact",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

package todo

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/go-server-utils/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Todo_Controller_Create(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("FindOneByTitle", ValidTodo.Title).Return("", "", nil, nil).Once()
	uuidMock.On("New").Return(ValidTodoID).Once()
	storageMock.On("Set", ValidTodoID, "", &ValidTodo).Return("some-rev", nil).Once()

	id, err := handler.Create(context.Background(), &CreateCmd{
		Title: ValidTodo.Title,
	})

	assert.NoError(t, err)
	assert.Equal(t, ValidTodoID, id)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_Create_with_validationError(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	id, err := handler.Create(context.Background(), &CreateCmd{
		Title: "i", // title too short
	})

	assert.Empty(t, id)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"title": "TOO_SHORT"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_Create_with_storage_get_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("FindOneByTitle", ValidTodo.Title).Return("", "", nil, fmt.Errorf("some-error")).Once()

	id, err := handler.Create(context.Background(), &CreateCmd{
		Title: ValidTodo.Title,
	})

	assert.Empty(t, id)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to check if the title is already taken",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_Create_with_title_already_taken(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("FindOneByTitle", ValidTodo.Title).Return("some-id", "some-rev", &ValidTodo, nil).Once()

	id, err := handler.Create(context.Background(), &CreateCmd{
		Title: ValidTodo.Title,
	})

	assert.Empty(t, id)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"title": "ALREADY_USED"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_Create_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("FindOneByTitle", ValidTodo.Title).Return("", "", nil, nil).Once()
	uuidMock.On("New").Return(ValidTodoID).Once() // one time for the id / one time for the secret
	storageMock.On("Set", ValidTodoID, "", &ValidTodo).Return("", fmt.Errorf("some-error")).Once()

	id, err := handler.Create(context.Background(), &CreateCmd{
		Title: ValidTodo.Title,
	})

	assert.Empty(t, id)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to save a todo",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_Get(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("Get", ValidTodoID).Return("some-rev", &ValidTodo, nil).Once()

	res, err := handler.Get(context.Background(), &GetCmd{
		TodoID: ValidTodoID,
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidTodo, res)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_Get_with_validationError(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	res, err := handler.Get(context.Background(), &GetCmd{
		TodoID: "i", // too short
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind": "validationError",
		"errors": {
			"todoID": "TOO_SHORT"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_Get_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("Get", ValidTodoID).Return("", nil, fmt.Errorf("some-error")).Once()

	res, err := handler.Get(context.Background(), &GetCmd{
		TodoID: ValidTodoID,
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get a todo",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_Get_with_resource_notFound(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	handler := NewController(uuidMock, storageMock)

	storageMock.On("Get", ValidTodoID).Return("", nil, nil).Once()

	res, err := handler.Get(context.Background(), &GetCmd{
		TodoID: ValidTodoID,
	})

	assert.NoError(t, err)
	assert.Nil(t, res)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_GetAll(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, storageMock)

	storageMock.On("GetAll").Return(map[string]Todo{
		"some-rev":   ValidTodo,
		"some-rev-2": ValidTodo,
	}, nil).Once()

	res, err := controller.GetAll(context.Background(), &GetAllCmd{})

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.EqualValues(t, map[string]Todo{
		"some-rev":   ValidTodo,
		"some-rev-2": ValidTodo,
	}, res)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_GetAll_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, storageMock)

	storageMock.On("GetAll").Return(nil, errors.New("some-error")).Once()

	res, err := controller.GetAll(context.Background(), &GetAllCmd{})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get all todos",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_Delete(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, storageMock)

	storageMock.On("Delete", "some-id").Return(nil).Once()

	err := controller.Delete(context.Background(), &DeleteCmd{TodoID: "some-id"})

	assert.NoError(t, err)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_Delete_with_a_validation_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, storageMock)

	err := controller.Delete(context.Background(), &DeleteCmd{TodoID: ""})

	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"todoID": "MISSING_FIELD"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func Test_Todo_Controller_Delete_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, storageMock)

	storageMock.On("Delete", "some-id").Return(errors.New("some-error")).Once()

	err := controller.Delete(context.Background(), &DeleteCmd{TodoID: "some-id"})

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to delete the todo",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

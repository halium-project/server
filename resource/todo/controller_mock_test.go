package todo

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Todo_ControllerMock_Create(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Create", &CreateCmd{
		Title: ValidTodo.Title,
	}).Return("", fmt.Errorf("some-error")).Once()

	id, err := mock.Create(context.Background(), &CreateCmd{
		Title: ValidTodo.Title,
	})

	assert.Empty(t, id)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_Todo_ControllerMock_Get(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		TodoID: "some-todo-id",
	}).Return(&ValidTodo, nil).Once()

	todo, err := mock.Get(context.Background(), &GetCmd{
		TodoID: "some-todo-id",
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidTodo, todo)

	mock.AssertExpectations(t)
}

func Test_Todo_ControllerMock_Get_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		TodoID: "some-todo-id",
	}).Return(nil, fmt.Errorf("some-error")).Once()

	todo, err := mock.Get(context.Background(), &GetCmd{
		TodoID: "some-todo-id",
	})

	assert.Nil(t, todo)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_Todo_ControllerMock_GetAll(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("GetAll", &GetAllCmd{}).Return(map[string]Todo{
		"some-id": ValidTodo,
	}, nil).Once()

	res, err := mock.GetAll(context.Background(), &GetAllCmd{})

	assert.NoError(t, err)
	assert.EqualValues(t, map[string]Todo{
		"some-id": ValidTodo,
	}, res)

	mock.AssertExpectations(t)
}

func Test_Todo_ControllerMock_GetAll_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("GetAll", &GetAllCmd{}).Return(nil, fmt.Errorf("some-error")).Once()

	todo, err := mock.GetAll(context.Background(), &GetAllCmd{})

	assert.Nil(t, todo)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_Todo_ControllerMock_Delete(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Delete", &DeleteCmd{TodoID: "some-id"}).Return(nil).Once()

	err := mock.Delete(context.Background(), &DeleteCmd{TodoID: "some-id"})

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

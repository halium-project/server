package contact

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Contact_ControllerMock_Create(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Create", &CreateCmd{
		Name: ValidContact.Name,
	}).Return("", fmt.Errorf("some-error")).Once()

	id, err := mock.Create(context.Background(), &CreateCmd{
		Name: ValidContact.Name,
	})

	assert.Empty(t, id)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_Contact_ControllerMock_Get(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		ContactID: "some-contact-id",
	}).Return(&ValidContact, nil).Once()

	contact, err := mock.Get(context.Background(), &GetCmd{
		ContactID: "some-contact-id",
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidContact, contact)

	mock.AssertExpectations(t)
}

func Test_Contact_ControllerMock_Get_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		ContactID: "some-contact-id",
	}).Return(nil, fmt.Errorf("some-error")).Once()

	contact, err := mock.Get(context.Background(), &GetCmd{
		ContactID: "some-contact-id",
	})

	assert.Nil(t, contact)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_Contact_ControllerMock_GetAll(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("GetAll", &GetAllCmd{}).Return(map[string]Contact{
		"some-id": ValidContact,
	}, nil).Once()

	res, err := mock.GetAll(context.Background(), &GetAllCmd{})

	assert.NoError(t, err)
	assert.EqualValues(t, map[string]Contact{
		"some-id": ValidContact,
	}, res)

	mock.AssertExpectations(t)
}

func Test_Contact_ControllerMock_GetAll_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("GetAll", &GetAllCmd{}).Return(nil, fmt.Errorf("some-error")).Once()

	contact, err := mock.GetAll(context.Background(), &GetAllCmd{})

	assert.Nil(t, contact)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_Contact_ControllerMock_Delete(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Delete", &DeleteCmd{ContactID: "some-id"}).Return(nil).Once()

	err := mock.Delete(context.Background(), &DeleteCmd{ContactID: "some-id"})

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

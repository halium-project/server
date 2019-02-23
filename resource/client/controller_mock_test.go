package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Client_ControllerMock_Create(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Create", &CreateCmd{
		ID:            "some-app",
		Name:          "Some App",
		RedirectURIs:  []string{"http://mydomain/oauth/callback"},
		GrantTypes:    []string{"client_credentials", "authorize_code", "implicit", "refresh_token"},
		ResponseTypes: []string{"code", "invalid-response-type"},
		Scopes:        []string{"client", "admin"},
		Public:        true,
	}).Return("", "", fmt.Errorf("some-error")).Once()

	id, secret, err := mock.Create(context.Background(), &CreateCmd{
		ID:            "some-app",
		Name:          "Some App",
		RedirectURIs:  []string{"http://mydomain/oauth/callback"},
		GrantTypes:    []string{"client_credentials", "authorize_code", "implicit", "refresh_token"},
		ResponseTypes: []string{"code", "invalid-response-type"},
		Scopes:        []string{"client", "admin"},
		Public:        true,
	})

	assert.Empty(t, id)
	assert.Empty(t, secret)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_Client_ControllerMock_Get(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		ClientID: "some-client-id",
	}).Return(&ValidClient, nil).Once()

	client, err := mock.Get(context.Background(), &GetCmd{
		ClientID: "some-client-id",
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidClient, client)

	mock.AssertExpectations(t)
}

func Test_Client_ControllerMock_Get_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		ClientID: "some-client-id",
	}).Return(nil, fmt.Errorf("some-error")).Once()

	client, err := mock.Get(context.Background(), &GetCmd{
		ClientID: "some-client-id",
	})

	assert.Nil(t, client)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_Client_ControllerMock_GetAll(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("GetAll", &GetAllCmd{}).Return(map[string]Client{
		"some-id": ValidClient,
	}, nil).Once()

	res, err := mock.GetAll(context.Background(), &GetAllCmd{})

	assert.NoError(t, err)
	assert.EqualValues(t, map[string]Client{
		"some-id": ValidClient,
	}, res)

	mock.AssertExpectations(t)
}

func Test_Client_ControllerMock_GetAll_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("GetAll", &GetAllCmd{}).Return(nil, fmt.Errorf("some-error")).Once()

	client, err := mock.GetAll(context.Background(), &GetAllCmd{})

	assert.Nil(t, client)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_Client_ControllerMock_Delete(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Delete", &DeleteCmd{ClientID: "some-id"}).Return(nil).Once()

	err := mock.Delete(context.Background(), &DeleteCmd{ClientID: "some-id"})

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

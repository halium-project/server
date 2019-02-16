package accesstoken

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/server/util"
	"github.com/stretchr/testify/assert"
)

func Test_AccessToken_ControllerMock_Create(t *testing.T) {
	util.TestIs(t, util.Unit)

	mock := new(ControllerMock)

	mock.On("Create", &CreateCmd{
		ClientID:     ValidAccessToken.ClientID,
		AccessToken:  ValidAccessToken.AccessToken,
		RefreshToken: ValidAccessToken.RefreshToken,
		ExpiresIn:    ValidAccessToken.ExpiresIn,
		Scopes:       ValidAccessToken.Scopes,
	}).Return(nil).Once()

	err := mock.Create(context.Background(), &CreateCmd{
		ClientID:     ValidAccessToken.ClientID,
		AccessToken:  ValidAccessToken.AccessToken,
		RefreshToken: ValidAccessToken.RefreshToken,
		ExpiresIn:    ValidAccessToken.ExpiresIn,
		Scopes:       ValidAccessToken.Scopes,
	})

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

func Test_AccessToken_ControllerMock_Create_with_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	mock := new(ControllerMock)

	mock.On("Create", &CreateCmd{
		ClientID:     ValidAccessToken.ClientID,
		AccessToken:  ValidAccessToken.AccessToken,
		RefreshToken: ValidAccessToken.RefreshToken,
		ExpiresIn:    ValidAccessToken.ExpiresIn,
		Scopes:       ValidAccessToken.Scopes,
	}).Return(fmt.Errorf("some-error")).Once()

	err := mock.Create(context.Background(), &CreateCmd{
		ClientID:     ValidAccessToken.ClientID,
		AccessToken:  ValidAccessToken.AccessToken,
		RefreshToken: ValidAccessToken.RefreshToken,
		ExpiresIn:    ValidAccessToken.ExpiresIn,
		Scopes:       ValidAccessToken.Scopes,
	})

	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_AccessToken_ControllerMock_Get(t *testing.T) {
	util.TestIs(t, util.Unit)

	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		AccessToken: "some-access-token",
	}).Return(&ValidAccessToken, nil).Once()

	res, err := mock.Get(context.Background(), &GetCmd{
		AccessToken: "some-access-token",
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidAccessToken, res)

	mock.AssertExpectations(t)
}

func Test_AccessToken_ControllerMock_Get_with_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		AccessToken: "some-access-token",
	}).Return(nil, fmt.Errorf("some-error")).Once()

	session, err := mock.Get(context.Background(), &GetCmd{
		AccessToken: "some-access-token",
	})

	assert.Nil(t, session)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

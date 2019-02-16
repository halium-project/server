package authorizationcode

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AuthorizationCode_ControllerMock_Create(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Create", &CreateCmd{
		ClientID:            ValidAuthorizationCode.ClientID,
		Code:                "some-authorization-code",
		ExpiresIn:           ValidAuthorizationCode.ExpiresIn,
		Scopes:              ValidAuthorizationCode.Scopes,
		RedirectURI:         ValidAuthorizationCode.RedirectURI,
		State:               ValidAuthorizationCode.State,
		CodeChallenge:       ValidAuthorizationCode.CodeChallenge,
		CodeChallengeMethod: ValidAuthorizationCode.CodeChallengeMethod,
	}).Return(nil).Once()

	err := mock.Create(context.Background(), &CreateCmd{
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

	mock.AssertExpectations(t)
}

func Test_AuthorizationCode_ControllerMock_Create_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Create", &CreateCmd{
		ClientID:            ValidAuthorizationCode.ClientID,
		Code:                "some-authorization-code",
		ExpiresIn:           ValidAuthorizationCode.ExpiresIn,
		Scopes:              ValidAuthorizationCode.Scopes,
		RedirectURI:         ValidAuthorizationCode.RedirectURI,
		State:               ValidAuthorizationCode.State,
		CodeChallenge:       ValidAuthorizationCode.CodeChallenge,
		CodeChallengeMethod: ValidAuthorizationCode.CodeChallengeMethod,
	}).Return(fmt.Errorf("some-error")).Once()

	err := mock.Create(context.Background(), &CreateCmd{
		ClientID:            ValidAuthorizationCode.ClientID,
		Code:                "some-authorization-code",
		ExpiresIn:           ValidAuthorizationCode.ExpiresIn,
		Scopes:              ValidAuthorizationCode.Scopes,
		RedirectURI:         ValidAuthorizationCode.RedirectURI,
		State:               ValidAuthorizationCode.State,
		CodeChallenge:       ValidAuthorizationCode.CodeChallenge,
		CodeChallengeMethod: ValidAuthorizationCode.CodeChallengeMethod,
	})

	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_AuthorizationCode_ControllerMock_Get(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		Code: "some-authorization-code",
	}).Return(&ValidAuthorizationCode, nil).Once()

	res, err := mock.Get(context.Background(), &GetCmd{
		Code: "some-authorization-code",
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidAuthorizationCode, res)

	mock.AssertExpectations(t)
}

func Test_AuthorizationCode_ControllerMock_Get_with_error(t *testing.T) {
	mock := new(ControllerMock)

	mock.On("Get", &GetCmd{
		Code: "some-authorization-code",
	}).Return(nil, fmt.Errorf("some-error")).Once()

	session, err := mock.Get(context.Background(), &GetCmd{
		Code: "some-authorization-code",
	})

	assert.Nil(t, session)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

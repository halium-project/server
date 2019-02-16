package accesstoken

import (
	"context"
	"testing"

	"github.com/halium-project/server/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_AccessToken_StorageMock_Set(t *testing.T) {
	util.TestIs(t, util.Unit)

	mock := new(StorageMock)

	mock.On("Set", "some-id", "some-rev", &ValidAccessToken).Return("some-rev-2", nil)

	rev, err := mock.Set(context.Background(), "some-id", "some-rev", &ValidAccessToken)

	assert.NoError(t, err)
	assert.Equal(t, "some-rev-2", rev)

	mock.AssertExpectations(t)
}

func Test_AccessToken_StorageMock_Get(t *testing.T) {
	util.TestIs(t, util.Unit)

	mock := new(StorageMock)

	mock.On("Get", "some-id").Return("some-rev", &ValidAccessToken, nil)

	rev, res, err := mock.Get(context.Background(), "some-id")
	assert.NoError(t, err)
	assert.Equal(t, "some-rev", rev)
	assert.EqualValues(t, &ValidAccessToken, res)

	mock.AssertExpectations(t)
}

func Test_AccessToken_StorageMock_Get_with_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	mock := new(StorageMock)

	mock.On("Get", "some-id").Return("", nil, errors.New("some-error"))

	rev, res, err := mock.Get(context.Background(), "some-id")
	assert.EqualError(t, err, "some-error")
	assert.Empty(t, rev)
	assert.Nil(t, res)

	mock.AssertExpectations(t)
}

func Test_AccessToken_StorageMock_Delete(t *testing.T) {
	util.TestIs(t, util.Unit)

	mock := new(StorageMock)

	mock.On("Delete", "some-id", "some-rev").Return(errors.New("some-error"))

	err := mock.Delete(context.Background(), "some-id", "some-rev")
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

func Test_AccessToken_StorageMock_FindOneByRefreshToken(t *testing.T) {
	util.TestIs(t, util.Unit)

	mock := new(StorageMock)

	mock.On("FindOneByRefreshToken", "some-name").Return("some-id", "some-rev", &ValidAccessToken, nil)

	id, rev, res, err := mock.FindOneByRefreshToken(context.Background(), "some-name")

	assert.NoError(t, err)
	assert.Equal(t, "some-id", id)
	assert.Equal(t, "some-rev", rev)
	assert.EqualValues(t, &ValidAccessToken, res)

	mock.AssertExpectations(t)
}

func Test_AccessToken_StorageMock_FindOneByRefreshToken_with_error(t *testing.T) {
	util.TestIs(t, util.Unit)

	mock := new(StorageMock)

	mock.On("FindOneByRefreshToken", "some-name").Return("", "", nil, errors.New("some-error"))

	id, rev, res, err := mock.FindOneByRefreshToken(context.Background(), "some-name")

	assert.Empty(t, rev)
	assert.Empty(t, id)
	assert.Nil(t, res)
	assert.EqualError(t, err, "some-error")

	mock.AssertExpectations(t)
}

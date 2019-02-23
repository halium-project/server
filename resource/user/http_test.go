package user

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/server/resource/accesstoken"
	"github.com/halium-project/server/utils/permission"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_User_HTTPHandler_Create_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Create", &CreateCmd{
		Username: "John",
		Password: "password123",
		Role:     "admin",
	}).Return("some-user-id", nil).Once()

	r := httptest.NewRequest("POST", "http://example.com/users", strings.NewReader(`{
		"username": "John",
		"password": "password123",
		"role": "admin"
	}`))
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.JSONEq(t, `{ "id": "some-user-id" }`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_User_HTTPHandler_Create_with_an_invalid_json(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	r := httptest.NewRequest("POST", "http://example.com/users", strings.NewReader("not a json"))
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "invalidJSON",
		"message": "invalid character 'o' in literal null (expecting 'u')"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_User_HTTPHandler_Create_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Create", &CreateCmd{
		Username: "John",
		Password: "password123",
		Role:     "admin",
	}).Return("", errors.New(errors.NotFound, "user not found")).Once()

	r := httptest.NewRequest("POST", "http://example.com/users", strings.NewReader(`{
		"username": "John",
		"password": "password123",
		"role": "admin"
	}`))
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "notFound",
		"message": "user not found"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_User_HTTPHandler_Get_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Get", &GetCmd{UserID: "some-user-id"}).Return(&ValidUser, nil).Once()

	r := httptest.NewRequest("GET", "http://example.com/users/some-user-id", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, `{
		"username": "some username",
		"role": "admin"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_User_HTTPHandler_Get_with_the_user_not_found(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Get", &GetCmd{UserID: "some-user-id"}).Return(nil, nil).Once()

	r := httptest.NewRequest("GET", "http://example.com/users/some-user-id", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "notFound",
		"message": "user \"some-user-id\" not found"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_User_HTTPHandler_Get_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Get", &GetCmd{UserID: "some-user-id"}).Return(nil, errors.New(errors.Internal, "some error")).Once()

	r := httptest.NewRequest("GET", "http://example.com/users/some-user-id", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "internalError",
		"message": "some error"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_User_HTTPHandler_Update_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Update", &UpdateCmd{
		UserID:   "some-user-id",
		Username: "John",
		Role:     "admin",
	}).Return(nil).Once()

	r := httptest.NewRequest("PUT", "http://example.com/users/some-user-id", strings.NewReader(`{
		"username": "John",
		"role": "admin"
	}`))
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, `{}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_User_HTTPHandler_Update_with_an_invalid_json_request(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	r := httptest.NewRequest("PUT", "http://example.com/users/some-user-id", strings.NewReader("not a valid json"))
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "invalidJSON",
		"message": "invalid character 'o' in literal null (expecting 'u')"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_User_HTTPHandler_Update_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Update", &UpdateCmd{
		UserID:   "some-user-id",
		Username: "John",
		Role:     "admin",
	}).Return(errors.New(errors.Internal, "some error")).Once()

	r := httptest.NewRequest("PUT", "http://example.com/users/some-user-id", strings.NewReader(`{
		"username": "John",
		"role": "admin"
	}`))
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "internalError",
		"message": "some error"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_User_HTTPHandler_GetAll_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("GetAll", &GetAllCmd{}).Return(map[string]User{"some-user-id": ValidUser}, nil).Once()

	r := httptest.NewRequest("GET", "http://example.com/users", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, `{
		"some-user-id": {
			"username": "some username",
			"role": "admin"
		}
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_User_HTTPHandler_GetAll_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("GetAll", &GetAllCmd{}).Return(nil, errors.New(errors.Internal, "some error")).Once()

	r := httptest.NewRequest("GET", "http://example.com/users", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "internalError",
		"message": "some error"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

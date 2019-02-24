package client

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

func Test_Client_HTTPHandler_Create_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Create", &CreateCmd{
		ID:            "my-app",
		Name:          "My App",
		RedirectURIs:  []string{"url-1", "url-2"},
		GrantTypes:    []string{"password"},
		ResponseTypes: []string{"token"},
		Scopes:        []string{"scope-1", "scope-2"},
	}).Return("some-client-id", "some-secret", nil).Once()

	r := httptest.NewRequest("POST", "http://example.com/clients", strings.NewReader(`{
		"id": "my-app",
		"name": "My App",
		"redirectURIs": ["url-1", "url-2"],
		"grantTypes": ["password"],
		"responseTypes": ["token"],
		"scopes": ["scope-1", "scope-2"]
	}`))
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.JSONEq(t, `{
		"clientID": "some-client-id",
		"clientSecret": "some-secret"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_Client_HTTPHandler_Create_with_an_invalid_json(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	r := httptest.NewRequest("POST", "http://example.com/clients", strings.NewReader("not a json"))
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

func Test_Client_HTTPHandler_Create_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Create", &CreateCmd{
		ID:            "my-app",
		Name:          "My App",
		RedirectURIs:  []string{"url-1", "url-2"},
		GrantTypes:    []string{"password"},
		ResponseTypes: []string{"token"},
		Scopes:        []string{"scope-1", "scope-2"},
	}).Return("", "", errors.New(errors.NotFound, "client not found")).Once()

	r := httptest.NewRequest("POST", "http://example.com/clients", strings.NewReader(`{
		"id": "my-app",
		"name": "My App",
		"redirectURIs": ["url-1", "url-2"],
		"grantTypes": ["password"],
		"responseTypes": ["token"],
		"scopes": ["scope-1", "scope-2"]
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
		"message": "client not found"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_Client_HTTPHandler_Get_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Get", &GetCmd{ClientID: "some-client-id"}).Return(&ValidClient, nil).Once()

	r := httptest.NewRequest("GET", "http://example.com/clients/some-client-id", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, `{
		"id": "my-web-application",
		"name": "My Web Application",
		"redirectURIs":  ["http://mydomain/oauth/callback"],
		"grantTypes": ["client_credentials", "authorize_code", "implicit", "refresh_token", "password"],
		"responseTypes": ["code", "token"],
		"scopes": ["user", "admin"],
		"public": false
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_Client_HTTPHandler_Get_with_the_client_not_found(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Get", &GetCmd{ClientID: "some-client-id"}).Return(nil, nil).Once()

	r := httptest.NewRequest("GET", "http://example.com/clients/some-client-id", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "notFound",
		"message": "client \"some-client-id\" not found"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_Client_HTTPHandler_Get_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Get", &GetCmd{ClientID: "some-client-id"}).Return(nil, errors.New(errors.Internal, "some error")).Once()

	r := httptest.NewRequest("GET", "http://example.com/clients/some-client-id", nil)
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

func Test_Client_HTTPHandler_GetAll_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("GetAll", &GetAllCmd{}).Return(map[string]Client{"some-client-id": ValidClient}, nil).Once()

	r := httptest.NewRequest("GET", "http://example.com/clients", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, `{
		"some-client-id": {
			"id": "my-web-application",
			"name": "My Web Application",
			"redirectURIs":  ["http://mydomain/oauth/callback"],
			"grantTypes": ["client_credentials", "authorize_code", "implicit", "refresh_token", "password"],
			"responseTypes": ["code", "token"],
			"scopes": ["user", "admin"],
			"public": false
		}
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_Client_HTTPHandler_GetAll_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("GetAll", &GetAllCmd{}).Return(nil, errors.New(errors.Internal, "some error")).Once()

	r := httptest.NewRequest("GET", "http://example.com/clients", nil)
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

func Test_Client_HTTPHandler_Delete_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Delete", &DeleteCmd{ClientID: "some-client-id"}).Return(nil).Once()

	r := httptest.NewRequest("DELETE", "http://example.com/clients/some-client-id", nil)
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

func Test_Client_HTTPHandler_Delete_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Delete", &DeleteCmd{ClientID: "some-client-id"}).Return(errors.New(errors.BadRequest, "some-error")).Once()

	r := httptest.NewRequest("DELETE", "http://example.com/clients/some-client-id", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "badRequest",
		"message": "some-error"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

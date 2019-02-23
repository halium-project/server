package contact

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

func Test_Contact_HTTPHandler_Create_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Create", &CreateCmd{
		Name: "Jane Doe",
	}).Return("some-contact-id", nil).Once()

	r := httptest.NewRequest("POST", "http://example.com/contacts", strings.NewReader(`{
		"name": "Jane Doe"
	}`))
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.JSONEq(t, `{ "id": "some-contact-id" }`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_Contact_HTTPHandler_Create_with_an_invalid_json(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	r := httptest.NewRequest("POST", "http://example.com/contacts", strings.NewReader("not a json"))
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

func Test_Contact_HTTPHandler_Create_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Create", &CreateCmd{
		Name: "Jane Doe",
	}).Return("", errors.New(errors.NotFound, "contact not found")).Once()

	r := httptest.NewRequest("POST", "http://example.com/contacts", strings.NewReader(`{
		"name": "Jane Doe"
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
		"message": "contact not found"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_Contact_HTTPHandler_Get_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Get", &GetCmd{ContactID: "some-contact-id"}).Return(&ValidContact, nil).Once()

	r := httptest.NewRequest("GET", "http://example.com/contacts/some-contact-id", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, `{
		"name": "Jane Doe"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_Contact_HTTPHandler_Get_with_the_contact_not_found(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Get", &GetCmd{ContactID: "some-contact-id"}).Return(nil, nil).Once()

	r := httptest.NewRequest("GET", "http://example.com/contacts/some-contact-id", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.JSONEq(t, `{
		"kind": "notFound",
		"message": "contact \"some-contact-id\" not found"
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_Contact_HTTPHandler_Get_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Get", &GetCmd{ContactID: "some-contact-id"}).Return(nil, errors.New(errors.Internal, "some error")).Once()

	r := httptest.NewRequest("GET", "http://example.com/contacts/some-contact-id", nil)
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

func Test_Contact_HTTPHandler_GetAll_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("GetAll", &GetAllCmd{}).Return(map[string]Contact{"some-contact-id": ValidContact}, nil).Once()

	r := httptest.NewRequest("GET", "http://example.com/contacts", nil)
	r.Header.Add("Authorization", "Bearer foobar")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, `{
		"some-contact-id": {
			"name": "Jane Doe"
		}
	}`, string(body))

	controllerMock.AssertExpectations(t)
	accessTokenControllerMock.AssertExpectations(t)
}

func Test_Contact_HTTPHandler_GetAll_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("GetAll", &GetAllCmd{}).Return(nil, errors.New(errors.Internal, "some error")).Once()

	r := httptest.NewRequest("GET", "http://example.com/contacts", nil)
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

func Test_Contact_HTTPHandler_Delete_success(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Delete", &DeleteCmd{ContactID: "some-contact-id"}).Return(nil).Once()

	r := httptest.NewRequest("DELETE", "http://example.com/contacts/some-contact-id", nil)
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

func Test_Contact_HTTPHandler_Delete_with_an_error_from_the_usecase(t *testing.T) {
	accessTokenControllerMock := new(accesstoken.ControllerMock)
	perm := permission.NewController(context.Background(), accessTokenControllerMock)
	router := mux.NewRouter()
	controllerMock := new(ControllerMock)
	handler := NewHTTPHandler(controllerMock)
	handler.RegisterRoutes(router, perm)

	// Token introspection
	accessTokenControllerMock.On("Get", &accesstoken.GetCmd{AccessToken: "foobar"}).Return(&accesstoken.ValidAccessToken, nil).Once()

	controllerMock.On("Delete", &DeleteCmd{ContactID: "some-contact-id"}).Return(errors.New(errors.BadRequest, "some-error")).Once()

	r := httptest.NewRequest("DELETE", "http://example.com/contacts/some-contact-id", nil)
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

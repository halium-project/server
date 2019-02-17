package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/go-server-utils/response"
)

type HTTPHandler struct {
	client ControllerInterface
}

type ControllerInterface interface {
	Create(ctx context.Context, cmd *CreateCmd) (string, string, error)
	Get(ctx context.Context, cmd *GetCmd) (*Client, error)
	GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]Client, error)
}

func NewHTTPHandler(client ControllerInterface) *HTTPHandler {
	return &HTTPHandler{
		client: client,
	}
}

func (t *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name          string   `json:"name"`
		RedirectURIs  []string `json:"redirectURI"`
		GrantTypes    []string `json:"grantTypes"`
		ResponseTypes []string `json:"responseTypes"`
		Scopes        []string `json:"scopes"`
		Public        bool     `json:"public"`
	}

	type responseBody struct {
		ClientID     string `json:"clientID"`
		ClientSecret string `json:"clientSecret"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.IntoResponse(w, errors.New(errors.InvalidJSON, err.Error()))
		return
	}

	id, secret, err := t.client.Create(r.Context(), &CreateCmd{
		Name:          req.Name,
		RedirectURIs:  req.RedirectURIs,
		ResponseTypes: req.ResponseTypes,
		GrantTypes:    req.GrantTypes,
		Scopes:        req.Scopes,
		Public:        req.Public,
	})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	response.Write(w, http.StatusCreated, &responseBody{
		ClientID:     id,
		ClientSecret: secret,
	})
}

func (t *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["clientID"]
	client, err := t.client.Get(r.Context(), &GetCmd{
		ClientID: clientID,
	})

	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	if client == nil {
		errors.IntoResponse(w, errors.Errorf(errors.NotFound, "client %q not found", clientID))
		return
	}

	response.Write(w, http.StatusOK, &client)
}

func (t *HTTPHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	clients, err := t.client.GetAll(r.Context(), &GetAllCmd{})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	response.Write(w, http.StatusOK, clients)
}

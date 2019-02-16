package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/halium-project/server/util/errors"
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

	type response struct {
		ClientID     string `json:"clientID"`
		ClientSecret string `json:"clientSecret"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req request
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
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

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&response{
		ClientID:     id,
		ClientSecret: secret,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	type response Client

	client, err := t.client.Get(r.Context(), &GetCmd{
		ClientID: mux.Vars(r)["clientID"],
	})

	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	if client == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res := response(*client)

	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *HTTPHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	clients, err := t.client.GetAll(r.Context(), &GetAllCmd{})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(clients)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/go-server-utils/response"
	"github.com/halium-project/server/utils/permission"
)

type HTTPHandler struct {
	client ControllerInterface
}

type ControllerInterface interface {
	Create(ctx context.Context, cmd *CreateCmd) (string, string, error)
	Get(ctx context.Context, cmd *GetCmd) (*Client, error)
	GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]Client, error)
	Delete(ctx context.Context, cmd *DeleteCmd) error
}

func NewHTTPHandler(client ControllerInterface) *HTTPHandler {
	return &HTTPHandler{
		client: client,
	}
}

func (t *HTTPHandler) RegisterRoutes(router *mux.Router, perm *permission.Controller) {
	router.HandleFunc("/clients", perm.Check("clients.write", t.Create)).Methods("POST")
	router.HandleFunc("/clients", perm.Check("clients.read", t.GetAll)).Methods("GET")
	router.HandleFunc("/clients/{clientID}", perm.Check("clients.read", t.Get)).Methods("GET")
	router.HandleFunc("/clients/{clientID}", perm.Check("clients.write", t.Delete)).Methods("Delete")
}

func (t *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ID            string   `json:"id"`
		Name          string   `json:"name"`
		RedirectURIs  []string `json:"redirectURIs"`
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
		ID:            req.ID,
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

	// Do not return the secret
	client.Secret = ""

	response.Write(w, http.StatusOK, &client)
}

func (t *HTTPHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	clients, err := t.client.GetAll(r.Context(), &GetAllCmd{})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	for _, client := range clients {
		// Do not return the secret
		client.Secret = ""
	}

	response.Write(w, http.StatusOK, clients)
}

func (t *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	clientID := mux.Vars(r)["clientID"]

	err := t.client.Delete(r.Context(), &DeleteCmd{
		ClientID: clientID,
	})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	// Do not return the password and the salt.
	response.Write(w, http.StatusOK, struct{}{})
}

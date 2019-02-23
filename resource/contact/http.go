package contact

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
	contact ControllerInterface
}

type ControllerInterface interface {
	Create(ctx context.Context, cmd *CreateCmd) (string, error)
	Get(ctx context.Context, cmd *GetCmd) (*Contact, error)
	GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]Contact, error)
	Delete(ctx context.Context, cmd *DeleteCmd) error
}

func NewHTTPHandler(contact ControllerInterface) *HTTPHandler {
	return &HTTPHandler{
		contact: contact,
	}
}

func (t *HTTPHandler) RegisterRoutes(router *mux.Router, perm *permission.Controller) {
	router.HandleFunc("/contacts", perm.Check("contacts.write", t.Create)).Methods("POST")
	router.HandleFunc("/contacts", perm.Check("contacts.read", t.GetAll)).Methods("GET")
	router.HandleFunc("/contacts/{contactID}", perm.Check("contacts.read", t.Get)).Methods("GET")
	router.HandleFunc("/contacts/{contactID}", perm.Check("contacts.write", t.Delete)).Methods("Delete")
}

func (t *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name string `json:"name"`
	}

	type responseBody struct {
		ContactID string `json:"id"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.IntoResponse(w, errors.New(errors.InvalidJSON, err.Error()))
		return
	}

	id, err := t.contact.Create(r.Context(), &CreateCmd{
		Name: req.Name,
	})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	response.Write(w, http.StatusCreated, &responseBody{
		ContactID: id,
	})
}

func (t *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	contactID := mux.Vars(r)["contactID"]
	contact, err := t.contact.Get(r.Context(), &GetCmd{
		ContactID: contactID,
	})

	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	if contact == nil {
		errors.IntoResponse(w, errors.Errorf(errors.NotFound, "contact %q not found", contactID))
		return
	}

	response.Write(w, http.StatusOK, &contact)
}

func (t *HTTPHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	contacts, err := t.contact.GetAll(r.Context(), &GetAllCmd{})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	response.Write(w, http.StatusOK, contacts)
}

func (t *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	contactID := mux.Vars(r)["contactID"]

	err := t.contact.Delete(r.Context(), &DeleteCmd{
		ContactID: contactID,
	})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	// Do not return the password and the salt.
	response.Write(w, http.StatusOK, struct{}{})
}

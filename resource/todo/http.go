package todo

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
	todo ControllerInterface
}

type ControllerInterface interface {
	Create(ctx context.Context, cmd *CreateCmd) (string, error)
	Get(ctx context.Context, cmd *GetCmd) (*Todo, error)
	GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]Todo, error)
	Delete(ctx context.Context, cmd *DeleteCmd) error
}

func NewHTTPHandler(todo ControllerInterface) *HTTPHandler {
	return &HTTPHandler{
		todo: todo,
	}
}

func (t *HTTPHandler) RegisterRoutes(router *mux.Router, perm *permission.Controller) {
	router.HandleFunc("/todos", perm.Check("todos.write", t.Create)).Methods("POST")
	router.HandleFunc("/todos", perm.Check("todos.read", t.GetAll)).Methods("GET")
	router.HandleFunc("/todos/{todoID}", perm.Check("todos.read", t.Get)).Methods("GET")
	router.HandleFunc("/todos/{todoID}", perm.Check("todos.write", t.Delete)).Methods("Delete")
}

func (t *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title string `json:"name"`
	}

	type responseBody struct {
		TodoID string `json:"id"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.IntoResponse(w, errors.New(errors.InvalidJSON, err.Error()))
		return
	}

	id, err := t.todo.Create(r.Context(), &CreateCmd{
		Title: req.Title,
	})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	response.Write(w, http.StatusCreated, &responseBody{
		TodoID: id,
	})
}

func (t *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	todoID := mux.Vars(r)["todoID"]
	todo, err := t.todo.Get(r.Context(), &GetCmd{
		TodoID: todoID,
	})

	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	if todo == nil {
		errors.IntoResponse(w, errors.Errorf(errors.NotFound, "todo %q not found", todoID))
		return
	}

	response.Write(w, http.StatusOK, &todo)
}

func (t *HTTPHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	todos, err := t.todo.GetAll(r.Context(), &GetAllCmd{})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	response.Write(w, http.StatusOK, todos)
}

func (t *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	todoID := mux.Vars(r)["todoID"]

	err := t.todo.Delete(r.Context(), &DeleteCmd{
		TodoID: todoID,
	})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	// Do not return the password and the salt.
	response.Write(w, http.StatusOK, struct{}{})
}

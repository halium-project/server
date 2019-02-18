package user

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
	user ControllerInterface
}

type ControllerInterface interface {
	Get(ctx context.Context, cmd *GetCmd) (*User, error)
	Create(ctx context.Context, cmd *CreateCmd) (string, error)
	Update(ctx context.Context, cmd *UpdateCmd) error
	GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]User, error)
}

func NewHTTPHandler(user ControllerInterface) *HTTPHandler {
	return &HTTPHandler{
		user: user,
	}
}

func (t *HTTPHandler) Register(router *mux.Router, perm *permission.Controller) {
	router.HandleFunc("/users", perm.Check("users.write", t.Create)).Methods("POST")
	router.HandleFunc("/users", perm.Check("users.read", t.GetAll)).Methods("GET")
	router.HandleFunc("/users/{userID}", perm.Check("users.write", t.Update)).Methods("PUT")
	router.HandleFunc("/users/{userID}", perm.Check("users.read", t.Get)).Methods("GET")
}

func (t *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Role     string `json:"role"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type responseBody struct {
		UserID string `json:"id"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.IntoResponse(w, errors.New(errors.InvalidJSON, err.Error()))
		return
	}

	userID, err := t.user.Create(r.Context(), &CreateCmd{
		Role:     req.Role,
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	response.Write(w, http.StatusCreated, &responseBody{
		UserID: userID,
	})
}

func (t *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	type responseBody struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	userID := mux.Vars(r)["userID"]
	user, err := t.user.Get(r.Context(), &GetCmd{
		UserID: userID,
	})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	if user == nil {
		errors.IntoResponse(w, errors.Errorf(errors.NotFound, "user %q not found", userID))
		return
	}

	// Do not return the password and the salt.
	response.Write(w, http.StatusOK, &responseBody{
		Username: user.Username,
		Role:     user.Role,
	})
}

func (t *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Role     string `json:"role"`
		Username string `json:"username"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.IntoResponse(w, errors.New(errors.InvalidJSON, err.Error()))
		return
	}

	err = t.user.Update(r.Context(), &UpdateCmd{
		UserID:   mux.Vars(r)["userID"],
		Role:     req.Role,
		Username: req.Username,
	})

	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	response.Write(w, http.StatusOK, &struct{}{})
}

func (t *HTTPHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	type userRes struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	res, err := t.user.GetAll(r.Context(), &GetAllCmd{})
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	users := make(map[string]userRes, len(res))

	for id, user := range res {
		users[id] = userRes{
			Username: user.Username,
			Role:     user.Role,
		}
	}

	response.Write(w, http.StatusOK, users)
}

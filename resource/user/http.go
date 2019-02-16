package user

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/halium-project/server/util/errors"
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

func NewHTTPHandler(
	user ControllerInterface) *HTTPHandler {
	return &HTTPHandler{
		user: user,
	}
}

func (t *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Role     string `json:"role"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		UserID string `json:"id"`
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

	userID, err := t.user.Create(r.Context(), &CreateCmd{
		Role:     req.Role,
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&response{
		UserID: userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	type getUserResponse User

	user, err := t.user.Get(r.Context(), &GetCmd{
		UserID: mux.Vars(r)["userID"],
	})

	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := getUserResponse(*user)

	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Role     string `json:"role"`
		Username string `json:"username"`
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

	err = t.user.Update(r.Context(), &UpdateCmd{
		UserID:   mux.Vars(r)["userID"],
		Role:     req.Role,
		Username: req.Username,
	})

	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	_, err = w.Write([]byte("{}"))
	if err != nil {
		errors.IntoResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

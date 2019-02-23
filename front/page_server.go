package front

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/server/resource/user"
)

type HTMLRenderer interface {
	Render(w io.Writer, templateName string, params interface{}) error
}

type UserCreator interface {
	Create(ctx context.Context, cmd *user.CreateCmd) (string, error)
}

type PageServer struct {
	renderer HTMLRenderer
	user     UserCreator
}

func NewPageServer(renderer HTMLRenderer, user UserCreator) *PageServer {
	return &PageServer{
		renderer: renderer,
		user:     user,
	}
}

func (t *PageServer) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", t.Register)
}

func (t *PageServer) Register(w http.ResponseWriter, r *http.Request) {
	type registerTemplateParam struct {
		Errors map[string]string
	}

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		t.renderer.Render(w, "register.html", registerTemplateParam{})
		return
	}

	if r.Method != "POST" {
		errors.IntoResponse(w, errors.New(errors.InvalidMethod, "only GET and POST methods are accepted"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		errors.IntoResponse(w, errors.Errorf(errors.BadRequest, err.Error()))
		return
	}

	_, err = t.user.Create(r.Context(), &user.CreateCmd{
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
		Role:     user.Dev,
	})
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if errors.IsUnexpected(err) {
			w.WriteHeader(http.StatusInternalServerError)
			t.renderer.Render(w, "internal_error.html", err.Error())
			return
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		t.renderer.Render(w, "register.html", registerTemplateParam{
			Errors: err.(*errors.Error).Errors,
		})
		return
	}

	http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
}

package oauth2

import (
	"context"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/halium-project/server/resource/user"
	"github.com/openshift/osin"
	"gitlab.com/Peltoche/yaccc"
)

type Controller struct {
	inner *osin.Server
	html  TemplateRenderer
	user  UserValidater
}

type TemplateRenderer interface {
	Render(w io.Writer, templateName string, params interface{}) error
}

type UserValidater interface {
	Validate(ctx context.Context, cmd *user.ValidateCmd) (string, *user.User, error)
}

func InitController(
	ctx context.Context,
	server *yaccc.Server,
	html TemplateRenderer,
	user UserValidater,
	storage osin.Storage,
) *Controller {
	osinConfig := osin.NewServerConfig()
	osinConfig.AllowedAuthorizeTypes = []osin.AuthorizeRequestType{osin.CODE, osin.TOKEN}
	osinConfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE, osin.REFRESH_TOKEN, osin.CLIENT_CREDENTIALS, osin.IMPLICIT}
	osinConfig.AccessExpiration = 3
	osinConfig.ErrorStatusCode = http.StatusBadRequest

	osinServer := osin.NewServer(osinConfig, storage)
	osinServer.Logger = log.New(os.Stdout, "", log.LstdFlags)

	return NewController(osinServer, html, user, storage)

}

func NewController(
	osinServer *osin.Server,
	html TemplateRenderer,
	user UserValidater,
	storage osin.Storage,
) *Controller {
	return &Controller{
		inner: osinServer,
		html:  html,
		user:  user,
	}
}

func (t *Controller) Authorize(w http.ResponseWriter, r *http.Request) {
	resp := t.inner.NewResponse()
	defer resp.Close()

	ar := t.inner.HandleAuthorizeRequest(resp, r)
	if ar != nil {

		if r.Method == "GET" {
			t.renderAuthenticationPage(w, http.StatusOK, nil)
			return
		}

		err := r.ParseForm()
		if err != nil {
			resp.SetErrorState(osin.E_SERVER_ERROR, "", ar.State)
			resp.InternalError = err
			t.inner.FinishAuthorizeRequest(resp, r, ar)
			return
		}

		userID, _, err := t.user.Validate(r.Context(), &user.ValidateCmd{
			Username: r.PostForm.Get("username"),
			Password: r.PostForm.Get("password"),
		})
		if err != nil {
			resp.SetErrorState(osin.E_SERVER_ERROR, "", ar.State)
			resp.InternalError = err
			t.inner.FinishAuthorizeRequest(resp, r, ar)
			return
		}

		if userID == "" {
			t.renderAuthenticationPage(w, http.StatusBadRequest, nil)
			return
		}

		if ar.Type == "token" {
			// This is an implicit grant. It is used by insecure clients
			//like web apps and so no refresh token are required.
			ar.Expiration = int32(math.MaxInt32)
		}

		ar.Authorized = true
		t.inner.FinishAuthorizeRequest(resp, r, ar)
	}

	err := osin.OutputJSON(resp, w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (t *Controller) Token(w http.ResponseWriter, r *http.Request) {
	resp := t.inner.NewResponse()
	defer resp.Close()

	ar := t.inner.HandleAccessRequest(resp, r)
	if ar != nil {
		ar.Authorized = true
		t.inner.FinishAccessRequest(resp, r, ar)
	}

	err := osin.OutputJSON(resp, w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (t *Controller) Info(w http.ResponseWriter, r *http.Request) {
	resp := t.inner.NewResponse()
	defer resp.Close()

	if ir := t.inner.HandleInfoRequest(resp, r); ir != nil {
		t.inner.FinishInfoRequest(resp, r, ir)
	}
	err := osin.OutputJSON(resp, w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (t *Controller) renderAuthenticationPage(w http.ResponseWriter, HTTPStatus int, param interface{}) {
	w.WriteHeader(HTTPStatus)
	err := t.html.Render(w, "auth.html", param)
	if err != nil {
		log.Println(err)
	}
}

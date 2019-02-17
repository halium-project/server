package permission

import (
	"context"
	"net/http"
	"strings"

	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/server/resource/accesstoken"
)

type AccessTokenGetter interface {
	Get(ctx context.Context, cmd *accesstoken.GetCmd) (*accesstoken.AccessToken, error)
}

type Controller struct {
	accessToken AccessTokenGetter
}

func NewController(ctx context.Context, accessToken AccessTokenGetter) *Controller {
	return &Controller{
		accessToken: accessToken,
	}
}

func (t *Controller) Check(permission string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, parseErr := RetrieveTokenFromRequest(r)
		if parseErr != nil {
			errors.WriteError(w, parseErr)
			return
		}

		session, err := t.accessToken.Get(r.Context(), &accesstoken.GetCmd{
			AccessToken: token,
		})

		if err != nil {
			errors.WriteError(w, errors.Wrapf(err, "failed to retrieve session %q", token))
			return
		}

		if session == nil {
			errors.WriteError(w, errors.New(errors.NotAuthorized, "invalid session"))
			return
		}

		var isAuthorized bool
		for _, scope := range session.Scopes {
			if strings.HasPrefix(permission, scope) {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			errors.WriteError(w, errors.New(errors.NotAuthorized, "doesn't have required permission"))
			return
		}

		handler(w, r)

	}
}

func RetrieveTokenFromRequest(r *http.Request) (string, *errors.Error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New(errors.BadRequest, `missing "Authorization" header`)
	}

	elements := strings.Split(header, " ")

	if len(elements) != 2 {
		return "", errors.New(errors.BadRequest, `invalid "Authorization" header format`)
	}

	if elements[0] != "Bearer" {
		return "", errors.New(errors.BadRequest, `"Authorization" header must be of kind "Bearer"`)
	}

	return elements[1], nil
}

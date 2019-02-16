package authorizationcode

import (
	"context"
	"log"
	"time"

	"github.com/halium-project/server/db"
	"github.com/halium-project/server/util/errors"
	"github.com/halium-project/server/util/password"
	"github.com/halium-project/server/util/uuid"
	"github.com/halium-project/server/util/validator"
	"github.com/halium-project/server/util/validator/is"
	"gitlab.com/Peltoche/yaccc"
)

type Controller struct {
	uuid     uuid.Producer
	storage  StorageInterface
	password password.HashManager
}

type StorageInterface interface {
	Set(ctx context.Context, id string, rev string, value *AuthorizationCode) (string, error)
	Get(ctx context.Context, code string) (string, *AuthorizationCode, error)
	Delete(ctx context.Context, code string, rev string) error
}

func InitController(ctx context.Context, server *yaccc.Server) *Controller {
	database, err := server.ConnectDatabase(ctx, BucketName)
	if err != nil {
		database, err = SetupStorage(ctx, server)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "failed to setup %q storage", BucketName))
		}
	}

	storage := NewStorage(db.NewCouchdbDriver(database))
	uuidProducer := uuid.NewGoUUID()
	passwordProducer := password.NewPasswordHasher()

	controller := NewController(uuidProducer, passwordProducer, storage)

	return controller
}

func NewController(
	uuid uuid.Producer,
	password password.HashManager,
	storage StorageInterface,
) *Controller {
	return &Controller{
		uuid:     uuid,
		password: password,
		storage:  storage,
	}
}

func (t *Controller) Create(ctx context.Context, cmd *CreateCmd) error {
	err := validator.New().
		CheckString("clientId", cmd.ClientID, is.Required, is.StringInRange(3, 100)).
		CheckString("code", cmd.Code, is.Required, is.StringInRange(8, 256)).
		CheckNumber("expiresIn", cmd.ExpiresIn, is.Required, is.NumberPositif).
		CheckString("redirectURI", cmd.RedirectURI, is.Required, is.StringInRange(3, 512)).
		CheckString("state", cmd.RedirectURI, is.Optional).
		CheckString("codeChallenge", cmd.CodeChallenge, is.Optional).
		CheckString("codeChallengeMethod", cmd.CodeChallengeMethod, is.Optional).
		Run()
	if err != nil {
		return err
	}

	// Save the document
	_, err = t.storage.Set(ctx, cmd.Code, "", &AuthorizationCode{
		ClientID:            cmd.ClientID,
		ExpiresIn:           cmd.ExpiresIn,
		Scopes:              cmd.Scopes,
		RedirectURI:         cmd.RedirectURI,
		State:               cmd.State,
		CreatedAt:           time.Now(),
		CodeChallenge:       cmd.CodeChallenge,
		CodeChallengeMethod: cmd.CodeChallengeMethod,
	})
	if err != nil {
		return errors.Wrap(err, "failed to save the user")
	}

	return nil
}

func (t *Controller) Get(ctx context.Context, cmd *GetCmd) (*AuthorizationCode, error) {
	err := validator.New().
		CheckString("code", cmd.Code, is.Required, is.StringInRange(8, 256)).
		Run()
	if err != nil {
		return nil, err
	}

	_, user, err := t.storage.Get(ctx, cmd.Code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the user")
	}

	return user, nil
}

func (t *Controller) Delete(ctx context.Context, cmd *DeleteCmd) error {
	err := validator.New().
		CheckString("code", cmd.Code, is.Required, is.StringInRange(8, 256)).
		Run()
	if err != nil {
		return err
	}

	rev, _, err := t.storage.Get(ctx, cmd.Code)
	if err != nil {
		return errors.Wrap(err, "failed to get the user")
	}

	if rev != "" {
		err := t.storage.Delete(ctx, cmd.Code, rev)
		if err != nil {
			return errors.Wrap(err, "failed to delete the user")
		}
	}

	return nil
}

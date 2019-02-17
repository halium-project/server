package accesstoken

import (
	"context"
	"log"
	"time"

	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/go-server-utils/password"
	"github.com/halium-project/go-server-utils/uuid"
	"github.com/halium-project/go-server-utils/validator"
	"github.com/halium-project/go-server-utils/validator/is"
	"github.com/halium-project/server/db"
	"gitlab.com/Peltoche/yaccc"
)

type Controller struct {
	uuid     uuid.Producer
	storage  StorageInterface
	password password.HashManager
}

type StorageInterface interface {
	Set(ctx context.Context, id string, rev string, value *AccessToken) (string, error)
	Get(ctx context.Context, code string) (string, *AccessToken, error)
	Delete(ctx context.Context, code string, rev string) error
	FindOneByRefreshToken(ctx context.Context, refreshToken string) (string, string, *AccessToken, error)
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
		CheckString("accessToken", cmd.AccessToken, is.Required, is.StringInRange(10, 50)).
		CheckString("refreshToken", cmd.RefreshToken, is.Optional, is.StringInRange(10, 50)).
		CheckNumber("expiresIn", cmd.ExpiresIn, is.Required, is.NumberPositif).
		CheckArray("scopes", cmd.Scopes, is.ArrayInRange(1, 50)).
		CheckEachString("scopes", cmd.Scopes, is.StringInRange(5, 120)).
		Run()
	if err != nil {
		return err
	}

	// Save the document
	_, err = t.storage.Set(ctx, cmd.AccessToken, "", &AccessToken{
		ClientID:     cmd.ClientID,
		AccessToken:  cmd.AccessToken,
		RefreshToken: cmd.RefreshToken,
		ExpiresIn:    cmd.ExpiresIn,
		Scopes:       cmd.Scopes,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "failed to save the accessToken")
	}

	return nil
}

func (t *Controller) Get(ctx context.Context, cmd *GetCmd) (*AccessToken, error) {
	err := validator.New().
		CheckString("accessToken", cmd.AccessToken, is.Required, is.StringInRange(8, 256)).
		Run()
	if err != nil {
		return nil, err
	}

	_, accessToken, err := t.storage.Get(ctx, cmd.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the accessToken")
	}

	return accessToken, nil
}

func (t *Controller) FindOneByRefreshToken(ctx context.Context, cmd *FindOneByRefreshTokenCmd) (string, *AccessToken, error) {
	err := validator.New().
		CheckString("refreshToken", cmd.RefreshToken, is.Required, is.StringInRange(8, 256)).
		Run()
	if err != nil {
		return "", nil, err
	}

	id, _, accessToken, err := t.storage.FindOneByRefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to get the accessToken")
	}

	return id, accessToken, nil
}

func (t *Controller) Delete(ctx context.Context, cmd *DeleteCmd) error {
	err := validator.New().
		CheckString("accessToken", cmd.AccessToken, is.Required, is.StringInRange(8, 256)).
		Run()
	if err != nil {
		return err
	}

	rev, _, err := t.storage.Get(ctx, cmd.AccessToken)
	if err != nil {
		return errors.Wrap(err, "failed to get the accessToken")
	}

	if rev != "" {
		err := t.storage.Delete(ctx, cmd.AccessToken, rev)
		if err != nil {
			return errors.Wrap(err, "failed to delete the accessToken")
		}
	}

	return nil
}

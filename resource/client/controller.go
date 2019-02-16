package client

import (
	"context"
	"log"
	"strings"
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
	password password.HashManager
	storage  StorageInterface
}

type StorageInterface interface {
	Set(ctx context.Context, id string, rev string, value *Client) (string, error)
	Get(ctx context.Context, id string) (string, *Client, error)
	GetAll(ctx context.Context) (map[string]Client, error)
	FindOneByName(ctx context.Context, name string) (string, string, *Client, error)
}

func InitController(ctx context.Context, server *yaccc.Server) *Controller {
	database, err := server.ConnectDatabase(ctx, BucketName)
	if err != nil {
		// SetupStorage will also create the default client.
		//
		// This is done on the storage lvl in order to force the id
		// ("00000000-0000-4200-b000-000000000000")
		database, err = SetupStorage(ctx, server)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "failed to setup %q storage", BucketName))
		}
	}

	storage := NewStorage(db.NewCouchdbDriver(database))

	uuidProducer := uuid.NewGoUUID()
	passwordProducer := password.NewPasswordHasher()

	return NewController(uuidProducer, passwordProducer, storage)
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

func (t *Controller) Create(ctx context.Context, cmd *CreateCmd) (string, string, error) {
	err := validator.New().
		CheckString("name", cmd.Name, is.Required, is.StringInRange(3, 50)).
		CheckArray("redirectURIs", cmd.RedirectURIs, is.ArrayInRange(0, 20)).
		CheckEachString("redirectURIs", cmd.RedirectURIs, is.URL).
		CheckArray("grantTypes", cmd.GrantTypes, is.ArrayInRange(1, 50)).
		CheckEachString("grantTypes", cmd.GrantTypes, is.OnOfString("client_credentials", "authorize_code", "implicit", "refresh_token", "password")).
		CheckEachString("responseTypes", cmd.ResponseTypes, is.OnOfString("code", "token")).
		CheckArray("scopes", cmd.Scopes, is.ArrayInRange(0, 50)).
		CheckEachString("scopes", cmd.Scopes, is.MatchingString(`[a-zA-Z0-9\.]+`), is.StringInRange(3, 50)).
		Run()
	if err != nil {
		return "", "", err
	}

	_, _, existingClient, err := t.storage.FindOneByName(ctx, cmd.Name)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to check if the name is already taken")
	}

	if existingClient != nil {
		return "", "", errors.NewValidationError().AddError("name", is.AlreadyUsed).IntoError()
	}

	id := strings.Replace(strings.ToLower(cmd.Name), " ", "-", -1)

	var hash string
	var secret string
	if !cmd.Public {
		secret = t.uuid.New()
		hash, err = t.password.Hash(secret)
		if err != nil {
			return "", "", errors.Wrap(err, "failed to hash password")
		}
	}

	client := Client{
		Secret:        hash,
		Name:          cmd.Name,
		RedirectURIs:  cmd.RedirectURIs,
		GrantTypes:    cmd.GrantTypes,
		ResponseTypes: cmd.ResponseTypes,
		Scopes:        cmd.Scopes,
		Public:        cmd.Public,
	}

	_, err = t.storage.Set(ctx, id, "", &client)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to save a client")
	}

	return id, secret, nil
}

func (t *Controller) Get(ctx context.Context, cmd *GetCmd) (*Client, error) {
	err := validator.New().
		CheckString("clientId", cmd.ClientID, is.Required, is.ID).
		Run()
	if err != nil {
		return nil, err
	}

	_, client, err := t.storage.Get(ctx, cmd.ClientID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a client")
	}

	return client, nil
}

func (t *Controller) GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]Client, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	res, err := t.storage.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all clients")
	}

	return res, nil
}

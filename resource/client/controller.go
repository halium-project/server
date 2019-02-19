package client

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
	password password.HashManager
	storage  StorageInterface
}

type StorageInterface interface {
	Set(ctx context.Context, id string, rev string, value *Client) (string, error)
	Get(ctx context.Context, id string) (string, *Client, error)
	GetAll(ctx context.Context) (map[string]Client, error)
	FindOneByName(ctx context.Context, name string) (string, string, *Client, error)
	Delete(ctx context.Context, id string) error
}

func InitController(ctx context.Context, server *yaccc.Server) *Controller {
	var requireBootstrap bool

	database, err := server.ConnectDatabase(ctx, BucketName)
	if err != nil {
		database, err = SetupStorage(ctx, server)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "failed to setup %q storage", BucketName))
		}

		requireBootstrap = true
	}

	storage := NewStorage(db.NewCouchdbDriver(database))

	uuidProducer := uuid.NewGoUUID()
	passwordProducer := password.NewPasswordHasher()

	controller := NewController(uuidProducer, passwordProducer, storage)

	// Give access to the "dashboard" app.
	//
	// This is required in order to configure your server.
	if requireBootstrap {
		_, _, err := controller.Create(ctx, &CreateCmd{
			ID:            "controle-panel",
			Name:          "Controle Panel",
			RedirectURIs:  []string{"http://localhost:8080"},
			GrantTypes:    []string{"implicit", "refresh_token"},
			ResponseTypes: []string{"token", "code"},
			Scopes:        []string{"users", "clients.read"},
			Public:        true,
		})
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to create the dashboard app permission"))
		}
	}

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

func (t *Controller) Create(ctx context.Context, cmd *CreateCmd) (string, string, error) {
	err := validator.New().
		CheckString("id", cmd.ID, is.Required, is.StringInRange(3, 50)).
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
		ID:            cmd.ID,
		Secret:        hash,
		Name:          cmd.Name,
		RedirectURIs:  cmd.RedirectURIs,
		GrantTypes:    cmd.GrantTypes,
		ResponseTypes: cmd.ResponseTypes,
		Scopes:        cmd.Scopes,
		Public:        cmd.Public,
	}

	_, err = t.storage.Set(ctx, cmd.ID, "", &client)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to save a client")
	}

	return cmd.ID, secret, nil
}

func (t *Controller) Get(ctx context.Context, cmd *GetCmd) (*Client, error) {
	err := validator.New().
		CheckString("clientID", cmd.ClientID, is.Required, is.StringInRange(3, 100)).
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

func (t *Controller) Delete(ctx context.Context, cmd *DeleteCmd) error {
	err := validator.New().
		CheckString("clientID", cmd.ClientID, is.Required, is.StringInRange(3, 100)).
		Run()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err = t.storage.Delete(ctx, cmd.ClientID)
	if err != nil {
		return errors.Wrap(err, "failed to delete the client")
	}

	return nil
}

package contact

import (
	"context"
	"log"
	"time"

	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/go-server-utils/uuid"
	"github.com/halium-project/go-server-utils/validator"
	"github.com/halium-project/go-server-utils/validator/is"
	"github.com/halium-project/server/db"
	"gitlab.com/Peltoche/yaccc"
)

type Controller struct {
	uuid    uuid.Producer
	storage StorageInterface
}

type StorageInterface interface {
	Set(ctx context.Context, id string, rev string, value *Contact) (string, error)
	Get(ctx context.Context, id string) (string, *Contact, error)
	GetAll(ctx context.Context) (map[string]Contact, error)
	FindOneByName(ctx context.Context, name string) (string, string, *Contact, error)
	Delete(ctx context.Context, id string) error
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

	controller := NewController(uuidProducer, storage)

	return controller
}

func NewController(
	uuid uuid.Producer,
	storage StorageInterface,
) *Controller {
	return &Controller{
		uuid:    uuid,
		storage: storage,
	}
}

func (t *Controller) Create(ctx context.Context, cmd *CreateCmd) (string, error) {
	err := validator.New().
		CheckString("name", cmd.Name, is.Required, is.StringInRange(3, 50)).
		Run()
	if err != nil {
		return "", err
	}

	_, _, existingContact, err := t.storage.FindOneByName(ctx, cmd.Name)
	if err != nil {
		return "", errors.Wrap(err, "failed to check if the name is already taken")
	}

	if existingContact != nil {
		return "", errors.NewValidationError().AddError("name", is.AlreadyUsed).IntoError()
	}

	contact := Contact{
		Name: cmd.Name,
	}

	// Generate the UUID
	contactID := t.uuid.New()

	_, err = t.storage.Set(ctx, contactID, "", &contact)
	if err != nil {
		return "", errors.Wrap(err, "failed to save a contact")
	}

	return contactID, nil
}

func (t *Controller) Get(ctx context.Context, cmd *GetCmd) (*Contact, error) {
	err := validator.New().
		CheckString("contactID", cmd.ContactID, is.Required, is.StringInRange(3, 100)).
		Run()
	if err != nil {
		return nil, err
	}

	_, contact, err := t.storage.Get(ctx, cmd.ContactID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a contact")
	}

	return contact, nil
}

func (t *Controller) GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]Contact, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	res, err := t.storage.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all contacts")
	}

	return res, nil
}

func (t *Controller) Delete(ctx context.Context, cmd *DeleteCmd) error {
	err := validator.New().
		CheckString("contactID", cmd.ContactID, is.Required, is.StringInRange(3, 100)).
		Run()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err = t.storage.Delete(ctx, cmd.ContactID)
	if err != nil {
		return errors.Wrap(err, "failed to delete the contact")
	}

	return nil
}

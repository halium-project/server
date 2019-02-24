package todo

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
	Set(ctx context.Context, id string, rev string, value *Todo) (string, error)
	Get(ctx context.Context, id string) (string, *Todo, error)
	GetAll(ctx context.Context) (map[string]Todo, error)
	FindOneByTitle(ctx context.Context, title string) (string, string, *Todo, error)
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
		CheckString("title", cmd.Title, is.Required, is.StringInRange(3, 50)).
		Run()
	if err != nil {
		return "", err
	}

	_, _, existingTodo, err := t.storage.FindOneByTitle(ctx, cmd.Title)
	if err != nil {
		return "", errors.Wrap(err, "failed to check if the title is already taken")
	}

	if existingTodo != nil {
		return "", errors.NewValidationError().AddError("title", is.AlreadyUsed).IntoError()
	}

	todo := Todo{
		Title: cmd.Title,
	}

	// Generate the UUID
	todoID := t.uuid.New()

	_, err = t.storage.Set(ctx, todoID, "", &todo)
	if err != nil {
		return "", errors.Wrap(err, "failed to save a todo")
	}

	return todoID, nil
}

func (t *Controller) Get(ctx context.Context, cmd *GetCmd) (*Todo, error) {
	err := validator.New().
		CheckString("todoID", cmd.TodoID, is.Required, is.StringInRange(3, 100)).
		Run()
	if err != nil {
		return nil, err
	}

	_, todo, err := t.storage.Get(ctx, cmd.TodoID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a todo")
	}

	return todo, nil
}

func (t *Controller) GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	res, err := t.storage.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all todos")
	}

	return res, nil
}

func (t *Controller) Delete(ctx context.Context, cmd *DeleteCmd) error {
	err := validator.New().
		CheckString("todoID", cmd.TodoID, is.Required, is.StringInRange(3, 100)).
		Run()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err = t.storage.Delete(ctx, cmd.TodoID)
	if err != nil {
		return errors.Wrap(err, "failed to delete the todo")
	}

	return nil
}

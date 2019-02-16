package user

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
	Set(ctx context.Context, userID string, rev string, value *User) (string, error)
	Get(ctx context.Context, userID string) (string, *User, error)
	GetAll(ctx context.Context) (map[string]User, error)
	FindOneByUsername(ctx context.Context, email string) (string, string, *User, error)
	FindTotalUserCount(ctx context.Context) (int, error)
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

	if requireBootstrap {
		_, err := controller.Create(ctx, &CreateCmd{
			Username: "admin",
			Password: "admin1234",
			Role:     Admin,
		})
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to create the admin user"))
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

func (t *Controller) Create(ctx context.Context, cmd *CreateCmd) (string, error) {
	err := validator.New().
		CheckString("username", cmd.Username, is.Required, is.StringInRange(4, 128)).
		CheckString("password", cmd.Password, is.Required, is.StringInRange(8, 256)).
		CheckString("role", cmd.Role, is.Required, is.OnOfString(
			Admin,
			Dev,
		)).
		Run()
	if err != nil {
		return "", err
	}

	err = t.validateUsernameUniqueness(ctx, cmd.Username)
	if err != nil {
		return "", err
	}

	// Generate the UUID
	userID := t.uuid.New()
	hash, salt, err := t.password.HashWithSalt(cmd.Password)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash the password")
	}

	// Save the document
	_, err = t.storage.Set(ctx, userID, "", &User{
		Username: cmd.Username,
		Role:     cmd.Role,
		Password: hash,
		Salt:     salt,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to save the user")
	}

	return userID, nil
}

func (t *Controller) Get(ctx context.Context, cmd *GetCmd) (*User, error) {
	err := validator.New().
		CheckString("userID", cmd.UserID, is.Required, is.ID).
		Run()
	if err != nil {
		return nil, err
	}

	_, user, err := t.storage.Get(ctx, cmd.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the user")
	}

	return user, nil
}

func (t *Controller) GetAll(ctx context.Context, cmd *GetAllCmd) (map[string]User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	res, err := t.storage.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all users")
	}

	return res, nil
}

func (t *Controller) Update(ctx context.Context, cmd *UpdateCmd) error {
	err := validator.New().
		CheckString("userID", cmd.UserID, is.Required, is.ID).
		CheckString("username", cmd.Username, is.Required, is.StringInRange(4, 128)).
		CheckString("role", cmd.Role, is.Required, is.OnOfString(
			Admin,
			Dev,
		)).
		Run()
	if err != nil {
		return err
	}

	rev, user, err := t.storage.Get(ctx, cmd.UserID)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the user")
	}

	if user == nil {
		return errors.New(errors.NotFound, "")
	}

	if user.Username != cmd.Username {
		err = t.validateUsernameUniqueness(ctx, cmd.Username)
		if err != nil {
			return err
		}

	}

	_, err = t.storage.Set(ctx, cmd.UserID, rev, &User{
		Username: cmd.Username,
		Role:     cmd.Role,

		// Don't touch this fields
		Password: user.Password,
		Salt:     user.Salt,
	})
	if err != nil {
		return errors.Wrap(err, "failed to save the user")
	}

	return nil
}

func (t *Controller) Validate(ctx context.Context, cmd *ValidateCmd) (string, *User, error) {
	userID, _, user, err := t.storage.FindOneByUsername(ctx, cmd.Username)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to get the credentials")
	}

	if user == nil {
		return "", nil, nil
	}

	valid, err := t.password.ValidateWithSalt(cmd.Password, user.Salt, user.Password)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to compare the password with the hash")
	}

	if !valid {
		return "", nil, nil
	}

	return userID, user, nil
}

func (t *Controller) GetTotalUserCount(ctx context.Context) (int, error) {
	nbUsers, err := t.storage.FindTotalUserCount(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve the number of user")
	}

	return nbUsers, nil
}

func (t *Controller) validateUsernameUniqueness(ctx context.Context, email string) error {
	_, _, user, err := t.storage.FindOneByUsername(ctx, email)
	if err != nil {
		return errors.Wrap(err, "failed to check if the user email is already taken")
	}

	if user != nil {
		return errors.NewValidationError().AddError("email", "ALREADY_USED").IntoError()
	}

	return nil
}

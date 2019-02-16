package authorizationcode

import (
	"context"

	"github.com/halium-project/server/db"
	"github.com/halium-project/server/util/errors"
	"gitlab.com/Peltoche/yaccc"
)

const BucketName = "authorization_codes"

type Storage struct {
	driver db.Driver
}

func SetupStorage(ctx context.Context, server *yaccc.Server) (*yaccc.Database, error) {
	db, err := server.CreateDatabase(ctx, &yaccc.CreateDatabaseCmd{
		Name: BucketName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the database")
	}

	return db, nil
}

func NewStorage(driver db.Driver) *Storage {
	return &Storage{
		driver: driver,
	}
}

func (t *Storage) Set(ctx context.Context, code string, rev string, value *AuthorizationCode) (string, error) {
	rev, err := t.driver.Set(ctx, code, rev, value)
	if err != nil {
		return "", errors.Wrap(err, "failed to set the document into the storage")
	}

	return rev, nil
}

func (t *Storage) Get(ctx context.Context, code string) (string, *AuthorizationCode, error) {
	var authorization AuthorizationCode

	rev, err := t.driver.Get(ctx, code, &authorization)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to get the document from the storage")
	}

	if rev == "" {
		return "", nil, nil
	}

	return rev, &authorization, nil
}

func (t *Storage) Delete(ctx context.Context, code string, rev string) error {
	err := t.driver.Delete(ctx, code, rev)
	if err != nil {
		return errors.Wrap(err, "failed to delete the document from the storage")
	}

	return nil
}

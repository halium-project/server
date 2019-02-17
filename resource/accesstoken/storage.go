package accesstoken

import (
	"context"

	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/server/db"
	"gitlab.com/Peltoche/yaccc"
)

const BucketName = "access_tokens"

type Storage struct {
	driver db.Driver
}

func SetupStorage(ctx context.Context, server *yaccc.Server) (*yaccc.Database, error) {
	db, err := server.CreateDatabase(ctx, &yaccc.CreateDatabaseCmd{
		Name: BucketName,
		DesignDocuments: map[string]yaccc.DesignDocument{
			"default": {
				Language: yaccc.Javascript,
				Views: map[string]yaccc.View{
					"by_refresh_token": {
						Map: `function (doc, meta) {
							if (doc.refreshToken) {
								emit(doc.refreshToken, null);
							}
						}`,
					},
				},
			},
		},
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

func (t *Storage) Set(ctx context.Context, code string, rev string, value *AccessToken) (string, error) {
	rev, err := t.driver.Set(ctx, code, rev, value)
	if err != nil {
		return "", errors.Wrap(err, "failed to set the document into the storage")
	}

	return rev, nil
}

func (t *Storage) Get(ctx context.Context, code string) (string, *AccessToken, error) {
	var authorization AccessToken

	rev, err := t.driver.Get(ctx, code, &authorization)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to get the document from the storage")
	}

	if rev == "" {
		return "", nil, nil
	}

	return rev, &authorization, nil
}

func (t *Storage) FindOneByRefreshToken(ctx context.Context, refreshToken string) (string, string, *AccessToken, error) {
	res, err := t.driver.ExecuteViewQuery(ctx, &db.Query{
		IndexName: "by_refresh_token",
		Limit:     1,
		Equals:    []interface{}{refreshToken},
	})
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to query the view")
	}

	if len(res) == 0 {
		return "", "", nil, nil
	}

	var accessToken AccessToken
	rev, err := t.driver.Get(ctx, res[0].ID, &accessToken)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to get the document")
	}

	return res[0].ID, rev, &accessToken, nil
}

func (t *Storage) Delete(ctx context.Context, code string, rev string) error {
	err := t.driver.Delete(ctx, code, rev)
	if err != nil {
		return errors.Wrap(err, "failed to delete the document from the storage")
	}

	return nil
}

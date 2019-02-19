package client

import (
	"context"

	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/server/db"
	"gitlab.com/Peltoche/yaccc"
)

const BucketName = "clients"

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
					"by_name": {
						Map: `function (doc, meta) {
							if (doc.name) {
								emit(doc.name, null);
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

func (t *Storage) Set(ctx context.Context, id string, rev string, value *Client) (string, error) {
	rev, err := t.driver.Set(ctx, id, rev, value)
	if err != nil {
		return "", errors.Wrap(err, "failed to set the document into the storage")
	}

	return rev, nil
}

func (t *Storage) Get(ctx context.Context, id string) (string, *Client, error) {
	var client Client

	rev, err := t.driver.Get(ctx, id, &client)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to get the document from the storage")
	}

	if rev == "" {
		return "", nil, nil
	}

	return rev, &client, nil
}

func (t *Storage) Delete(ctx context.Context, id string) error {
	var client Client

	rev, err := t.driver.Get(ctx, id, &client)
	if err != nil {
		return errors.Wrap(err, "failed to get the document from the storage")
	}

	if rev == "" {
		return nil
	}

	err = t.driver.Delete(ctx, id, rev)
	if err != nil {
		return errors.Wrap(err, "failed to delete the document from the storage")
	}

	return nil
}

func (t *Storage) GetAll(ctx context.Context) (map[string]Client, error) {
	viewResult, err := t.driver.ExecuteViewQuery(ctx, &db.Query{
		IndexName: "by_name",
		Limit:     200,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query the view")
	}

	if len(viewResult) == 0 {
		return map[string]Client{}, nil
	}

	clientList := map[string]*Client{}

	for _, val := range viewResult {
		clientList[val.ID] = &Client{}
	}

	err = t.driver.GetMany(ctx, clientList)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the documents")
	}

	res := make(map[string]Client, len(clientList))

	for key, value := range clientList {
		res[key] = *value
	}

	return res, nil
}

func (t *Storage) FindOneByName(ctx context.Context, name string) (string, string, *Client, error) {
	res, err := t.driver.ExecuteViewQuery(ctx, &db.Query{
		IndexName: "by_name",
		Limit:     1,
		Equals:    []interface{}{name},
	})
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to query the view")
	}

	if len(res) == 0 {
		return "", "", nil, nil
	}

	var client Client
	rev, err := t.driver.Get(ctx, res[0].ID, &client)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to get the document")
	}

	return res[0].ID, rev, &client, nil
}

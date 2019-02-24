package todo

import (
	"context"

	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/server/db"
	"gitlab.com/Peltoche/yaccc"
)

const BucketName = "todos"

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
					"by_title": {
						Map: `function (doc, meta) {
							if (doc.title) {
								emit(doc.title, null);
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

func (t *Storage) Set(ctx context.Context, id string, rev string, value *Todo) (string, error) {
	rev, err := t.driver.Set(ctx, id, rev, value)
	if err != nil {
		return "", errors.Wrap(err, "failed to set the document into the storage")
	}

	return rev, nil
}

func (t *Storage) Get(ctx context.Context, id string) (string, *Todo, error) {
	var todo Todo

	rev, err := t.driver.Get(ctx, id, &todo)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to get the document from the storage")
	}

	if rev == "" {
		return "", nil, nil
	}

	return rev, &todo, nil
}

func (t *Storage) Delete(ctx context.Context, id string) error {
	var todo Todo

	rev, err := t.driver.Get(ctx, id, &todo)
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

func (t *Storage) GetAll(ctx context.Context) (map[string]Todo, error) {
	viewResult, err := t.driver.ExecuteViewQuery(ctx, &db.Query{
		IndexName: "by_title",
		Limit:     200,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query the view")
	}

	if len(viewResult) == 0 {
		return map[string]Todo{}, nil
	}

	todoList := map[string]*Todo{}

	for _, val := range viewResult {
		todoList[val.ID] = &Todo{}
	}

	err = t.driver.GetMany(ctx, todoList)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the documents")
	}

	res := make(map[string]Todo, len(todoList))

	for key, value := range todoList {
		res[key] = *value
	}

	return res, nil
}

func (t *Storage) FindOneByTitle(ctx context.Context, title string) (string, string, *Todo, error) {
	res, err := t.driver.ExecuteViewQuery(ctx, &db.Query{
		IndexName: "by_title",
		Limit:     1,
		Equals:    []interface{}{title},
	})
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to query the view")
	}

	if len(res) == 0 {
		return "", "", nil, nil
	}

	var todo Todo
	rev, err := t.driver.Get(ctx, res[0].ID, &todo)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to get the document")
	}

	return res[0].ID, rev, &todo, nil
}

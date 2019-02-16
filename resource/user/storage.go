package user

import (
	"context"

	"github.com/halium-project/server/db"
	"github.com/halium-project/server/util/errors"
	"gitlab.com/Peltoche/yaccc"
)

const BucketName = "users"

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
					"by_username": {
						Map: `function (doc, meta) {
							if (doc.username) {
								emit(doc.username, null);
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

func (t *Storage) Set(ctx context.Context, id string, rev string, value *User) (string, error) {
	rev, err := t.driver.Set(ctx, id, rev, value)
	if err != nil {
		return "", errors.Wrap(err, "failed to set the document into the storage")
	}

	return rev, nil
}

func (t *Storage) Get(ctx context.Context, id string) (string, *User, error) {
	var user User

	rev, err := t.driver.Get(ctx, id, &user)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to get the document from the storage")
	}

	if rev == "" {
		return "", nil, nil
	}

	return rev, &user, nil
}

func (t *Storage) GetAll(ctx context.Context) (map[string]User, error) {
	viewResult, err := t.driver.ExecuteViewQuery(ctx, &db.Query{
		IndexName: "by_username",
		Limit:     200,
		Range: &db.Range{
			Start: []string{},
			End:   nil,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query the view")
	}

	if len(viewResult) == 0 {
		return map[string]User{}, nil
	}

	userList := map[string]*User{}

	for _, val := range viewResult {
		userList[val.ID] = &User{}
	}

	err = t.driver.GetMany(ctx, userList)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the documents")
	}

	res := make(map[string]User, len(userList))

	for key, value := range userList {
		res[key] = *value
	}

	return res, nil
}

func (t *Storage) FindOneByUsername(ctx context.Context, username string) (string, string, *User, error) {
	res, err := t.driver.ExecuteViewQuery(ctx, &db.Query{
		IndexName: "by_username",
		Limit:     1,
		Equals:    []interface{}{username},
	})
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to query the view")
	}

	if len(res) == 0 {
		return "", "", nil, nil
	}

	var user User
	rev, err := t.driver.Get(ctx, res[0].ID, &user)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to get the document")
	}

	return res[0].ID, rev, &user, nil
}

func (t *Storage) FindTotalUserCount(ctx context.Context) (int, error) {
	nbRows, err := t.driver.GetTotalRow(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to query the driver")
	}

	// The schema definitions are counted as row. Remove the number of schema
	// in order to keep only the number of user.
	nbRows--

	return nbRows, nil
}

package contact

import (
	"context"

	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/server/db"
	"gitlab.com/Peltoche/yaccc"
)

const BucketName = "contacts"

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

func (t *Storage) Set(ctx context.Context, id string, rev string, value *Contact) (string, error) {
	rev, err := t.driver.Set(ctx, id, rev, value)
	if err != nil {
		return "", errors.Wrap(err, "failed to set the document into the storage")
	}

	return rev, nil
}

func (t *Storage) Get(ctx context.Context, id string) (string, *Contact, error) {
	var contact Contact

	rev, err := t.driver.Get(ctx, id, &contact)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to get the document from the storage")
	}

	if rev == "" {
		return "", nil, nil
	}

	return rev, &contact, nil
}

func (t *Storage) Delete(ctx context.Context, id string) error {
	var contact Contact

	rev, err := t.driver.Get(ctx, id, &contact)
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

func (t *Storage) GetAll(ctx context.Context) (map[string]Contact, error) {
	viewResult, err := t.driver.ExecuteViewQuery(ctx, &db.Query{
		IndexName: "by_name",
		Limit:     200,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query the view")
	}

	if len(viewResult) == 0 {
		return map[string]Contact{}, nil
	}

	contactList := map[string]*Contact{}

	for _, val := range viewResult {
		contactList[val.ID] = &Contact{}
	}

	err = t.driver.GetMany(ctx, contactList)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the documents")
	}

	res := make(map[string]Contact, len(contactList))

	for key, value := range contactList {
		res[key] = *value
	}

	return res, nil
}

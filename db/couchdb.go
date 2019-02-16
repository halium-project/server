package db

import (
	"context"
	"reflect"
	"sync"

	"github.com/pkg/errors"
	"gitlab.com/Peltoche/yaccc"
)

type CouchdbDriver struct {
	bucket *yaccc.Database
}

func InitCouchdbServer(ctx context.Context, server *yaccc.Server) error {
	_, err := server.ConnectDatabase(ctx, "_users")
	if err != nil {
		_, err = server.CreateDatabase(context.Background(), &yaccc.CreateDatabaseCmd{
			Name: "_users",
		})
		if err != nil {
			return errors.Wrap(err, "failed to create the base _user table")
		}
	}

	return nil
}

func NewCouchdbDriver(bucket *yaccc.Database) *CouchdbDriver {
	return &CouchdbDriver{
		bucket: bucket,
	}
}

func (t *CouchdbDriver) Set(ctx context.Context, key string, rev string, value interface{}) (string, error) {
	return t.bucket.Set(ctx, key, rev, value)
}

func (t *CouchdbDriver) Get(ctx context.Context, key string, valuePtr interface{}) (string, error) {
	return t.bucket.Get(ctx, key, valuePtr)
}

func (t *CouchdbDriver) Delete(ctx context.Context, key string, rev string) error {
	return t.bucket.Delete(ctx, key, rev)
}

func (t *CouchdbDriver) ExecuteViewQuery(ctx context.Context, query *Query) ([]ViewRow, error) {
	yacccQuery := yaccc.ViewQuery{
		ViewName:       query.IndexName,
		DesignDocument: "default",
		Limit:          query.Limit,
		Keys:           query.Equals,
	}

	if query.Order == Descending {
		yacccQuery.Order = yaccc.Descending
	}

	if query.Range != nil {
		yacccQuery.Range = &yaccc.Range{
			Start: query.Range.Start,
			End:   query.Range.End,
		}
	}

	yacccRes, err := t.bucket.QueryView(ctx, &yacccQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query the yaccc view")
	}

	res := make([]ViewRow, len(yacccRes))
	for idx, yacccRow := range yacccRes {
		res[idx] = ViewRow{
			ID:    yacccRow.ID,
			Key:   yacccRow.Key,
			Value: yacccRow.Value,
		}
	}

	return res, nil
}

func (t *CouchdbDriver) GetMany(ctx context.Context, valueMapPtr interface{}) error {

	v := reflect.ValueOf(valueMapPtr)

	if v.Kind() != reflect.Map {
		return errors.New("valuesPtr is not a map[string]interface{}")
	}

	var wg sync.WaitGroup
	var commonErr error

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, key := range v.MapKeys() {
		wg.Add(1)

		strct := v.MapIndex(key)

		go func(key string) {
			defer wg.Done()

			_, err := t.Get(ctx, key, strct.Interface())
			if err != nil {
				if commonErr == nil {
					commonErr = errors.Wrap(err, "failed during bulk retrieve")
				}
				cancel()
			}
		}(key.String())
	}

	wg.Wait()

	return commonErr
}

func (t *CouchdbDriver) GetTotalRow(ctx context.Context) (int, error) {
	bucketInfo, err := t.bucket.Info(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to retrieve bucket info")
	}

	return bucketInfo.DocCount, nil
}

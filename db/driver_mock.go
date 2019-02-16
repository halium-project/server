package db

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type DriverMock struct {
	mock.Mock
}

func (t *DriverMock) Set(ctx context.Context, id string, rev string, value interface{}) (string, error) {
	args := t.Called(id, rev, value)

	return args.String(0), args.Error(1)
}

func (t *DriverMock) Get(ctx context.Context, key string, valuePtr interface{}) (string, error) {
	args := t.Called(key)

	if args.Get(1) == nil {
		return "", args.Error(2)
	}

	res, err := json.Marshal(args.Get(1))
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(res, valuePtr)
	if err != nil {
		return "", err
	}

	return args.String(0), args.Error(2)
}

func (t *DriverMock) Delete(ctx context.Context, id string, rev string) error {
	return t.Called(id, rev).Error(0)
}

func (t *DriverMock) ExecuteViewQuery(ctx context.Context, query *Query) ([]ViewRow, error) {
	args := t.Called(query)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]ViewRow), args.Error(1)
}

func (t *DriverMock) GetMany(ctx context.Context, valuesPtr interface{}) error {

	valueList := map[string]interface{}{}

	v := reflect.ValueOf(valuesPtr)

	if v.Kind() != reflect.Map {
		return errors.New("valuesPtr is not a map[string]interface{}")
	}

	keyList := make([]string, v.Len())
	for idx, key := range v.MapKeys() {
		keyList[idx] = key.String()

		strct := v.MapIndex(key)
		valueList[key.String()] = strct.Interface()
	}

	// Sort by name in order to force a fixed order between tests
	sort.Strings(keyList)
	args := t.Called(keyList)

	if args.Get(0) == nil {
		return args.Error(1)
	}

	ret := reflect.ValueOf(args.Get(0))

	for _, key := range ret.MapKeys() {
		strct := ret.MapIndex(key)

		res, err := json.Marshal(strct.Interface())
		if err != nil {
			return err
		}

		tmp := valueList[key.String()]
		err = json.Unmarshal(res, &tmp)
		if err != nil {
			return err
		}

	}

	return args.Error(1)
}

func (t *DriverMock) GetTotalRow(_ context.Context) (int, error) {
	args := t.Called()

	return args.Int(0), args.Error(1)
}

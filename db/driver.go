package db

import (
	"context"
	"encoding/json"
)

type Query struct {
	IndexName string
	Limit     uint
	Order     SortOrder
	Equals    []interface{}
	Range     *Range
}

type SortOrder int

type Range struct {
	Start interface{}
	End   interface{}
}

const (
	Ascending SortOrder = iota
	Descending
)

type ViewRow struct {
	ID    string
	Key   interface{}
	Value json.RawMessage
}

type Driver interface {
	Set(ctx context.Context, id string, rev string, value interface{}) (string, error)
	Delete(ctx context.Context, id string, rev string) error
	Get(ctx context.Context, key string, valuePtr interface{}) (string, error)
	GetMany(ctx context.Context, valuesPtr interface{}) error
	ExecuteViewQuery(ctx context.Context, query *Query) ([]ViewRow, error)
	GetTotalRow(ctx context.Context) (int, error)
}

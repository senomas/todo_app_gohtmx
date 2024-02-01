package store

import (
	"encoding/json"
	"errors"
	"log/slog"
)

type key int

const (
	StoreCtxDB key = iota
	StoreCtxCookie
)

var ErrNoData = errors.New("sql: no data")

type JsonLogValue struct {
	V []any
}

func (a *JsonLogValue) LogValue() slog.Value {
	v, _ := json.Marshal(a.V)
	return slog.StringValue(string(v))
}

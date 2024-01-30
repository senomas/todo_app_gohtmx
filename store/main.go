package store

import "errors"

type key int

const StoreCtxDB key = iota

var ErrNoData = errors.New("sql: no data")

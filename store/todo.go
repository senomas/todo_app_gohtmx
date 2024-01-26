package store

type Todo struct {
	Title     string
	ID        int64
	Completed bool
}

type TodoFilter struct {
	Title     FilterString
	ID        FilterInt64
	Completed FilterBool
}

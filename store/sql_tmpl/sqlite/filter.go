package sqlite

import (
	"fmt"

	"github.com/senomas/todo_app/store"
)

func filterToString(where []string, args []any, field string, filter any) ([]string, []any) {
	switch f := filter.(type) {
	case store.FilterInt64:
		switch f.Op {
		case store.OP_NOP:
		case store.OP_EQ:
			args = append(args, f.Value)
			where = append(where, fmt.Sprintf("%s = $%d", field, len(args)))
		default:
			panic(fmt.Sprintf("unknown filter int64 op: %+v", f))
		}
	case store.FilterBool:
		switch f.Op {
		case store.OP_NOP:
		case store.OP_EQ:
			args = append(args, f.Value)
			where = append(where, fmt.Sprintf("%s = $%d", field, len(args)))
		default:
			panic(fmt.Sprintf("unknown filter bool op: %+v", f))
		}
	case store.FilterString:
		switch f.Op {
		case store.OP_NOP:
		case store.OP_EQ:
			args = append(args, f.Value)
			where = append(where, fmt.Sprintf("%s = $%d", field, len(args)))
		case store.OP_LIKE:
			args = append(args, f.Value)
			where = append(where, fmt.Sprintf("%s LIKE $%d", field, len(args)))
		default:
			panic(fmt.Sprintf("unknown filter string op: %+v", f))
		}
	default:
		panic(fmt.Sprintf("unknown filter type: %+v", f))
	}
	return where, args
}

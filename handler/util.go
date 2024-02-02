package handler

import "fmt"

var f = fmt.Sprintf

func fv(v any) string {
	return fmt.Sprintf("%v", v)
}

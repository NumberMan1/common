package example

import "cmp"

type Member[T, U cmp.Ordered] struct {
	UId   U
	Score T
	Extra any
}

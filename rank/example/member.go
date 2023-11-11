package example

import "golang.org/x/exp/constraints"

type Member[T, U constraints.Ordered] struct {
	UId   U
	Score T
	Extra any
}

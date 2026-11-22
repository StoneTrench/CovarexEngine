package math

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Float | constraints.Integer
}

func Abs[T Number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
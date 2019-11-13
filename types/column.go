package types

import "reflect"

type Column struct {
	Name  string
	Skip  bool
	ColNo int
	Kind  reflect.Kind
	Label string
}

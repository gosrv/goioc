package gioc

import "reflect"

const (
	PrioritySystem = 0
	PriorityHigh   = 10000
	PriorityMiddle = 100000
	PriorityLow    = 1000000
)

type IPriority interface {
	GetPriority() int
}

var IPriorityType = reflect.TypeOf((*IPriority)(nil)).Elem()

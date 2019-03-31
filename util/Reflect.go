package util

import "reflect"

type rfl struct {}

var Rfl = &rfl{}

func (this *rfl)UnptrValue(ins interface{}) reflect.Value {
	val := reflect.ValueOf(ins)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}

func (this *rfl)UnptrType(ins interface{}) reflect.Type {
	val := reflect.TypeOf(ins)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}


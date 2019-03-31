package util

import (
	"fmt"
	"log"
	"reflect"
)

func Assert(condition bool, msg string) {
	if !condition {
		Panic(msg)
	}
}

func VerifyNotNull(ins interface{}) {
	if ins == nil {
		Panic("nil interface")
	}
}

func VerifyNoError(err error) {
	if err != nil {
		Panic(err.Error())
	}
}

func Panic(format string, a ...interface{}) {
	log.Panic(fmt.Sprintf(format, a...))
}

func IsNilAnonymousField(ins interface{}, ftype reflect.Type) bool {
	bcValue := reflect.ValueOf(ins)
	bcType := reflect.TypeOf(ins)

	for bcType.Kind() == reflect.Ptr {
		if bcValue.IsNil() {
			return false
		}
		bcValue = bcValue.Elem()
		bcType = bcType.Elem()
	}

	for i := 0; i < bcValue.NumField(); i++ {
		fieldValue := bcValue.Field(i)
		fieldType := bcType.Field(i)
		if !fieldType.Anonymous {
			continue
		}
		if fieldType.Type.AssignableTo(ftype) {
			fieldValue = Hack.ValuePatchWrite(fieldValue)
			return fieldValue.IsNil()
		}
	}
	return false
}

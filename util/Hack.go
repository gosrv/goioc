package util

import (
	"reflect"
	"unsafe"
)

const flagRO = uintptr((1 << 5) | (1 << 6))

var Hack = struct {
	ValuePatchWrite           func(val reflect.Value) reflect.Value
	ValuePatchWriteRecurrence func(val reflect.Value) reflect.Value
}{
	ValuePatchWrite: valuePatchWrite,
}

func valuePatchWrite(val reflect.Value) reflect.Value {
	if val.CanSet() {
		return val
	}
	// reflect reflect.Value
	rval := reflect.ValueOf(&val)
	flagField := rval.Elem().FieldByName("flag")
	// fetch Value.flag field.
	ptrField := unsafe.Pointer(flagField.UnsafeAddr())
	// cast to pointer.
	flagPtr := (*uintptr)(ptrField)
	// modify flag
	*flagPtr = *flagPtr &^ flagRO
	return val
}

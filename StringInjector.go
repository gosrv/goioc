package goioc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gosrv/util"
	"reflect"
	"strconv"
)

type PropertyInjectFunc struct {
	injectors map[reflect.Kind]func(value reflect.Value, val string) error
}

func (this *PropertyInjectFunc) StringValueInjector(value reflect.Value, val string) error {
	for value.Kind() == reflect.Ptr {
		value.Set(reflect.New(value.Type().Elem()))
		value = value.Elem()
	}
	return this.injectors[value.Kind()](value, val)
}

type IStringValueInjector interface {
	StringValueInjector(value reflect.Value, val string) error
}

func (this *PropertyInjectFunc) Inject(value reflect.Value, val string) error {
	rkind := value.Type().Kind()
	injfunc, ok := this.injectors[rkind]
	if !ok {
		return errors.New(fmt.Sprintf("no injector for type: %s", rkind))
	}

	return injfunc(value, val)
}

func BoolInject(value reflect.Value, val string) error {
	pv, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	value.SetBool(pv)
	return nil
}

func IntInject(value reflect.Value, val string) error {
	pv, err := strconv.ParseInt(val, 0, 64)
	if err != nil {
		return err
	}
	value.SetInt(pv)
	return nil
}

func FloatInject(value reflect.Value, val string) error {
	pv, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	value.SetFloat(pv)
	return nil
}

func UintInject(value reflect.Value, val string) error {
	pv, err := strconv.ParseUint(val, 0, 64)
	if err != nil {
		return err
	}
	value.SetUint(pv)
	return nil
}

func StringInject(value reflect.Value, val string) error {
	value.SetString(val)
	return nil
}

func SliceInject(value reflect.Value, val string) error {
	return errors.New("slice value inject not implement")
}

func StructInject(value reflect.Value, val string) error {
	util.Hack.ValuePatchWrite(value)
	return json.Unmarshal([]byte(val), value.Addr().Interface())
}

func PtrInject(value reflect.Value, val string) error {
	return errors.New("ptr value inject not implement")
}

var StringPropertyInjects = &PropertyInjectFunc{
	map[reflect.Kind]func(value reflect.Value, val string) error{
		reflect.Bool:    BoolInject,
		reflect.Int:     IntInject,
		reflect.Int8:    IntInject,
		reflect.Int16:   IntInject,
		reflect.Int32:   IntInject,
		reflect.Int64:   IntInject,
		reflect.Uint:    UintInject,
		reflect.Uint8:   UintInject,
		reflect.Uint16:  UintInject,
		reflect.Uint32:  UintInject,
		reflect.Uint64:  UintInject,
		reflect.Float32: FloatInject,
		reflect.Float64: FloatInject,
		reflect.Slice:   SliceInject,
		reflect.String:  StringInject,
		reflect.Struct:  StructInject,
		reflect.Ptr:     PtrInject,
	},
}

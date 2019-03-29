package goioc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gosrv/goioc/util"
	"reflect"
	"strconv"
)

type PropertyInjectGroup struct {
	injectors map[reflect.Kind]IStringValueInjector
}

func (this *PropertyInjectGroup) StringValueInjector(value reflect.Value, val string) error {
	for value.Kind() == reflect.Ptr {
		value.Set(reflect.New(value.Type().Elem()))
		value = value.Elem()
	}
	return this.injectors[value.Kind()].StringValueInjector(value, val)
}

type IStringValueInjector interface {
	StringValueInjector(value reflect.Value, val string) error
}
type FuncStringValueInjector func(value reflect.Value, val string) error

func (this FuncStringValueInjector) StringValueInjector(value reflect.Value, val string) error {
	return this(value, val)
}

func (this *PropertyInjectGroup) Inject(value reflect.Value, val string) error {
	rkind := value.Type().Kind()
	injfunc, ok := this.injectors[rkind]
	if !ok {
		return errors.New(fmt.Sprintf("no injector for type: %s", rkind))
	}

	return injfunc.StringValueInjector(value, val)
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

var StringPropertyInjects = &PropertyInjectGroup{
	map[reflect.Kind]IStringValueInjector{
		reflect.Bool:    FuncStringValueInjector(BoolInject),
		reflect.Int:     FuncStringValueInjector(IntInject),
		reflect.Int8:    FuncStringValueInjector(IntInject),
		reflect.Int16:   FuncStringValueInjector(IntInject),
		reflect.Int32:   FuncStringValueInjector(IntInject),
		reflect.Int64:   FuncStringValueInjector(IntInject),
		reflect.Uint:    FuncStringValueInjector(UintInject),
		reflect.Uint8:   FuncStringValueInjector(UintInject),
		reflect.Uint16:  FuncStringValueInjector(UintInject),
		reflect.Uint32:  FuncStringValueInjector(UintInject),
		reflect.Uint64:  FuncStringValueInjector(UintInject),
		reflect.Float32: FuncStringValueInjector(FloatInject),
		reflect.Float64: FuncStringValueInjector(FloatInject),
		reflect.Slice:   FuncStringValueInjector(SliceInject),
		reflect.String:  FuncStringValueInjector(StringInject),
		reflect.Struct:  FuncStringValueInjector(StructInject),
		reflect.Ptr:     FuncStringValueInjector(PtrInject),
	},
}

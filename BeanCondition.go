package gioc

import (
	"github.com/gosrv/goioc/util"
	"reflect"
)

const (
	beanInjectTag = "bean"
)

/**
bean的装配条件，只有IsConditionPass才会进入装配，并放入到BeanContainer
*/
type IBeanCondition interface {
	IsConditionPass() bool
}

var BeanConditionType = reflect.TypeOf((*IBeanCondition)(nil)).Elem()
var BeanConditionArrayType = reflect.TypeOf((*[]IBeanCondition)(nil)).Elem()

// and 所有条件通过才算通过
type conditionAnd struct {
	conditions []IBeanCondition
}

func (this *conditionAnd) IsConditionPass() bool {
	for _, condion := range this.conditions {
		if !condion.IsConditionPass() {
			return false
		}
	}
	return true
}

func NewConditionAnd(conditions ...IBeanCondition) IBeanCondition {
	return &conditionAnd{
		conditions: conditions,
	}
}

// or 有一个条件通过就算通过
type conditionOr struct {
	conditions []IBeanCondition
}

func (this *conditionOr) IsConditionPass() bool {
	for _, condion := range this.conditions {
		if condion.IsConditionPass() {
			return true
		}
	}
	return false
}

func NewConditionOr(conditions ...IBeanCondition) IBeanCondition {
	return &conditionOr{
		conditions: conditions,
	}
}

type conditionNot struct {
	condition IBeanCondition
}

func (this *conditionNot) IsConditionPass() bool {
	return !this.condition.IsConditionPass()
}

func NewConditionNot(condition IBeanCondition) IBeanCondition {
	return &conditionNot{
		condition: condition,
	}
}

// 配置文件中存在配置项或不存在配置项
type conditionOnValue struct {
	value  string
	exist  bool
	loader IConfigLoader `bean`
}

func (this *conditionOnValue) IsConditionPass() bool {
	if this.loader == nil {
		return false
	}
	_, ok := this.loader.Config().Map()[this.value]
	return (!ok && !this.exist) || (ok && this.exist)
}

func NewConditionOnValue(value string, exist bool) IBeanCondition {
	return &conditionOnValue{
		value: value,
		exist: exist,
	}
}

// 已经有bean存在或不存在
type conditionOnBeanType struct {
	beanType      reflect.Type
	exit          bool
	beanContainer IBeanContainer `bean`
}

func (this *conditionOnBeanType) IsConditionPass() bool {
	if this.beanContainer == nil {
		return false
	}
	realExist := len(this.beanContainer.GetBeanByType(this.beanType)) > 0
	return (this.exit && realExist) || (!this.exit && !realExist)
}

func NewConditionOnBeanType(beanType reflect.Type, exit bool) *conditionOnBeanType {
	return &conditionOnBeanType{
		beanType: beanType,
		exit:     exit,
	}
}

type conditionOnBeanName struct {
	beanName      string
	exit          bool
	beanContainer IBeanContainer `bean`
}

func (this *conditionOnBeanName) IsConditionPass() bool {
	if this.beanContainer == nil {
		return false
	}
	realExist := this.beanContainer.GetBeanByName(this.beanName) != nil
	return (this.exit && realExist) || (!this.exit && !realExist)
}

func NewConditionOnBeanName(beanName string, exit bool) *conditionOnBeanName {
	return &conditionOnBeanName{
		beanName: beanName,
		exit:     exit,
	}
}

type IBeanConditionInjector interface {
	BeanConditionInjector(bcondition IBeanCondition, beanContainer IBeanContainer)
}

type BeanBeanConditionInjector struct {
}

func NewBeanBeanConditionInjector() *BeanBeanConditionInjector {
	return &BeanBeanConditionInjector{}
}
func (this *BeanBeanConditionInjector) BeanConditionInjector(bcondition IBeanCondition, beanContainer IBeanContainer) {
	beanConditionInjectorLoopCheck(bcondition, beanContainer, 1024)
}

func beanConditionInjectorLoopCheck(bcondition IBeanCondition,
	beanContainer IBeanContainer, leftDepth int) {
	if leftDepth < 0 {
		util.Panic("loop in condition inject %v", reflect.TypeOf(bcondition))
	}
	bcValue := reflect.ValueOf(bcondition)
	bcType := reflect.TypeOf(bcondition)

	for bcType.Kind() == reflect.Ptr {
		if bcValue.IsNil() {
			return
		}
		bcValue = bcValue.Elem()
		bcType = bcType.Elem()
	}

	for i := 0; i < bcValue.NumField(); i++ {
		fieldValue := bcValue.Field(i)
		fieldType := bcType.Field(i)

		if fieldType.Type.AssignableTo(BeanConditionType) {
			fieldValue = util.Hack.ValuePatchWrite(fieldValue)
			beanConditionInjectorLoopCheck(fieldValue.Interface().(IBeanCondition), beanContainer, leftDepth-1)
			continue
		}
		if fieldType.Type.AssignableTo(BeanConditionArrayType) {
			fieldValue = util.Hack.ValuePatchWrite(fieldValue)
			for _, fieldCondition := range fieldValue.Interface().([]IBeanCondition) {
				beanConditionInjectorLoopCheck(fieldCondition, beanContainer, leftDepth-1)
			}
			continue
		}
		if fieldType.Tag != beanInjectTag {
			continue
		}
		beans := beanContainer.GetBeanByType(fieldType.Type)
		switch len(beans) {
		case 0:
			continue
		case 1:
			fieldValue = util.Hack.ValuePatchWrite(fieldValue)
			fieldValue.Set(reflect.ValueOf(beans[0]))
		default:
			util.Panic("error")
		}
	}
}

var BeanConditionInjectorHelper = struct {
	GetBeanConditionInjector func(beanContainer IBeanContainer) []IBeanConditionInjector
	BeanConditionInjector    func(bcondition IBeanCondition, injecots []IBeanConditionInjector, beanContainer IBeanContainer)
}{
	GetBeanConditionInjector: getBeanConditionInjector,
	BeanConditionInjector:    beanConditionInjector,
}

func getBeanConditionInjector(beanContainer IBeanContainer) []IBeanConditionInjector {
	injectorsInss := beanContainer.GetBeanByType(reflect.TypeOf((*IBeanConditionInjector)(nil)).Elem())
	injectors := make([]IBeanConditionInjector, len(injectorsInss), len(injectorsInss))
	for i := 0; i < len(injectorsInss); i++ {
		injectors[i] = injectorsInss[i].(IBeanConditionInjector)
	}
	return injectors
}

func beanConditionInjector(bcondition IBeanCondition, injecots []IBeanConditionInjector, beanContainer IBeanContainer) {
	for _, injector := range injecots {
		injector.BeanConditionInjector(bcondition, beanContainer)
	}
}

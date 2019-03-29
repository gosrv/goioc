package gioc

import (
	"github.com/gosrv/goioc/util"
	"math"
	"reflect"
	"sort"
)

type beanField struct {
	tags        map[string]string
	fieldValue  reflect.Value
	fieldType   reflect.Type
	structField *reflect.StructField
}

func newBeanField(fieldValue reflect.Value, structField *reflect.StructField) *beanField {
	return &beanField{
		fieldValue:  fieldValue,
		fieldType:   structField.Type,
		structField: structField,
	}
}

type beanCtx struct {
	beanIns   interface{}
	condition IBeanCondition
	name      string

	beanValue reflect.Value
	beanType  reflect.Type

	eleValue reflect.Value
	eleType  reflect.Type

	fields []*beanField
}

func newBeanCtx(name string, beanIns interface{}) *beanCtx {
	ctx := &beanCtx{
		name:      name,
		beanIns:   beanIns,
		beanValue: reflect.ValueOf(beanIns),
		beanType:  reflect.TypeOf(beanIns),
	}

	ctx.eleValue = ctx.beanValue
	ctx.eleType = ctx.beanType

	for ctx.eleType.Kind() == reflect.Ptr {
		ctx.eleValue = ctx.eleValue.Elem()
		ctx.eleType = ctx.eleType.Elem()
	}

	if ctx.eleType.Kind() != reflect.Struct {
		return ctx
	}

	for i := 0; i < ctx.eleValue.NumField(); i++ {
		fieldType := ctx.eleType.Field(i)
		fieldValue := ctx.eleValue.Field(i)
		field := newBeanField(fieldValue, &fieldType)
		ctx.fields = append(ctx.fields, field)
		if fieldType.Anonymous {
			if fieldType.Type.AssignableTo(BeanConditionType) {
				if ctx.condition != nil {
					util.Panic("duplicate bean condition in %v", ctx.beanType)
				}
				ctx.condition = fieldValue.Interface().(IBeanCondition)
			}
			if fieldType.Type.AssignableTo(IBeanNameType) {
				if len(ctx.name) > 0 {
					util.Panic("duplicate bean name in %v", ctx.beanType)
				}
				ctx.name = fieldValue.Interface().(IBeanName).BeanName()
			}
		}
	}

	if ctx.condition == nil && ctx.beanType.AssignableTo(BeanConditionType) {
		ctx.condition = beanIns.(IBeanCondition)
	}

	return ctx
}

type IBeanContainerBuilder interface {
	AddBean(bean ...interface{})
	AddNamedBean(name string, bean interface{})
	GetBeanContainer() IBeanContainer
	Build()
}

type defaultBeanContainerBuilder struct {
	beans         []*beanCtx
	beanContainer IBeanContainer
}

func NewBeanContainerBuilder() IBeanContainerBuilder {
	return &defaultBeanContainerBuilder{
		beanContainer: NewBeanContainer(),
	}
}

func (this *defaultBeanContainerBuilder) GetBeanContainer() IBeanContainer {
	return this.beanContainer
}

func (this *defaultBeanContainerBuilder) doAddBean(bean *beanCtx) {
	for _, cbean := range this.beans {
		if bean.beanIns == cbean.beanIns {
			util.Panic("duplicate bean instance")
		}
	}
	this.beans = append(this.beans, bean)
}

func (this *defaultBeanContainerBuilder) AddNamedBean(name string, bean interface{}) {
	util.VerifyNotNull(bean)
	this.doAddBean(newBeanCtx(name, bean))
}

func (this *defaultBeanContainerBuilder) AddBean(beans ...interface{}) {
	for _, bean := range beans {
		util.VerifyNotNull(bean)
		this.doAddBean(newBeanCtx("", bean))
	}
}

/**
build 步骤
1. 通过condition过滤所有不符合条件的bean，通过条件的bean加入beanContainer
2. IBeanBeforeTagProcess
3. 获取所有的TagProcessor，并用其处理所有的bean
4. IBeanAfterTagProcess
*/
func (this *defaultBeanContainerBuilder) Build() {
	this.buildCondition()

	tagProcessors := TagProcessorHelper.GetTagProcessor(this.beanContainer)
	if len(tagProcessors) == 0 {
		return
	}
	// nil check
	for _, tp := range tagProcessors {
		if reflect.TypeOf(tp).AssignableTo(IPriorityType) {
			if util.IsNilAnonymousField(tp, IPriorityType) {
				util.Panic("nil IPriority interface in bean %v", reflect.TypeOf(tp))
			}
		}
	}
	sort.Slice(tagProcessors, func(i, j int) bool {
		ip := math.MaxInt32
		jp := math.MaxInt32
		if reflect.TypeOf(tagProcessors[i]).AssignableTo(IPriorityType) {
			priority := tagProcessors[i].(IPriority)
			if priority == nil {
				util.Panic("nil interface %v", reflect.TypeOf(tagProcessors[i]))
			}
			ip = tagProcessors[i].(IPriority).GetPriority()
		}
		if reflect.TypeOf(tagProcessors[j]).AssignableTo(IPriorityType) {
			jp = tagProcessors[j].(IPriority).GetPriority()
		}
		return ip < jp
	})

	tagParser := TagParserHelper.GetTagParser(this.beanContainer)
	if tagParser == nil {
		util.Panic("no tag parser found!!!, you must add exace 1 tag parser")
	}

	for _, tagProcessor := range tagProcessors {
		for _, bean := range this.beanContainer.GetAllBeans() {
			var process IBeanProcess = nil
			if reflect.TypeOf(bean).AssignableTo(IBeanProcessType) {
				process = bean.(IBeanProcess)
			}
			if process != nil {
				process.BeanBeforeTagProcess(tagProcessor, this.beanContainer)
			}
			TagProcessorHelper.BeanTagProcess(bean, tagParser, tagProcessor)
			if process != nil {
				process.BeanAfterTagProcess(tagProcessor, this.beanContainer)
			}
		}
	}
}

func (this *defaultBeanContainerBuilder) buildCondition() {
	// 反复收集可用bean
	isSomeThingDo := true
	leftBeans := this.beans

	for isSomeThingDo {
		isSomeThingDo = false
		for i := 0; i < 2; i++ {
			// condition check
			processBeans := leftBeans
			leftBeans = make([]*beanCtx, 0, len(processBeans))
			conditionInjectors := BeanConditionInjectorHelper.GetBeanConditionInjector(this.beanContainer)

			for _, bean := range processBeans {
				if bean.condition == nil || bean.condition.IsConditionPass() {
					if len(bean.name) > 0 {
						this.beanContainer.AddNamedBean(bean.name, bean.beanIns)
					} else {
						this.beanContainer.AddBean(bean.beanIns)
					}
					isSomeThingDo = true
				} else {
					BeanConditionInjectorHelper.BeanConditionInjector(bean.condition,
						conditionInjectors, this.beanContainer)
					leftBeans = append(leftBeans, bean)
				}
			}
		}
	}
}

func (this *defaultBeanContainerBuilder) BeanTagProcess(bean interface{}) {
	tagProcessors := TagProcessorHelper.GetTagProcessor(this.beanContainer)
	tagParser := TagParserHelper.GetTagParser(this.beanContainer)

	TagProcessorHelper.BeanTagProcess(bean, tagParser, tagProcessors...)
}

package gioc

import (
	"github.com/gosrv/goioc/util"
	"reflect"
)

type ITagProcessor interface {
	TagProcessorName() string
	TagProcess(bean interface{}, field reflect.Value, tags map[string]string)
}

var ITagProcessorType = reflect.TypeOf((*ITagProcessor)(nil)).Elem()

type ITagProcessorPriority interface {
	GetTagProcessorPriority() int
}

var ITagProcessorPriorityType = reflect.TypeOf((*ITagProcessorPriority)(nil)).Elem()

type tagProcessorHelper struct{}

var TagProcessorHelper = &tagProcessorHelper{}

func (this *tagProcessorHelper) GetTagProcessor(beanContainer IBeanContainer) []ITagProcessor {
	tagProcessorInss := beanContainer.GetBeanByType(reflect.TypeOf((*ITagProcessor)(nil)).Elem())
	tagProcessors := make([]ITagProcessor, len(tagProcessorInss), len(tagProcessorInss))
	for i := 0; i < len(tagProcessorInss); i++ {
		tagProcessors[i] = tagProcessorInss[i].(ITagProcessor)
	}
	return tagProcessors
}

func (this *tagProcessorHelper) BeanTagProcess(bean interface{}, parser ITagParser, processors ...ITagProcessor) {
	beanVal := util.Rfl.UnptrValue(bean)
	beanType := util.Rfl.UnptrType(bean)
	if beanType.Kind() != reflect.Struct {
		return
	}
	beanVal = util.Hack.ValuePatchWrite(beanVal)

	for i := 0; i < beanVal.NumField(); i++ {
		fieldValue := beanVal.Field(i)
		fieldType := beanType.Field(i)
		if len(fieldType.Tag) == 0 {
			continue
		}
		tags := parser.Parse(fieldType.Tag)
		fieldValue = util.Hack.ValuePatchWrite(fieldValue)
		for _, tagProcessor := range processors {
			tagProcessor.TagProcess(bean, fieldValue, tags)
		}
	}
}

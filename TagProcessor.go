package goioc

import (
	"github.com/gosrv/goioc/util"
	"reflect"
)

type ITagProcessor interface {
	TagProcessorName() string
	TagProcess(bean interface{}, field reflect.Value, tags map[string]string)
}

var TagProcessorHelper = struct {
	GetTagProcessor func(beanContainer IBeanContainer) []ITagProcessor
	BeanTagProcess  func(bean interface{}, parser ITagParser, processors ...ITagProcessor)
}{
	GetTagProcessor: getTagProcessor,
	BeanTagProcess:  beanTagProcess,
}

func getTagProcessor(beanContainer IBeanContainer) []ITagProcessor {
	tagProcessorInss := beanContainer.GetBeanByType(reflect.TypeOf((*ITagProcessor)(nil)).Elem())
	tagProcessors := make([]ITagProcessor, len(tagProcessorInss), len(tagProcessorInss))
	for i := 0; i < len(tagProcessorInss); i++ {
		tagProcessors[i] = tagProcessorInss[i].(ITagProcessor)
	}
	return tagProcessors
}

func beanTagProcess(bean interface{}, parser ITagParser, processors ...ITagProcessor) {
	beanVal := reflect.ValueOf(bean)
	beanType := reflect.TypeOf(bean)

	for beanVal.Kind() == reflect.Ptr {
		beanVal = beanVal.Elem()
		beanType = beanType.Elem()
	}
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

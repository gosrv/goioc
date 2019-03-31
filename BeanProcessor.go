package gioc

import "reflect"

type IBeanProcess interface {
	// tag前置处理
	BeanBeforeTagProcess(tagProcessor ITagProcessor, beanContainer IBeanContainer)
	// tag后缀处理
	BeanAfterTagProcess(tagProcessor ITagProcessor, beanContainer IBeanContainer)
}

var IBeanProcessType = reflect.TypeOf((*IBeanProcess)(nil)).Elem()

type beanProcessHelper struct{}

func (this *beanProcessHelper) GetBeanProcess(beanContainer IBeanContainer) []IBeanProcess {
	tagProcessorInss := beanContainer.GetBeanByType(reflect.TypeOf((*IBeanProcess)(nil)).Elem())
	tagProcessors := make([]IBeanProcess, len(tagProcessorInss), len(tagProcessorInss))
	for i := 0; i < len(tagProcessorInss); i++ {
		tagProcessors[i] = tagProcessorInss[i].(IBeanProcess)
	}
	return tagProcessors
}

var BeanProcessHelper = &beanProcessHelper{}

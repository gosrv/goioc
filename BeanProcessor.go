package gioc

import "reflect"

type IBeanProcess interface {
	// tag前置处理
	BeanBeforeTagProcess(tagProcessor ITagProcessor, beanContainer IBeanContainer)
	// tag后缀处理
	BeanAfterTagProcess(tagProcessor ITagProcessor, beanContainer IBeanContainer)
}

var IBeanProcessType = reflect.TypeOf((*IBeanProcess)(nil)).Elem()

var BeanProcessHelper = struct {
	GetBeanProcess func(beanContainer IBeanContainer) []IBeanProcess
}{
	GetBeanProcess: func(beanContainer IBeanContainer) []IBeanProcess {
		tagProcessorInss := beanContainer.GetBeanByType(reflect.TypeOf((*IBeanProcess)(nil)).Elem())
		tagProcessors := make([]IBeanProcess, len(tagProcessorInss), len(tagProcessorInss))
		for i := 0; i < len(tagProcessorInss); i++ {
			tagProcessors[i] = tagProcessorInss[i].(IBeanProcess)
		}
		return tagProcessors
	},
}

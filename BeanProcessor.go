package gioc

import "reflect"

type IBeanStartAssembly interface {
	// tag前置处理
	BeanStartAssembly(beanContainer IBeanContainer)
}

type IBeanBeforeTagProcess interface {
	// tag前置处理
	BeanBeforeTagProcess(tagProcessor ITagProcessor, beanContainer IBeanContainer)
}

type IBeanAfterTagProcess interface {
	// tag后缀处理
	BeanAfterTagProcess(tagProcessor ITagProcessor, beanContainer IBeanContainer)
}

type IBeanFinishAssembly interface {
	// tag前置处理
	BeanFinishAssembly(beanContainer IBeanContainer)
}

var IBeanStartAssemblyType = reflect.TypeOf((*IBeanStartAssembly)(nil)).Elem()
var IBeanBeforeTagProcessType = reflect.TypeOf((*IBeanBeforeTagProcess)(nil)).Elem()
var IBeanAfterTagProcessType = reflect.TypeOf((*IBeanAfterTagProcess)(nil)).Elem()
var IBeanFinishAssemblyType = reflect.TypeOf((*IBeanFinishAssembly)(nil)).Elem()

type beanProcessHelper struct{}

func (this *beanProcessHelper) GetBeanStartAssembly(beanContainer IBeanContainer) []IBeanStartAssembly {
	tagProcessorInss := beanContainer.GetBeanByType(IBeanStartAssemblyType)
	tagProcessors := make([]IBeanStartAssembly, len(tagProcessorInss), len(tagProcessorInss))
	for i := 0; i < len(tagProcessorInss); i++ {
		tagProcessors[i] = tagProcessorInss[i].(IBeanStartAssembly)
	}
	return tagProcessors
}

func (this *beanProcessHelper) GetBeanFinishAssembly(beanContainer IBeanContainer) []IBeanFinishAssembly {
	tagProcessorInss := beanContainer.GetBeanByType(IBeanFinishAssemblyType)
	tagProcessors := make([]IBeanFinishAssembly, len(tagProcessorInss), len(tagProcessorInss))
	for i := 0; i < len(tagProcessorInss); i++ {
		tagProcessors[i] = tagProcessorInss[i].(IBeanFinishAssembly)
	}
	return tagProcessors
}

var BeanProcessHelper = &beanProcessHelper{}

package gioc

import (
	"github.com/gosrv/goioc/util"
	"reflect"
)

// bean容器
type IBeanContainer interface {
	// 用这个加入的bean会直接加入容器，忽略装配，装配在builder中完成
	AddBean(beans ...interface{})
	AddNamedBean(name string, bean interface{})
	// 通过类型获取bean
	GetBeanByType(pt reflect.Type) []interface{}
	// 通过名字获取bean
	GetBeanByName(name string) interface{}
	// 获取所有的bean
	GetAllBeans() []interface{}
}

type EmptyBeanContainer struct {
}

type defaultBeanContainer struct {
	allBeans   []interface{}
	namedBeans map[string]interface{}
}

func NewBeanContainer() *defaultBeanContainer {
	return &defaultBeanContainer{
		namedBeans: make(map[string]interface{}),
	}
}

func (this *defaultBeanContainer) AddBean(beans ...interface{}) {
	for _, bean := range beans {
		if beanName, ok := bean.(IBeanName); ok {
			this.AddNamedBean(beanName.BeanName(), bean)
		} else {
			this.doAddBean(bean)
		}
	}
}

func (this *defaultBeanContainer) AddNamedBean(name string, bean interface{}) {
	if bean == nil {
		util.Panic("nil bean interface")
	}

	if beanName, ok := bean.(IBeanName); ok {
		if name != beanName.BeanName() {
			util.Panic("ambiguous bean name %v:%v", name, bean.(IBeanName).BeanName())
		}
	}

	if this.GetBeanByName(name) != nil {
		util.Panic("duplicate bean with name %v", name)
	}
	this.namedBeans[name] = bean
	this.doAddBean(bean)
}

func (this *defaultBeanContainer) doAddBean(bean interface{}) {
	for _, ebean := range this.allBeans {
		if bean == ebean {
			util.Panic("duplicate bean instance")
		}
	}
	this.allBeans = append(this.allBeans, bean)
}

func (this *defaultBeanContainer) GetBeanByName(name string) interface{} {
	bean, ok := this.namedBeans[name]
	if !ok {
		return nil
	}
	return bean
}

func (this *defaultBeanContainer) GetBeanByType(pt reflect.Type) []interface{} {
	tbeans := make([]interface{}, 0, 0)
	for _, bean := range this.allBeans {
		if reflect.TypeOf(bean).AssignableTo(pt) {
			tbeans = append(tbeans, bean)
		}
	}
	return tbeans
}

func (this *defaultBeanContainer) GetAllBeans() []interface{} {
	return this.allBeans
}

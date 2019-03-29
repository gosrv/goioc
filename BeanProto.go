package gioc

import "reflect"

type IBeanName interface {
	BeanName() string
}

var IBeanNameType = reflect.TypeOf((*IBeanName)(nil)).Elem()

type beanName struct {
	name string
}

func NewBeanName(name string) IBeanName {
	return &beanName{name: name}
}

func (this *beanName) BeanName() string {
	return this.name
}

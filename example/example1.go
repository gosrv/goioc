package main

import (
	"fmt"
	"github.com/gosrv/goioc"
	"reflect"
)
// beanA///////////////////////////////////////////////
type IBeanA interface {
	helloA()
}
type BeanA struct {
	msg string
	// beanB的实例由ioc容器注入
	beanB IBeanB	`bean`
}
func (this *BeanA)helloA() {
	fmt.Print(this.msg)
	this.beanB.helloB()
}
// beanB///////////////////////////////////////////////
type IBeanB interface {
	helloB()
}
type BeanB struct {
	msg string
}

func (this *BeanB)helloB()  {
	fmt.Println(this.msg)
}

func main() {
	// init
	builder := goioc.NewBeanContainerBuilder()
	builder.AddBean(goioc.NewBeanTagProcessor(builder.GetBeanContainer()))
	builder.AddBean(goioc.NewTagParser())

	// add beans and build
	builder.AddBean(&BeanA{msg:"hello "}, &BeanB{"world"})
	builder.Build()

	// find assembly bean
	container := builder.GetBeanContainer()
	beanA := container.GetBeanByType(reflect.TypeOf((*IBeanA)(nil)).Elem())[0].(IBeanA)
	beanA.helloA()
}
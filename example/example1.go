package main

import (
	"fmt"
	"github.com/gosrv/goioc"
	"reflect"
)

// beanA///////////////////////////////////////////////
type BeanC struct {
	beanA IBeanA `bean:"beana" bean.required:"true"`
}

type IBeanA interface {
	helloA()
}
type BeanA struct {
	msg string
	// beanB的实例由ioc容器注入
	beanB IBeanB `bean`
}

func (this *BeanA) helloA() {
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

func (this *BeanB) helloB() {
	fmt.Println(this.msg)
}

func main() {
	// 创建容器构建器
	builder := gioc.NewBeanContainerBuilder()
	// 加入tag处理器和tag解析器
	builder.AddBean(gioc.NewBeanTagProcessor(builder.GetBeanContainer()))
	builder.AddBean(gioc.NewTagParser())
	// 加入bean
	builder.AddBean(&BeanA{msg: "hello "}, &BeanB{"world"})
	// 构建容器
	builder.Build()
	// 获取bean
	container := builder.GetBeanContainer()
	beanA := container.GetBeanByType(reflect.TypeOf((*IBeanA)(nil)).Elem())[0].(IBeanA)
	beanA.helloA()
}

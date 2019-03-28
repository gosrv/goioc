package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosrv/goioc"
	"github.com/gosrv/goioc/util"
)

type BeanConfig struct {
	Name string
	Level int
	Age int
}

type Bean struct {
	ConfigA *BeanConfig `cfg:"cfg.a"`
	ConfigB *BeanConfig `cfg:"cfg.b"`
}

func main() {
	loader := goioc.NewConfigLoader()
	err := loader.Load("example/conf/config.json")
	util.VerifyNoError(err)

	// init
	builder := goioc.NewBeanContainerBuilder()
	builder.AddBean(goioc.NewBeanTagProcessor(builder.GetBeanContainer()))
	builder.AddBean(goioc.NewConfigValueTagProcessor(loader))
	builder.AddBean(goioc.NewTagParser())

	builder.AddNamedBean("bean", &Bean{})
	builder.Build()
	bean := builder.GetBeanContainer().GetBeanByName("bean").(*Bean)
	data, _ := json.Marshal(bean)
	fmt.Println(string(data))
}

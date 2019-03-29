package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosrv/goioc"
	"github.com/gosrv/goioc/util"
)

type BeanConfig struct {
	Name  string
	Level int
	Age   int
}

type Bean struct {
	ConfigA *BeanConfig `cfg:"cfg.a"`
	ConfigB *BeanConfig `cfg:"cfg.b"`
}

func main() {
	loader := gioc.NewConfigLoader()
	err := loader.Load("example/conf/config.json")
	util.VerifyNoError(err)

	// init
	builder := gioc.NewBeanContainerBuilder()
	builder.AddBean(gioc.NewBeanTagProcessor(builder.GetBeanContainer()))
	builder.AddBean(gioc.NewConfigValueTagProcessor(loader))
	builder.AddBean(gioc.NewTagParser())

	builder.AddNamedBean("bean", &Bean{})
	builder.Build()
	bean := builder.GetBeanContainer().GetBeanByName("bean").(*Bean)
	data, _ := json.Marshal(bean)
	fmt.Println(string(data))
}

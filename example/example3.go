package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosrv/goioc"
	"github.com/gosrv/goioc/util"
)

type BeanConfigBase struct {
	Name  string
	Level int
	Age   int
}

type BeanBase struct {
	gioc.IConfigBase
	ConfigA *BeanConfigBase `cfg:"cfg.a"`
	ConfigB *BeanConfigBase `cfg:"cfg.b"`
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

	builder.AddNamedBean("bean", &BeanBase{IConfigBase: gioc.NewConfigBase("cfg.base")})
	builder.Build()
	bean := builder.GetBeanContainer().GetBeanByName("bean").(*BeanBase)
	data, _ := json.Marshal(bean)
	fmt.Println(string(data))
}

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
	goioc.IConfigBase
	ConfigA *BeanConfigBase `cfg:"cfg.a"`
	ConfigB *BeanConfigBase `cfg:"cfg.b"`
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

	builder.AddNamedBean("bean", &BeanBase{IConfigBase: goioc.NewConfigBase("cfg.base")})
	builder.Build()
	bean := builder.GetBeanContainer().GetBeanByName("bean").(*BeanBase)
	data, _ := json.Marshal(bean)
	fmt.Println(string(data))
}

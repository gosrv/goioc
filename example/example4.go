package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosrv/goioc"
	"github.com/gosrv/goioc/util"
)

type BeanConfigCondition struct {
	Name  string
	Level int
	Age   int
}

type BeanCondition struct {
	gioc.IBeanCondition
	gioc.IConfigBase
	ConfigA *BeanConfigCondition `cfg.d:"cfg.a"`
	ConfigB *BeanConfigCondition `cfg:"cfg.b"`
}

func main() {
	loader := gioc.NewConfigLoader()
	err := loader.Load("example/conf/config.json")
	util.VerifyNoError(err)

	// init
	builder := gioc.NewBeanContainerBuilder()
	builder.AddBean(loader)
	builder.AddBean(gioc.NewBeanTagProcessor(builder.GetBeanContainer()))
	builder.AddBean(gioc.NewConfigValueTagProcessor(loader))
	builder.AddBean(gioc.NewTagParser())
	builder.AddBean(gioc.NewBeanBeanConditionInjector())

	builder.AddNamedBean("bean", &BeanCondition{
		IBeanCondition: gioc.NewConditionOnValue("cfg.base", true),
		IConfigBase:    gioc.NewConfigBase("cfg.base"),
	})
	builder.Build()
	bean := builder.GetBeanContainer().GetBeanByName("bean").(*BeanCondition)
	data, _ := json.Marshal(bean)
	fmt.Println(string(data))
}

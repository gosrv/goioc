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
	// 需要加载配置文件
	loader := gioc.NewConfigLoader()
	err := loader.Load("example/conf/config.json")
	util.VerifyNoError(err)
	// 创建容器构建器
	builder := gioc.NewBeanContainerBuilder()
	// 加入bean和config的tag处理器和tag解析器
	builder.AddBean(loader)
	builder.AddBean(gioc.NewBeanTagProcessor(builder.GetBeanContainer()))
	builder.AddBean(gioc.NewConfigValueTagProcessor(loader))
	builder.AddBean(gioc.NewTagParser())
	builder.AddBean(gioc.NewBeanBeanConditionInjector())
	// 加入bean，这个bean有条件：cfg.base存在时才生效
	builder.AddNamedBean("bean", &BeanCondition{
		IBeanCondition: gioc.NewConditionOnValue("cfg.base", true),
		IConfigBase:    gioc.NewConfigBase("cfg.base"),
	})
	// 构建容器
	builder.Build()
	// 获取bean
	bean := builder.GetBeanContainer().GetBeanByName("bean").(*BeanCondition)
	data, _ := json.Marshal(bean)
	fmt.Println(string(data))
}

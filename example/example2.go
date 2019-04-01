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
	// 注入配置文件的配置项"cfg.a"
	ConfigA *BeanConfig `cfg:"cfg.a"`
	// 注入配置文件的配置项"cfg.b"
	ConfigB *BeanConfig `cfg:"cfg.b"`
}

func main() {
	// 需要加载配置文件
	loader := gioc.NewConfigLoader()
	err := loader.Load("example/conf/config.json")
	util.VerifyNoError(err)
	// 创建容器构建器
	builder := gioc.NewBeanContainerBuilder()
	// 加入bean和config的tag处理器和tag解析器
	builder.AddBean(gioc.NewBeanTagProcessor(builder.GetBeanContainer()))
	builder.AddBean(gioc.NewConfigValueTagProcessor(loader))
	builder.AddBean(gioc.NewTagParser())
	// 加入bean
	builder.AddNamedBean("bean", &Bean{})
	// 构建容器
	builder.Build()
	// 获取bean
	bean := builder.GetBeanContainer().GetBeanByName("bean").(*Bean)
	data, _ := json.Marshal(bean)
	fmt.Println(string(data))
}

package goioc

import (
	"github.com/gosrv/goioc/util"
	"reflect"
	"strings"
)

/**
cfg相关tag处理
*/
const (
	ConfigTagProcessor = "config"
	// 从配置文件中读取的配置项
	ValueConfigTag = "cfg"
	// 默认配置项
	ValueDefaultTag = "cfg.default"
)

type configValueTagProcessor struct {
	conf                IConfigLoader
	stringValueInjector IStringValueInjector
}

func NewConfigValueTagProcessor(conf IConfigLoader) ITagProcessor {
	return &configValueTagProcessor{
		conf:                conf,
		stringValueInjector: StringPropertyInjects,
	}
}

func (this *configValueTagProcessor) TagProcessorName() string {
	return ConfigTagProcessor
}

func (this *configValueTagProcessor) GetPriority() int {
	return PrioritySystem
}

func (this *configValueTagProcessor) TagProcess(bean interface{}, field reflect.Value, tags map[string]string) {
	valDefault, defok := tags[ValueDefaultTag]
	valConfig, cfgok := tags[ValueConfigTag]
	if defok {
		// 注入默认值，目前只支持基本类型
		this.stringValueInjector.StringValueInjector(field, valDefault)
	}
	if cfgok {
		// 注入配置文件中的配置
		_, ok := this.conf.Config().Map()[valConfig]
		if !ok {
			subcfgs := strings.Split(valConfig, ".")
			err := this.conf.Config().Get(subcfgs...).Scan(field.Addr().Interface())
			if err != nil {
				util.Panic("config scan failed %v:%v", valConfig, err)
			}
		} else {
			err := this.conf.Config().Get(valConfig).Scan(field.Addr().Interface())
			if err != nil {
				util.Panic("config scan failed %v:%v", valConfig, err)
			}
		}
	}
}

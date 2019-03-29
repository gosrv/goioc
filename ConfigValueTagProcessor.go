package gioc

import (
	"github.com/gosrv/goioc/util"
	"github.com/micro/go-config/reader"
	"reflect"
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

func (this *configValueTagProcessor) Priority() int {
	return PrioritySystem
}

func getConfigValue(values reader.Values, path string) reader.Value {
	return values.Get(path)
}

func getConfigValueWithBase(values reader.Values, base, path string) reader.Value {
	if len(base) == 0 {
		return getConfigValue(values, path)
	}
	return values.Get(base, path)
}

func (this *configValueTagProcessor) TagProcess(bean interface{}, field reflect.Value, tags map[string]string) {
	valDefault, defok := tags[ValueDefaultTag]
	valConfig, cfgok := tags[ValueConfigTag]
	if defok {
		// 注入默认值，目前只支持基本类型
		this.stringValueInjector.StringValueInjector(field, valDefault)
	}
	if cfgok {
		var value reader.Value
		if reflect.TypeOf(bean).AssignableTo(IConfigBaseType) {
			cfgBaseName := bean.(IConfigBase).ConfigBase()
			if len(cfgBaseName) > 0 {
				value = this.conf.Config().Get(cfgBaseName, valConfig)
			} else {
				value = this.conf.Config().Get(valConfig)
			}
		} else {
			value = this.conf.Config().Get(valConfig)
		}

		err := value.Scan(field.Addr().Interface())
		if err != nil {
			util.Panic("config scan failed %v:%v", valConfig, err)
		}
	}
}

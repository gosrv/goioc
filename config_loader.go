package gioc

import (
	"github.com/gosrv/goioc/util"
	"github.com/micro/go-config"
	"github.com/micro/go-config/reader"
	"github.com/micro/go-config/source/file"
	"reflect"
)

type IConfigBase interface {
	ConfigBase() string
}

var IConfigBaseType = reflect.TypeOf((*IConfigBase)(nil)).Elem()

type configBase struct {
	base string
}

func NewConfigBase(base string) *configBase {
	return &configBase{base: base}
}

func (this configBase) ConfigBase() string {
	return this.base
}

type IConfigLoader interface {
	Config() reader.Values
	Load(cfgFileName string) error
	AutoLoad(func())
}

func NewConfigLoader() *configAutoLoader {
	return &configAutoLoader{
		conf: config.NewConfig(),
	}
}

type configAutoLoader struct {
	conf config.Config
}

func (this *configAutoLoader) AutoLoad(loader func()) {
	watch, err := this.RawConf().Watch()
	util.VerifyNoError(err)
	watch.Next()

	loader()
}

func (this *configAutoLoader) RawConf() config.Config {
	return this.conf
}

func (this *configAutoLoader) Config() reader.Values {
	return this.conf
}

func (this *configAutoLoader) Load(cfgFileName string) error {
	return this.conf.Load(file.NewSource(
		file.WithPath(cfgFileName),
	))
}

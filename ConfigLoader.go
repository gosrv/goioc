package goioc

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/reader"
	"github.com/micro/go-config/source/file"
)

type IConfigLoader interface {
	Config() reader.Values
	Load(cfgFileName string) error
}

type IConfigLoaderAware interface {
	SetConfigLoader(configLoader IConfigLoader)
}

func NewConfigLoader() *DefaultConfigAutoLoader {
	return &DefaultConfigAutoLoader{
		conf: config.NewConfig(),
	}
}

type DefaultConfigAutoLoader struct {
	conf config.Config
}

func (this *DefaultConfigAutoLoader) RawConf() config.Config {
	return this.conf
}

func (this *DefaultConfigAutoLoader) Config() reader.Values {
	return this.conf
}

func (this *DefaultConfigAutoLoader) Load(cfgFileName string) error {
	return this.conf.Load(file.NewSource(
		file.WithPath(cfgFileName),
	))
}

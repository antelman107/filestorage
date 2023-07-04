package providers

import (
	"github.com/antelman107/filestorage/internal/config"
	"github.com/antelman107/filestorage/pkg/domain"
)

type defaultConfigLoader struct {
	additionalPaths []string
}

func NewDefaultConfigLoader(additionalPaths ...string) domain.ConfigLoader {
	return &defaultConfigLoader{
		additionalPaths: additionalPaths,
	}
}

func (p *defaultConfigLoader) Load(name string, in interface{}) error {
	return config.Load(name, in, p.additionalPaths...)
}

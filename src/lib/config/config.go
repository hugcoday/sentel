package config

import (
	"fmt"
	"os"

	mc "github.com/koding/multiconfig"
)

type Config interface {
  GetKey(name string)string
}

type ConfigLoader struct {
	mc.DefaultLoader
}

func NewWithPath(path string) *ConfigLoader {
	loader := &ConfigLoader{}
	loader.DefaultLoader = *mc.NewWithPath(path)
	return loader
}

func MustLoadWithPath(path string, conf interface{}) {
	d := NewWithPath(path)
	d.MustLoad(conf)
}

func (c *ConfigLoader) MustLoad(conf interface{}) {
	if err := c.Load(conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func (c *ConfigLoader) MustValidate(conf interface{}) {
	c.MustValidate(conf)
}

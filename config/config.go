package config

import (
	"fmt"

	"os"
	"sync"

	"github.com/go-kit/log"
	"gopkg.in/yaml.v2"
)

func LoadFile(filename string) (*Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}
	cfg := &Config{}
	err = yaml.UnmarshalStrict(content, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

type Config struct {
	Dingtalk map[string]Dingtalk
	Wechat   map[string]Wechat
}

type SafeConfig struct {
	sync.RWMutex
	C *Config
}

func (sc *SafeConfig) ReloadConfig(confFile string, logger log.Logger) error {
	c, err := LoadFile(confFile)
	if err != nil {
		return err
	}

	sc.Lock()
	sc.C = c
	sc.Unlock()

	return nil
}

type Dingtalk struct {
	URL    string `yaml:"url,omitempty"`
	Secret string `yaml:"secret,omitempty"`
}
type Wechat struct {
	URL string
}

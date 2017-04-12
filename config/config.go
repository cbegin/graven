package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"github.com/bgentry/speakeasy"
	"fmt"
)

const DefaultConfigFileName = ".graven.yaml"

type Config struct {
	configFileName string
	data map[string]map[string]string
}

func NewConfig() Config {
	config := Config{}
	config.data = map[string]map[string]string{}
	config.configFileName = DefaultConfigFileName
	return config
}

func (c Config) Set(group, name string, value string) {
	if g, ok := c.data[group]; ok {
		g[name] = value
	} else {
		c.data[group] = map[string]string{}
		c.data[group][name] = value
	}
}

func (c Config) Get(group, name string) string {
	if g, ok := c.data[group]; ok {
		return g[name]
	}
	return ""
}

func (c Config) SetSecret(group, name, prompt string) error {
	password, err := speakeasy.Ask(prompt)
	if err != nil {
		return fmt.Errorf("Error reading secret from terminal: %v", err)
	}
	c.Set(group, name, password)
	return nil
}

func (c Config) Read() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	file, err := os.Open(path.Join(usr.HomeDir, c.configFileName))
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, c.data)
	if err != nil {
		return err
	}
	return nil
}

func (c Config) Write() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	bytes, err := yaml.Marshal(c.data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(usr.HomeDir, c.configFileName), bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

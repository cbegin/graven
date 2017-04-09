package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

const DefaultConfigFileName = ".graven.yaml"

type Config struct {
	configFileName string
	data map[string]interface{}
}

func NewConfig() Config {
	config := Config{}
	config.data = map[string]interface{}{}
	config.configFileName = DefaultConfigFileName
	return config
}

func (c Config) Set(name string, value interface{}) {
	c.data[name] = value
}

func (c Config) Get(name string) interface{} {
	return c.data[name]
}

func (c Config) GetString(name string) string {
	if v, ok := c.data[name]; ok {
		return v.(string)
	}
	return ""
}

func (c Config) GetMap(name string) map[string]interface{} {
	if v, ok := c.data[name]; ok {
		return v.(map[string]interface{})
	}
	return map[string]interface{}{}
}

func (c Config) GetMaps(name string) []map[string]interface{} {
	if v, ok := c.data[name]; ok {
		return v.([]map[string]interface{})
	}
	return []map[string]interface{}{}
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

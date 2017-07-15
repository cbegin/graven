package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"bufio"

	"github.com/bgentry/speakeasy"
	"gopkg.in/yaml.v2"
)

const DefaultConfigFileName = ".graven.yaml"

type Config struct {
	configFileName string
	data           map[string]map[string]string
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

func (c Config) SetPlainText(group, name, prompt string) error {
	fmt.Print(prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	if scanner.Err() != nil {
		return fmt.Errorf("Error reading input from terminal: %v", scanner.Err())
	}
	c.Set(group, name, scanner.Text())
	return nil
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
	file, err := os.Open(filepath.Join(usr.HomeDir, c.configFileName))
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
	err = ioutil.WriteFile(filepath.Join(usr.HomeDir, c.configFileName), bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

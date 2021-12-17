// This is a minimalist config package. It uses a flat structure
// and simple obfuscation to avoid over-the-shoulder or casual viewing
// of passwords on the file system. Future implemenations may improve upon
// the data structure and allow for user provided passwords that will
// better protect the locally stored password.
// Or we can try to find a config library that supports:
//   - getting / setting values in a structured way (like Viper)
//   - securely prompting for values (I suppose this could be externalized)
//   - encrypting / decrypting / obfuscating stored values
//   - loading and saving configuration (not just reading)
package config

import (
	"bufio"
	"fmt"
	"github.com/cbegin/graven/internal/util"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

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

func (c Config) GetSecret(group, name string) (string, error) {
	if cipherText := c.Get(group, name); cipherText != "" {
		plainText, err := util.Uncloak(cipherText)
		if err != nil {
			return "", err
		}
		return string(plainText), nil
	}
	return "", nil
}

func (c Config) PromptPlainText(group, name, prompt string) error {
	fmt.Print(prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	if scanner.Err() != nil {
		return fmt.Errorf("Error reading input from terminal: %v", scanner.Err())
	}
	c.Set(group, name, scanner.Text())
	return nil
}

func (c Config) PromptSecret(group, name, prompt string) error {
	plainText, err := speakeasy.Ask(prompt)
	if err != nil {
		return fmt.Errorf("Error reading secret from terminal: %v", err)
	}

	cipherText, err := util.Cloak(plainText)
	if err != nil {
		return fmt.Errorf("Error encrypting secret: %v", err)
	}

	c.Set(group, name, cipherText)
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

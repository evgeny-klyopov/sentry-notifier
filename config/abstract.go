package config

import (
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"os"
)

type builder interface {
	getContent() ([]byte, error)
	build(content []byte) error
	validate() error
	getConfig() *Config
}

type abstract struct {
	path   string
	config *Config
}

func (a *abstract) validate() error {
	var validate = validator.New()
	return validate.Struct(a.config)
}

func (a *abstract) getConfig() *Config {
	return a.config
}

func (a *abstract) getContent() ([]byte, error) {
	file, err := os.Open(a.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

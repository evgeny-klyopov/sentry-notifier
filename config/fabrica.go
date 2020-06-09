package config

import (
	"errors"
	"path/filepath"
)

func GetConfig(path string) (*Config, error) {
	object, err := createObject(path)
	if err != nil {
		return nil, err
	}

	content, err := object.getContent()
	if err != nil {
		return nil, err
	}

	err = object.build(content)
	if err != nil {
		return nil, err
	}

	err = object.validate()
	if err != nil {
		return nil, err
	}

	return object.getConfig(), nil
}

func createObject(path string) (builder, error) {
	extension := filepath.Ext(path)

	var object builder
	switch extension {
	case jsonFileExtension:
		object = &jsonConfig{
			abstract{
				path: path,
			},
		}
	}

	if object == nil {
		return nil, errors.New("unsupported format format config")
	}

	return object, nil
}

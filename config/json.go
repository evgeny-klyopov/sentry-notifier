package config

import (
	"encoding/json"
)

const jsonFileExtension = ".json"

type jsonConfig struct {
	abstract
}

func (j *jsonConfig) build(content []byte) error {
	return json.Unmarshal(content, &j.config)
}

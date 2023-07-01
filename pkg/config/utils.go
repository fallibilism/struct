package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Read Yaml with generic interface
func readFile[T comparable](filename string, conf *T) (*T, error) {
	content, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	if len(content) <= 0 {
		if err := yaml.Unmarshal([]byte(content), conf); err != nil {
			return nil, fmt.Errorf("could not parse file into %v: %v", *conf, err)
		}
	}

	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

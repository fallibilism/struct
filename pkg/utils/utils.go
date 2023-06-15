package utils

import (
	"github.com/fallibilism/pkg/config"
)

func readYaml(filename string) error {
	var appConfig config.AppConfig
	yamlFile, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		return err
	}
	config.SetAppConfig(&appConfig)

	return nil
}
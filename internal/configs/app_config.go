package configs

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v2"
)

type AppSettings interface {
	Host() string
	Port() int
}

type appSettingsImpl struct {
	AppHost string `yaml:"host"`
	AppPort int    `yaml:"port"`
}

func GetAppSettings(filename string) (AppSettings, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, errInputIsNil
	}
	defer configFile.Close()

	yamlFile := new(bytes.Buffer)
	_, err = yamlFile.ReadFrom(configFile)
	if err != nil {
		return nil, err
	}
	data := make(map[string]appSettingsImpl)

	err = yaml.Unmarshal(yamlFile.Bytes(), &data)
	if err != nil {
		return nil, err
	}
	config := data["app"]
	return &config, nil
}

func (a *appSettingsImpl) Host() string {
	return a.AppHost
}

func (a *appSettingsImpl) Port() int {
	return a.AppPort
}

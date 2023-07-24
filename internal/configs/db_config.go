package configs

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v2"
)

type DBConnectionConfig interface {
	Host() string
	Port() int
	DatabaseName() string
}

type dbConnectionImpl struct {
	HostDB string `yaml:"host"`
	PortDB int    `yaml:"port"`
	NameDB string `yaml:"name"`
}

func GetDBConnectionConfig(filename string) (DBConnectionConfig, error) {
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
	data := make(map[string]dbConnectionImpl)

	err = yaml.Unmarshal(yamlFile.Bytes(), &data)
	if err != nil {
		return nil, err
	}
	config := data["database"]
	return &config, nil
}

func (c *dbConnectionImpl) Host() string {
	return c.HostDB
}

func (c *dbConnectionImpl) Port() int {
	return c.PortDB
}

func (c *dbConnectionImpl) DatabaseName() string {
	return c.NameDB
}

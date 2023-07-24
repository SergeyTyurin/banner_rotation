package configs

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v2"
)

type MessageBrokerConfig interface {
	Host() string
	Port() int
	URL() string
}

type messageBrokerImpl struct {
	HostBR string `yaml:"host"`
	PortBR int    `yaml:"port"`
	UrlBR  string `yaml:"url"`
}

func GetMessageBrokerConfig(filename string) (MessageBrokerConfig, error) {
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
	data := make(map[string]messageBrokerImpl)

	err = yaml.Unmarshal(yamlFile.Bytes(), &data)
	if err != nil {
		return nil, err
	}
	config := data["message_broker"]
	return &config, nil
}

func (c *messageBrokerImpl) Host() string {
	return c.HostBR
}

func (c *messageBrokerImpl) Port() int {
	return c.PortBR
}

func (c *messageBrokerImpl) URL() string {
	return c.UrlBR
}

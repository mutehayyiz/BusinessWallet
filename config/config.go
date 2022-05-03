package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type Config struct {
	Host      string         `json:"host"`
	Port      int            `json:"port"`
	APISecret string         `json:"api_secret"`
	Storage   StorageOptions `json:"storage"`
}

type StorageOptions struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DB       string `json:"db"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var Global Config

func (c *Config) Load(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.WithError(err).Error("Couldn't read config file")
		return err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		logrus.WithError(err).Error("Couldn't unmarshal configuration")
		return err
	}

	return err
}

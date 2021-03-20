package config

import (
	"github.com/gophers0/users/internal/repository/postgres"
	gaarx "github.com/zergu1ar/Gaarx"
)

const (
	TypeOfConfigLocal      = "local"
	TypeOfConfigTest       = "test"
	TypeOfConfigProduction = "prod"
)

type Config struct {
	Name   string `json:"name"`
	System struct {
		DB  *DB `json:"db"`
		Log struct {
			Filename string `json:"filename"`
			Level    int    `json:"level"`
		} `json:"log"`
	} `json:"system"`
	Api struct {
		Port string `json:"port"`
	} `json:"api"`
}
type DB struct {
	Dialect  string `json:"postgres"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func (c *Config) GetConnString() string {
	return postgres.GetConnString(
		c.System.DB.User,
		c.System.DB.Password,
		c.System.DB.Host,
		c.System.DB.Port,
		c.System.DB.Database,
	)
}

func (c *Config) GetLogWay() gaarx.LogWay {
	if c.System.Log.Filename != "" {
		return gaarx.FileLog
	}
	return ""
}

func (c *Config) GetLogDestination() string {
	return c.System.Log.Filename
}

func (c *Config) GetLogApplicationName() string {
	return "Users"
}

func Init(cfgType string) (string, error) {
	configFile := ""
	switch cfgType {
	case TypeOfConfigLocal:
		configFile = "config_local.json"
	case TypeOfConfigTest:
		configFile = "config_test.json"
	case TypeOfConfigProduction:
		configFile = "config.json"
	}

	if configFile == "" {
		panic("invalid config Type: " + cfgType)
	}

	return configFile, nil
}

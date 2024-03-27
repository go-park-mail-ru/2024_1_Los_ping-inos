package config

import (
	"github.com/spf13/viper"
	"time"
)

const PersonTableName = "person"
const InterestTableName = "interest"

var Cfg Config

type Config struct {
	ApiPath    string
	Server     ServerConfig     `json:"server"`
	Database   DatabaseConfig   `json:"database"`
	FilesPaths FilesPathsConfig `json:"filesPaths"`
}

type ServerConfig struct {
	Host        string        `json:"host"`
	Port        string        `json:"port"`
	SwaggerPort string        `json:"swaggerPort"`
	Timeout     time.Duration `json:"timeout"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"username"`
	Password string `json:"password"`
}

type FilesPathsConfig struct {
}

func LoadConfig(path string) (*Config, error) {
	var err error
	var config Config

	viper.SetConfigFile(path)

	err = viper.BindEnv("server.host", "SERVER_HOST")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("server.port", "SERVER_PORT")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("server.timeout", "SERVER_TIMEOUT")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.host", "DB_HOST")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.port", "DB_PORT")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.user", "DB_USER")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.password", "DB_PASSWORD")
	if err != nil {
		return nil, err
	}
	err = viper.BindEnv("database.dbname", "DB_NAME")
	if err != nil {
		return nil, err
	}
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	Cfg = config

	return &config, nil
}

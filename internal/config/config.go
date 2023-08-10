package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	HTTP ServerConfig `yaml:"http"`
	DB   DBConfig     `yaml:"db"`
	Test TestConfig   `json:"test"`
}

type ServerConfig struct {
	Port            string        `yaml:"port"`
	Timeout         time.Duration `yaml:"timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
}

type Collections struct {
	Task string `yaml:"task"`
}

type DBConfig struct {
	Host        string      `yaml:"host"`
	Port        string      `yaml:"port"`
	DBName      string      `yaml:"db_name"`
	Username    string      `yaml:"username"`
	Password    string      `yaml:"password"`
	Collections Collections `yaml:"collections"`
}

type TestConfig struct {
	DB DBConfig
}

func InitConfig(path string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)

	return cfg, nil
}

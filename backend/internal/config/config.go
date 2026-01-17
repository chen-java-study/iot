package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Wechat   WechatConfig   `yaml:"wechat"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type JWTConfig struct {
	SecretKey   string `yaml:"secret_key"`
	ExpireHours int    `yaml:"expire_hours"`
}

type WechatConfig struct {
	AppID          string `yaml:"app_id"`
	MchID          string `yaml:"mch_id"`
	APIV3Key       string `yaml:"api_v3_key"`
	SerialNo       string `yaml:"serial_no"`
	PrivateKeyPath string `yaml:"private_key_path"`
	NotifyURL      string `yaml:"notify_url"`
}

func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

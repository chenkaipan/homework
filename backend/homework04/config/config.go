package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port       string `json:"port"`
	DBType     string `json:"db_type"`
	SqlitePath string `json:"sqlite_path"`
	MySQLDSN   string `json:"mysql_dsn"`
	JWTSecret  string `json:"jwt_secret"`
}

var AppConfig *Config

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	AppConfig = &cfg
	return nil
}

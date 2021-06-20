package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ConnString string `json:"connstring"`
}

func InitConfig() (Config, error) {
	file, err := os.ReadFile("configs/cfg.json")
	if err != nil {
		return Config{}, fmt.Errorf("Can't open file: %w", err)
	}
	var cfg Config
	if err = json.Unmarshal(file, &cfg); err != nil {
		return Config{}, fmt.Errorf("Can't unmarshall json file: %w", err)
	}

	return cfg, nil
}

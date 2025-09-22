package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Addr string
	Port int
}

func New() *Config {
	addr := getEnv("APP_ADDR", "localhost")
	port := getEnvAsInt("APP_PORT", 3000)

	return &Config{
		Addr: addr,
		Port: port,
	}
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	if valStr, ok := os.LookupEnv(name); ok {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}

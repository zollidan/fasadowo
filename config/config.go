package config

type Config struct {
	Addr string
	Port int
}

func (c *Config) New() {
	return nil
}
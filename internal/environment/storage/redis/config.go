package redis

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewConfig(v *viper.Viper) Config {
	return Config{v: v}
}

type Config struct {
	v *viper.Viper
}

func (c *Config) GetHost() string {
	return c.v.GetString("redis.host")
}

func (c *Config) GetPort() int {
	return c.v.GetInt("redis.port")
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.GetHost(), c.GetPort())
}

func (c *Config) GetUsername() string {
	return c.v.GetString("redis.username")
}

func (c *Config) GetPassword() string {
	return c.v.GetString("redis.password")
}

func (c *Config) GetDatabase() int {
	return c.v.GetInt("redis.database")
}

package internalgrpc

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
	return c.v.GetString("grpc.host")
}

func (c *Config) GetPort() int {
	return c.v.GetInt("grpc.port")
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.GetHost(), c.GetPort())
}

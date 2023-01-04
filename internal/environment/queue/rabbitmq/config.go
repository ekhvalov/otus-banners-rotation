package rabbitmq

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func NewConfig(v *viper.Viper) *Config {
	return &Config{v: v}
}

type Config struct {
	v *viper.Viper
}

func (c *Config) GetDSN() string {
	dsnBuilder := strings.Builder{}
	dsnBuilder.WriteString("amqp://")
	if username := c.GetUsername(); username != "" {
		dsnBuilder.WriteString(fmt.Sprintf("%s:%s@", username, c.GetPassword()))
	}
	dsnBuilder.WriteString(fmt.Sprintf("%s:%d/", c.GetHost(), c.GetPort()))
	return dsnBuilder.String()
}

func (c *Config) GetHost() string {
	return c.v.GetString("rabbitmq.host")
}

func (c *Config) GetPort() int {
	return c.v.GetInt("rabbitmq.port")
}

func (c *Config) GetUsername() string {
	return c.v.GetString("rabbitmq.username")
}

func (c *Config) GetPassword() string {
	return c.v.GetString("rabbitmq.password")
}

func (c *Config) GetQueueName() string {
	return c.v.GetString("rabbitmq.queue_name")
}

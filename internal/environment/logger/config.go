package logger

import (
	"strings"

	"github.com/spf13/viper"
)

func NewConfig(v *viper.Viper) Config {
	return Config{v: v}
}

type Config struct {
	v *viper.Viper
}

func (c *Config) GetLevel() LogLevel {
	level := c.v.GetString("logger.level")
	switch strings.ToLower(level) {
	case string(LevelDebug):
		return LevelDebug
	case "warn":
		fallthrough
	case string(LevelWarning):
		return LevelWarning
	case string(LevelError):
		return LevelError
	default:
		return LevelInfo
	}
}

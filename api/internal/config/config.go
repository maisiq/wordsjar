package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var config *Config

type Database struct {
	DSN string `mapstructure:"dsn"`
}

type App struct {
	Debug bool   `mapstructure:"debug"`
	Host  string `mapstrcture:"host"`
	Port  string `mapstrcture:"port"`
}

type Cookies struct {
	BaseDomain string `mapstructure:"base_domain"`
	MaxAge     int    `mapstructure:"max_age"`
	Secure     bool   `mapstructure:"secure"`
}

type HTTP struct {
	Cookies               Cookies  `mapstructure:"cookies"`
	AllowedHosts          []string `mapstructure:"allowed_hosts"`
	ServerShutdownTimeout int      `mapstructure:"server_shutdown_timeout"`
}

type Config struct {
	Database Database `mapstructure:"database"`
	App      App      `mapstructure:"app"`
	HTTP     HTTP     `mapstructure:"http"`
}

func Load(path string) Config {
	v := viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
	)
	v.SetConfigFile(path)

	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	cfg := new(Config)
	if err := v.Unmarshal(cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %s", err)
	}

	config = cfg
	return *config
}

func GetHTTP() HTTP {
	return config.HTTP
}

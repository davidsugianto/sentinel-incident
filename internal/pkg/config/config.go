package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Auth     Auth     `yaml:"auth"`
	CORS     CORS     `yaml:"cors"`
}

type Server struct {
	Port int `yaml:"port"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

type Auth struct {
	JWTSecret string `yaml:"jwt_secret"`
}

type CORS struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

var (
	cfg  *Config
	once sync.Once
	mu   sync.RWMutex
)

func LoadConfig(path string) (*Config, error) {
	var err error
	once.Do(func() {
		v := viper.New()
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(path)
		v.AutomaticEnv()

		if err = v.ReadInConfig(); err != nil {
			err = fmt.Errorf("error reading config file: %w", err)
			return
		}

		var c Config
		if err = v.Unmarshal(&c); err != nil {
			err = fmt.Errorf("unable to decode config: %w", err)
			return
		}

		cfg = &c

		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("Config file changed: %s", e.Name)

			var newCfg Config
			if err := v.Unmarshal(&newCfg); err != nil {
				log.Printf("Failed to reload config: %v", err)
				return
			}

			mu.Lock()
			cfg = &newCfg
			mu.Unlock()
			log.Println("Config reloaded successfully")
		})
	})

	return cfg, err
}

func GetConfig() *Config {
	mu.RLock()
	defer mu.RUnlock()
	return cfg
}

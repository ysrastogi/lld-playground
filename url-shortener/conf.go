package urlshortener

import "github.com/spf13/viper"

type Config struct {
	Redis    Redis    `mapstructure:"redis"`
	Database Database `mapstructure:"database"`
}

type Redis struct {
	RedisURL string `mapstructure:"redis_url"`
}

type Database struct {
	DBPath string `mapstructure:"db_path"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

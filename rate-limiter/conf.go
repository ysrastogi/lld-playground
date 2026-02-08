package ratelimiter

import "github.com/spf13/viper"

type Config struct {
	Redis 
}

type Redis struct {
	redis_url string `mapstructure:"redis_url"`
}

func Load()(*Config, error){
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
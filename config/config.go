package config

import (
    "log"

    "github.com/spf13/viper"
)

type Config struct {
    Server struct {
        Port         string `mapstructure:"port"`
        ReadTimeout  int    `mapstructure:"read_timeout"`
        WriteTimeout int    `mapstructure:"write_timeout"`
    } `mapstructure:"server"`
    BestChange struct {
        APIKey string `mapstructure:"api_key"`
    } `mapstructure:"bestchange"`
}

var Cfg *Config

func Load() error {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")

    viper.AutomaticEnv() 
    viper.SetDefault("server.port", "8080")
    viper.SetDefault("server.read_timeout", 10)
    viper.SetDefault("server.write_timeout", 10)

    if err := viper.ReadInConfig(); err != nil {
        log.Printf("Config file not found, using defaults/env: %v", err)
    }

    Cfg = &Config{}
    if err := viper.Unmarshal(Cfg); err != nil {
        return err
    }
    return nil
}

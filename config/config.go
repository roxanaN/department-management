package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	ServerPort    string `mapstructure:"SERVER_PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	EmailFrom     string `mapstructure:"EMAIL_FROM"`
	EmailPassword string `mapstructure:"EMAIL_PASSWORD"`
	SMTPHost      string `mapstructure:"SMTP_HOST"`
	SMTPPort      string `mapstructure:"SMTP_PORT"`
}

// Get environment variables from .env file
func LoadConfig() (Config, error) {
	var config Config
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}

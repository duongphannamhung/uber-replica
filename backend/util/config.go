package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver                  string        `mapstructure:"DB_DRIVER"`
	DBSource                  string        `mapstructure:"DB_SOURCE"`
	ServerAddress             string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey         string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDurationDriver time.Duration `mapstructure:"ACCESS_TOKEN_DURATION_DRIVER"`
	AccessTokenDurationUser   time.Duration `mapstructure:"ACCESS_TOKEN_DURATION_USER"`

	TwiolioAccountSid string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwiolioAuthToken  string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwiolioFrom       string `mapstructure:"TWILIO_FROM"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

package configs

import (
	"github.com/spf13/viper"
)

type Conf struct {
	RateLimiterQtyToken  string `mapstructure:"RATE_LIMITER_QTY_TOKEN"`
	RateLimiterTimeToken string `mapstructure:"RATE_LIMITER_TIME_TOKEN"`
	RateLimiterQtyIp     string `mapstructure:"RATE_LIMITER_QTY_IP"`
	RateLimiterTimeIp    string `mapstructure:"RATE_LIMITER_TIME_IP"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}

package config

type Config struct {
	DBUri                string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri             string `mapstructure:"REDIS_URL"`
	Port                 string `mapstructure:"PORT"`
	Origin               string `mapstructure:"CLIENT_ORIGIN"`
	AccessTokenPublicKey string `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
}

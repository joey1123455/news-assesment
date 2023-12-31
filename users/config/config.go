package config

import "time"

type Config struct {
	DBUri                  string        `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri               string        `mapstructure:"REDIS_URL"`
	Port                   string        `mapstructure:"PORT"`
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
	Origin                 string        `mapstructure:"CLIENT_ORIGIN"`
	EmailFrom              string        `mapstructure:"EMAIL_FROM"`
	SMTPHost               string        `mapstructure:"SMTP_HOST"`
	SMTPPass               string        `mapstructure:"SMTP_PASS"`
	SMTPPort               int           `mapstructure:"SMTP_PORT"`
	SMTPUser               string        `mapstructure:"SMTP_USER"`
}

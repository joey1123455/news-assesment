package config

type Config struct {
	ApiKey    string `mapstructure:"NEWS_DATA_API_KEY"`
	DBUri     string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri  string `mapstructure:"REDIS_URL"`
	Port      string `mapstructure:"PORT"`
	Origin    string `mapstructure:"CLIENT_ORIGIN"`
	EmailFrom string `mapstructure:"EMAIL_FROM"`
	SMTPHost  string `mapstructure:"SMTP_HOST"`
	SMTPPass  string `mapstructure:"SMTP_PASS"`
	SMTPPort  int    `mapstructure:"SMTP_PORT"`
	SMTPUser  string `mapstructure:"SMTP_USER"`
}

package config

import (
	"errors"
	"time"

	"github.com/joho/godotenv"
	"github.com/shanisharrma/gopher-social/internal/env"
	"github.com/shanisharrma/gopher-social/internal/ratelimiter"
)

type Config struct {
	Env         string
	Addr        string
	ApiUrl      string
	FrontendURL string
	Db          dbConfig
	Mail        mailConfig
	Auth        authConfig
	RedisCfg    redisConfig
	Ratelimiter ratelimiter.Config
}

type redisConfig struct {
	Addr    string
	Pw      string
	DB      int
	Enabled bool
}

type authConfig struct {
	Basic basicConfig
	Token tokenConfig
}

type tokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
}

type basicConfig struct {
	User string
	Pass string
}

type mailConfig struct {
	Sendgrid  SendgridConfig
	Exp       time.Duration
	FromEmail string
	Mailtrap  MailtrapConfig
}

type dbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type SendgridConfig struct {
	ApiKey string
}

type MailtrapConfig struct {
	ApiKey   string
	Username string
}

func MustLoad() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, errors.New("error loading .env file")
	}

	return Config{
		Env:         env.GetString("ENV", "development"),
		Addr:        env.GetString("ADDR", ":8080"),
		ApiUrl:      env.GetString("EXTERNAL_URL", "localhost:8000"),
		FrontendURL: env.GetString("FRONTEND_URL", "localhost:5173"),
		Db: dbConfig{
			Addr: env.GetString(
				"DB_ADDR",
				"postgres://admin:adminpassword@localhost/social?sslmode=disable",
			),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15min"),
		},
		RedisCfg: redisConfig{
			Addr:    env.GetString("REDDIS_ADDR", "localhost:6379"),
			Pw:      env.GetString("REDIS_PW", ""),
			DB:      env.GetInt("REDIS_DB", 0),
			Enabled: env.GetBool("REDIS_ENABLED", false),
		},
		Mail: mailConfig{
			Exp:       time.Hour * 24 * 3,
			FromEmail: env.GetString("FROM_EMAIL", ""),
			Sendgrid: SendgridConfig{
				ApiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			Mailtrap: MailtrapConfig{
				ApiKey:   env.GetString("MAILTRAP_API_KEY", ""),
				Username: env.GetString("MAILTRAP_USERNAME", ""),
			},
		},
		Auth: authConfig{
			Basic: basicConfig{
				User: env.GetString("BASIC_AUTH_USER", "admin"),
				Pass: env.GetString("BASIC_AUTH_PASS", "pass"),
			},
			Token: tokenConfig{
				Secret: env.GetString("AUTH_TOKEN_SECRET", "defaultsecret12345678"),
				Exp:    time.Hour * 24 * 3, // 3 days
				Iss:    "gophersocial",
			},
		},
		Ratelimiter: ratelimiter.Config{
			RequestPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame:           time.Second * 5,
			Enabled:             env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}, nil
}

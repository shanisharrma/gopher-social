package main

import (
	"expvar"
	"log"
	"runtime"

	"github.com/go-redis/redis/v8"
	"github.com/shanisharrma/gopher-social/cmd/api/config"
	"github.com/shanisharrma/gopher-social/internal/auth"
	"github.com/shanisharrma/gopher-social/internal/db"
	"github.com/shanisharrma/gopher-social/internal/mailer"
	"github.com/shanisharrma/gopher-social/internal/ratelimiter"
	"github.com/shanisharrma/gopher-social/internal/store"
	"github.com/shanisharrma/gopher-social/internal/store/cache"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			GopherSocial
//	@description	An API for GopherSocial, social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v2
//
// @SecurityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// configuring Database connection
	db, err := db.New(cfg.Db.Addr, cfg.Db.MaxOpenConns, cfg.Db.MaxIdleConns, cfg.Db.MaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("Database connection pool established")

	// Cache
	var rdb *redis.Client
	if cfg.RedisCfg.Enabled {
		rdb = cache.NewRedisClient(cfg.RedisCfg.Addr, cfg.RedisCfg.Pw, cfg.RedisCfg.DB)
		logger.Info("redis cache connection established")
	}

	// Rate Limiter
	ratelimiter := ratelimiter.NewFixedWindowRateLimiter(cfg.Ratelimiter.RequestPerTimeFrame, cfg.Ratelimiter.TimeFrame)

	// configuring database store (relational and cache)
	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	// Configuring the mailer
	// mailer := mailer.NewSendGrid(cfg.Mail.FromEmail, cfg.Mail.Sendgrid.ApiKey)
	mailtrap, err := mailer.NewMailtrap(cfg.Mail.Mailtrap.ApiKey, cfg.Mail.Mailtrap.Username, cfg.Mail.FromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	// JWT Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.Auth.Token.Secret, cfg.Auth.Token.Iss, cfg.Auth.Token.Iss)

	// injecting dependencies to application
	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailtrap,
		authenticator: jwtAuthenticator,
		ratelimiter:   ratelimiter,
	}

	// metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	logger.Fatal(app.run(mux))
}

package main

import (
	"expvar"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/theluminousartemis/socialnews/internal/auth"
	"github.com/theluminousartemis/socialnews/internal/db"
	"github.com/theluminousartemis/socialnews/internal/env"
	"github.com/theluminousartemis/socialnews/internal/mailer"
	"github.com/theluminousartemis/socialnews/internal/ratelimiter"
	"github.com/theluminousartemis/socialnews/internal/store"
	"github.com/theluminousartemis/socialnews/internal/store/cache"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Wise.ly
//	@description	API for Wise.ly, a community driven Q&A platform.

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("API_URL", "localhost:3000"),

		db: &dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/wisely?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "production"),
		mail: mailConfig{
			exp:       time.Hour * 24,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			username:  env.GetString("MAILTRAP_USERNAME", ""),
			password:  env.GetString("MAILTRAP_PASSWORD", ""),
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: jwtConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3,
				iss:    "wise.ly",
			},
		},
		redisCfg: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
			enabled:  env.GetBool("REDIS_ENABLED", true),
		},
		ratelimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20), //test
			Timeframe:            time.Minute * 2,
			Enabled:              env.GetBool("RATELIMITER_ENABLED", true),
		},
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
	}
	//logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")
	store := store.NewPostgresStorage(db)

	//redis init
	var redisClient *redis.Client
	if cfg.redisCfg.enabled {
		redisClient = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.password, cfg.redisCfg.db)
		logger.Info("redis connection established")
	}
	cache := cache.NewRedisStorage(redisClient)

	//mailer
	// mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	mailtrap, err := mailer.NewMailTrapClient(cfg.mail.username, cfg.mail.fromEmail, cfg.mail.password)
	if err != nil {
		log.Fatal(err)
	}

	//ratelimiter
	ratelimiter := ratelimiter.NewRedisFixedWindowRateLimiter(
		cache, cfg.ratelimiter.RequestsPerTimeFrame, cfg.ratelimiter.Timeframe,
	)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	app := &application{
		config:        cfg,
		storage:       store,
		l:             logger,
		mailer:        mailtrap,
		authenticator: jwtAuthenticator,
		cache:         cache,
		rateLimiter:   ratelimiter,
	}

	mux := app.mount()
	app.start(mux)
}

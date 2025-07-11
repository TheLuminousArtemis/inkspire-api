package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/theluminousartemis/inkspire/docs"
	"github.com/theluminousartemis/inkspire/internal/auth"
	"github.com/theluminousartemis/inkspire/internal/mailer"
	"github.com/theluminousartemis/inkspire/internal/ratelimiter"
	"github.com/theluminousartemis/inkspire/internal/store"
	"github.com/theluminousartemis/inkspire/internal/store/cache"
	"go.uber.org/zap"
)

type application struct {
	config        config
	storage       store.Storage
	l             *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	cache         cache.Storage
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	addr        string
	apiURL      string
	db          *dbConfig
	env         string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisCfg    redisConfig
	ratelimiter ratelimiter.Config
}

type redisConfig struct {
	addr     string
	password string
	db       int
	enabled  bool
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type mailConfig struct {
	exp       time.Duration
	fromEmail string
	username  string
	password  string
}

type jwtConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type authConfig struct {
	basic basicConfig
	token jwtConfig
}

type basicConfig struct {
	user string
	pass string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:4000")},
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	if app.config.ratelimiter.Enabled {
		r.Use(app.RateLimiterMiddleware)
	}
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.HealthCheck)
		r.With(app.BasicAuthMiddleware()).Get("/metrics", expvar.Handler().ServeHTTP)
		docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.createCommentHandler)
					r.Route("/{commentID}", func(r chi.Router) {
						r.Use(app.commentsContextMiddleware)
						r.Delete("/", app.checkcommentOwnership("admin", app.deleteCommentHandler))
					})
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		//users auth & registration
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/token", app.createTokenHandler)
			r.Post("/user", app.registerUserHandler)
		})
	})
	return r
}

func (app *application) start(mux http.Handler) error {
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		app.l.Infow("signal caught", "signal", s.String())
		shutdown <- srv.Shutdown(ctx)
	}()

	app.l.Infow("Starting server at", zap.String("addr", app.config.addr), zap.String("env", app.config.env))
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdown
	if err != nil {
		return err
	}
	app.l.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)
	return nil
}

// if err != nil {
// 	log.Fatalf("Error starting server: %v", err)
// }

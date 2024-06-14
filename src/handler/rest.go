package handler

import (
	"context"
	"fmt"
	"loverly/lib/jwt"
	"loverly/src/business/usecase"
	"loverly/src/config"
	"net/http"
	"sync"
	"time"

	"loverly/lib/log"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
)

var (
	once   = &sync.Once{}
	Verify = validator.New()
	Log    log.Interface
)

func Init(ctx context.Context, log log.Interface, cfg config.Configuration, uc *usecase.Usecases, jwt *jwt.TokenProvider) {

	once.Do(func() {
		address := fmt.Sprintf(":%d", cfg.BindAddress)
		log.Info(ctx, fmt.Sprintf("Starting loverly service on %s", address))

		r := chi.NewRouter()
		r.Use(chimiddleware.Recoverer)
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		r.Use(chimiddleware.Logger)
		r.Use(chimiddleware.RealIP)
		r.Use(chimiddleware.Timeout(60 * time.Second))
		r.Use(addFieldsToContext)
		r.Use(bodyLogger(log))

		// Initialize routes
		Router(r, uc, jwt)

		// Initalize Log
		Log = log

		err := http.ListenAndServe(address, r)
		if err != nil {
			log.Error(ctx, fmt.Sprintf("ListenAndServe err %s", err))
		}
	})
}

func Router(r *chi.Mux, usecase *usecase.Usecases, jwt *jwt.TokenProvider) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/v1", func(v1 chi.Router) {
		// Authentication
		v1.Post("/login", SignIn(usecase))
		v1.Post("/register", SignUp(usecase))

		auth := v1.With(authentication(jwt, Log))

		// dating in action
		auth.Get("/discovery", Discovery(usecase))
		auth.Get("/match", Match(usecase))
		auth.Post("/swipe", Swipe(usecase))

		// profile
		auth.Get("/profile", GetProfile(usecase))

		// subscription
		auth.Post("/subscription", Subscribe(usecase))
		auth.Get("/subscription", GetSubscribe(usecase))

	})

}

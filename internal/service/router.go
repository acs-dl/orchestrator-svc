package service

import (
	"context"

	auth "github.com/acs-dl/auth-svc/middlewares"
	"github.com/acs-dl/orchestrator-svc/internal/data/postgres"
	"github.com/acs-dl/orchestrator-svc/internal/receiver"
	"github.com/acs-dl/orchestrator-svc/internal/sender"
	"github.com/acs-dl/orchestrator-svc/internal/service/handlers"
	"github.com/acs-dl/orchestrator-svc/internal/service/helpers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()
	ctx := context.Background()

	s.startSender(ctx)
	s.startListener(ctx)

	secret := s.cfg.JwtParams().Secret

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxLog(s.log),
			helpers.CtxModulesQ(postgres.NewModuleQ(s.cfg.DB().Clone())),
			helpers.CtxRequestsQ(postgres.NewRequestsQ(s.cfg.DB().Clone())),
			helpers.CtxRequestTransactionsQ(postgres.NewRequestTransactionsQ(s.cfg.DB().Clone())),
			helpers.CtxSender(sender.NewSender(s.cfg)),
			helpers.CtxRawDB(s.cfg.RawDB()),
			helpers.CtxPublisher(s.cfg.Publisher()),
			helpers.CtxSubscriber(s.cfg.Subscriber()),
			helpers.CtxAmqpCfg(s.cfg.Amqp()),
		),
	)

	r.Route("/integrations/orchestrator", func(r chi.Router) {
		r.Route("/health", func(r chi.Router) {
			r.Get("/live", handlers.CheckHealthLive)
			r.Get("/ready", handlers.CheckHealthReady)
			r.Get("/status", handlers.CheckHealthStatus)
		})

		r.With(auth.Jwt(secret, "orchestrator", []string{"read", "write"}...)).
			Post("/estimate_refresh", handlers.GetEstimatedRefreshTime)

		r.With(auth.Jwt(secret, "orchestrator", []string{"write"}...)).
			Post("/refresh", handlers.RefreshAllModules)

		r.Route("/modules", func(r chi.Router) {
			r.Post("/", handlers.RegisterModule)           // comes from modules
			r.Delete("/{name}", handlers.UnregisterModule) // comes from modules

			r.With(auth.Jwt(secret, "orchestrator", []string{"read", "write"}...)).
				Get("/", handlers.GetModules)
		})

		r.With(auth.Jwt(secret, "orchestrator", []string{"read", "write"}...)).
			Get("/role", handlers.GetRole)

		r.With(auth.Jwt(secret, "orchestrator", []string{"read", "write"}...)).
			Get("/roles", handlers.GetRoles)

		r.With(auth.Jwt(secret, "orchestrator", []string{"read", "write"}...)).
			Get("/submodule", handlers.CheckSubmodule)

		r.Route("/requests", func(r chi.Router) {
			r.With(auth.Jwt(secret, "orchestrator", []string{"write"}...)).
				Post("/", handlers.CreateRequest)
			r.With(auth.Jwt(secret, "orchestrator", []string{"read", "write"}...)).
				Get("/", handlers.GetRequests)
			r.With(auth.Jwt(secret, "orchestrator", []string{"read", "write"}...)).
				Get("/{id}", handlers.GetRequest)
		})

		r.Route("/users", func(r chi.Router) {
			r.With(auth.Jwt(secret, "orchestrator", []string{"read", "write"}...)).
				Get("/{id}", handlers.GetUserById)
			r.With(auth.Jwt(secret, "orchestrator", []string{"read", "write"}...)).
				Delete("/", handlers.DeleteUserById)
		})
	})

	return r
}

func (s *service) startListener(ctx context.Context) error {
	s.log.Info("Starting listener")

	receiver.NewReceiver(s.cfg).Run(ctx)
	return nil
}

func (s *service) startSender(ctx context.Context) error {
	s.log.Info("Starting sender")
	sender.NewSender(s.cfg).Run(ctx)
	return nil
}

package service

import (
	"context"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/receiver"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/sender"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/handlers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()
	ctx := context.Background()

	s.startSender(ctx)
	s.startListener(ctx)

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxLog(s.log),
			helpers.CtxModulesQ(s.modulesQ),
			helpers.CtxRequestsQ(s.requestsQ),
		),
	)
	r.Route("/integrations/orchestrator", func(r chi.Router) {
		r.Route("/modules", func(r chi.Router) {
			r.Post("/", handlers.RegisterModule)
			r.Delete("/{name}", handlers.UnregisterModule)
		})
		r.Route("/requests", func(r chi.Router) {
			r.Post("/", handlers.CreateRequest)
		})
	})

	return r
}

func (s *service) startListener(ctx context.Context) error {
	s.log.Info("Starting listener")
	receiver.NewReceiver(s.subscriber, s.modulesQ, s.requestsQ).Run(ctx)
	return nil
}

func (s *service) startSender(ctx context.Context) error {
	s.log.Info("Starting sender")
	sender.NewSender(s.publisher, s.requestsQ).Run(ctx)
	return nil
}
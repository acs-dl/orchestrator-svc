package service

import (
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data/postgres"
	"net"
	"net/http"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/config"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log        *logan.Entry
	copus      types.Copus
	listener   net.Listener
	modulesQ   data.ModuleQ
	requestsQ  data.RequestQ
	publisher  *amqp.Publisher
	subscriber *message.Subscriber
}

func (s *service) run() error {
	s.log.Info("Service started")
	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) *service {
	return &service{
		log:        cfg.Log(),
		copus:      cfg.Copus(),
		listener:   cfg.Listener(),
		modulesQ:   postgres.NewModuleQ(cfg.DB().Clone()),
		requestsQ:  postgres.NewRequestsQ(cfg.DB().Clone()),
		publisher:  cfg.Publisher(),
		subscriber: cfg.Subscriber(),
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(); err != nil {
		panic(err)
	}
}

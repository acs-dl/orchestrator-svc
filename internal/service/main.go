package service

import (
	"database/sql"
	"net"
	"net/http"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data/postgres"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/config"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log                  *logan.Entry
	copus                types.Copus
	listener             net.Listener
	modulesQ             data.ModuleQ
	requestsQ            data.RequestQ
	requestTransactionsQ data.RequestTransactions
	publisher            *amqp.Publisher
	subscriber           *amqp.Subscriber
	jwt                  *config.JwtCfg
	rawDB                *sql.DB
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
		log:                  cfg.Log(),
		copus:                cfg.Copus(),
		listener:             cfg.Listener(),
		modulesQ:             postgres.NewModuleQ(cfg.DB().Clone()),
		requestsQ:            postgres.NewRequestsQ(cfg.DB().Clone()),
		requestTransactionsQ: postgres.NewRequestTransactionsQ(cfg.DB().Clone()),
		publisher:            cfg.Publisher(),
		subscriber:           cfg.Subscriber(),
		jwt:                  cfg.JwtParams(),
		rawDB:                cfg.RawDB(),
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(); err != nil {
		panic(err)
	}
}

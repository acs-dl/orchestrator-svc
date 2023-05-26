package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AmqpCfg struct {
	Auth         string `figure:"auth,required"`
	Orchestrator string `figure:"orchestrator,required"`
}

func (c *config) Amqp() *AmqpCfg {
	return c.amqp.Do(func() interface{} {
		var cfg AmqpCfg
		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "amqp")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out amqp params from config"))
		}

		return &cfg
	}).(*AmqpCfg)
}

package config

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

type PublisherConfig struct {
	AmqpUrl string `fig:"amqp_url,required"`
}

func (c *config) Publisher() *message.Publisher {
	return c.publisher.Do(func() interface{} {
		var cfg PublisherConfig

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "publisher")).
			Please()

		if err != nil {
			panic(errors.Wrap(err, "failed to figure out publisher config"))
		}

		amqpConfig := amqp.NewDurablePubSubConfig(cfg.AmqpUrl, nil)
		// TODO: implement custom logger
		watermillLogger := watermill.NewStdLogger(false, false)

		publisher, err := amqp.NewPublisher(amqpConfig, watermillLogger)
		if err != nil {
			panic(errors.Wrap(err, "failed to create publisher"))
		}

		return publisher
	}).(*message.Publisher)
}

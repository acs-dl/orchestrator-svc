package config

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type SubscriberConfig struct {
	AmqpUrl string `fig:"amqp_url,required"`
}

func (c *config) Subscriber() *amqp.Subscriber {
	return c.subscriber.Do(func() interface{} {
		var cfg SubscriberConfig

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "subscriber")).
			Please()

		if err != nil {
			panic(errors.Wrap(err, "failed to figure out subscriber config"))
		}

		amqpConfig := amqp.NewDurablePubSubConfig(cfg.AmqpUrl, amqp.GenerateQueueNameTopicNameWithSuffix("SUFFIX"))
		if err != nil {
			panic(errors.Wrap(err, "failed to create subscriber config"))
		}

		// TODO: implement custom logger
		watermillLogger := watermill.NewStdLogger(false, false)

		subscriber, err := amqp.NewSubscriber(amqpConfig, watermillLogger)
		if err != nil {
			panic(errors.Wrap(err, "failed to create subscriber"))
		}

		return subscriber
	}).(*amqp.Subscriber)
}

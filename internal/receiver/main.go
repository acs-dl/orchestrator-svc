package receiver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"time"
)

const serviceName = "receiver"

type Receiver struct {
	subscriber *amqp.Subscriber
	log        *logan.Entry
	modulesQ   data.ModuleQ
	requestsQ  data.RequestQ
}

func NewReceiver(subscriber *amqp.Subscriber, modulesQ data.ModuleQ, requestsQ data.RequestQ) *Receiver {
	return &Receiver{
		subscriber: subscriber,
		log:        logan.New().WithField("service", serviceName),
		modulesQ:   modulesQ,
		requestsQ:  requestsQ,
	}
}

func (r *Receiver) Run(ctx context.Context) {
	go running.WithBackOff(ctx, r.log,
		serviceName,
		r.listenMessages,
		30*time.Second,
		30*time.Second,
		30*time.Second,
	)
}

func (r *Receiver) listenMessages(ctx context.Context) error {
	r.log.Debug("started listening messages")
	modules, err := r.modulesQ.Select()
	if err != nil {
		return errors.Wrap(err, "failed to select modules")
	}

	for _, module := range modules {
		r.log.Debug("started listening messages for module ", module.Name)
		r.startSubscriber(ctx, *module.Endpoint)
	}
	return nil
}

func (r *Receiver) startSubscriber(ctx context.Context, topic string) {
	go running.WithBackOff(ctx, r.log,
		fmt.Sprint(serviceName, "_", topic),
		func(ctx context.Context) error {
			return r.subscribeForTopic(ctx, "orchestrator") //topic)
		},
		30*time.Second,
		30*time.Second,
		30*time.Second,
	)
}

func (r *Receiver) subscribeForTopic(ctx context.Context, topic string) error {
	msgChan, err := (*r.subscriber).Subscribe(ctx, topic)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe for topic "+topic)
	}
	r.log.Debug("successfully subscribed for topic ", topic)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgChan:
			r.log.Debug("received message ", msg.UUID)
			err = r.processMessage(msg)
			if err != nil {
				r.log.WithError(err).Error("failed to process message ", msg.UUID)
			} else {
				msg.Ack()
			}
		}
	}
}

func (r *Receiver) processMessage(msg *message.Message) error {
	r.log.Debug("started processing message ", msg.UUID)
	var queueOutput types.QueueOutput
	err := json.Unmarshal(msg.Payload, &queueOutput)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message "+msg.UUID)
	}

	err = r.requestsQ.FilterByIDs(queueOutput.ID).SetStatusError(queueOutput.Status.ToRequestStatus(), queueOutput.Error)
	if err != nil {
		return errors.Wrap(err, "failed to update state for notification: "+queueOutput.ID)
	}

	r.log.Debug("finished processing message ", msg.UUID)
	return nil
}

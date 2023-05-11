package sender

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"time"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/acs-dl/orchestrator-svc/internal/data"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/running"
)

const serviceName = "sender"

type Sender struct {
	publisher *amqp.Publisher
	requestQ  data.RequestQ
	moduleQ   data.ModuleQ
	log       *logan.Entry
}

func NewSender(publisher *amqp.Publisher, requestQ data.RequestQ, moduleQ data.ModuleQ) *Sender {
	return &Sender{
		publisher: publisher,
		requestQ:  requestQ,
		moduleQ:   moduleQ,
		log:       logan.New().WithField("service", serviceName),
	}
}

func (s *Sender) Run(ctx context.Context) error {
	go running.WithBackOff(ctx, s.log,
		serviceName,
		s.processMessages,
		30*time.Second,
		30*time.Second,
		30*time.Second,
	)
	return nil
}

func (s *Sender) processMessages(ctx context.Context) error {
	s.log.Debug("started processing messages")
	messagesToSend, err := s.requestQ.FilterByStatuses(data.CREATED).Select()
	if err != nil {
		return errors.Wrap(err, "failed to select messages")
	}

	for _, message := range messagesToSend {
		s.log.Debug("started processing notification with id ", message.ID)

		module, err := s.moduleQ.FilterByNames(message.ModuleName).Get()
		if err != nil {
			return errors.Wrap(err, "failed to get full module for notification: "+message.ID)
		}

		if module == nil {
			return errors.Errorf("no module was found for notification:" + message.ID)
		}

		err = s.publisher.Publish(module.Topic, s.buildMessage(message))
		if err != nil {
			return errors.Wrap(err, "failed to process notification: "+message.ID)
		}

		err = s.requestQ.FilterByIDs(message.ID).SetStatus(data.PENDING)
		if err != nil {
			return errors.Wrap(err, "failed to update state for notification: "+message.ID)
		}
		s.log.Debug("finished processing notification with id ", message.ID)
	}

	s.log.Debug("finished processing messages")
	return nil
}

func (s *Sender) buildMessage(request data.Request) *message.Message {
	return &message.Message{
		UUID:     request.ID,
		Metadata: nil,
		Payload:  message.Payload(request.Payload),
	}
}

func (s *Sender) SendMessageToCustomChannel(topic string, msg *message.Message) error {
	err := (*s.publisher).Publish(topic, msg)
	if err != nil {
		s.log.WithError(err).Errorf("failed to send msg `%s to `%s`", msg.UUID, topic)
		return errors.Wrap(err, "failed to send msg: "+msg.UUID)
	}

	return nil
}

func (s *Sender) BuildPermissionsMessage(uuid string, payload []byte) *message.Message {
	return &message.Message{
		UUID:     uuid,
		Metadata: nil,
		Payload:  payload,
	}
}

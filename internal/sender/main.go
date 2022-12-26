package sender

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"time"
)

const serviceName = "sender"

type Sender struct {
	publisher *message.Publisher
	requestQ  data.RequestQ
	log       *logan.Entry
}

func NewSender(publisher *message.Publisher, requestQ data.RequestQ) *Sender {
	return &Sender{
		publisher: publisher,
		requestQ:  requestQ,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	running.WithBackOff(ctx, s.log,
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
		err = (*s.publisher).Publish(*message.Module.Endpoint, s.buildMessage(message))
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

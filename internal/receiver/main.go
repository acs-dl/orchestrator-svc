package receiver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/processor"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/sender"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const serviceName = "receiver"

type Receiver struct {
	subscriber           *amqp.Subscriber
	log                  *logan.Entry
	modulesQ             data.ModuleQ
	requestsQ            data.RequestQ
	requestTransactionsQ data.RequestTransactions
	processor            processor.Processor
}

func NewReceiver(subscriber *amqp.Subscriber, modulesQ data.ModuleQ, requestsQ data.RequestQ, sender *sender.Sender, requestTransactionsQ data.RequestTransactions) *Receiver {
	return &Receiver{
		subscriber:           subscriber,
		log:                  logan.New().WithField("service", serviceName),
		modulesQ:             modulesQ,
		requestsQ:            requestsQ,
		requestTransactionsQ: requestTransactionsQ,
		processor:            processor.NewProcessor(modulesQ, sender),
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
	return r.subscribeForTopic(ctx, "orchestrator") //topic)
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
			}
			msg.Ack()
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

	if queueOutput.Action != nil {
		err = r.processor.HandleNewMessage(queueOutput)
		if err != nil {
			return errors.Wrap(err, "failed to handle message")
		}

		return nil
	}

	err = r.requestsQ.FilterByIDs(queueOutput.ID).SetStatusError(queueOutput.Status.ToRequestStatus(), queueOutput.Error)
	if err != nil {
		return errors.Wrap(err, "failed to update state for notification: "+queueOutput.ID)
	}

	transaction, err := r.requestTransactionsQ.FilterByRequestID(queueOutput.ID).Get()
	if err != nil {
		return errors.Wrap(err, "failed to get request transaction")
	}
	if transaction == nil {
		return errors.New("something wrong with transaction")
	}

	var transactionRequests map[string]bool
	err = json.Unmarshal(transaction.Requests, &transactionRequests)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal transaction requests")
	}

	transactionRequests[queueOutput.ID] = true
	containsUnhandled := false
	for _, isHandled := range transactionRequests {
		if !isHandled {
			containsUnhandled = true
		}
	}

	if !containsUnhandled {
		switch transaction.Action {
		case data.Single:
			err = r.handleSingleAction(*transaction)
			if err != nil {
				return errors.Wrap(err, "failed to handle single request")
			}
		case data.DeleteUserAction:
			err = r.handleDeleteUserAction(queueOutput, *transaction)
			if err != nil {
				return errors.Wrap(err, "failed to handle delete user request")
			}
		}
		return nil
	}

	transaction.Requests, err = json.Marshal(transactionRequests)
	if err != nil {
		return errors.Wrap(err, "failed to marshal transaction requests")
	}

	err = r.requestTransactionsQ.FilterByIDs(queueOutput.ID).Update(*transaction)
	if err != nil {
		return errors.Wrap(err, "failed to update transaction request")
	}

	r.log.Debug("finished processing message ", msg.UUID)
	return nil
}

func (r *Receiver) handleSingleAction(transaction data.RequestTransaction) error {
	err := r.requestTransactionsQ.FilterByIDs(transaction.ID).Delete()
	if err != nil {
		return errors.Wrap(err, "failed to delete request transaction")
	}

	return nil
}

func (r *Receiver) handleDeleteUserAction(queueOutput types.QueueOutput, transaction data.RequestTransaction) error {
	module, err := r.modulesQ.FilterByNames("identity").Get()
	if err != nil {
		return errors.Wrap(err, "failed to get identity module")
	}

	if module == nil {
		return errors.Errorf("no module with name `identity`")
	}

	req, err := r.requestsQ.FilterByIDs(queueOutput.ID).Get()
	if err != nil {
		return errors.Wrap(err, "failed to get request")
	}

	if req == nil {
		return errors.Errorf("no request with such id")
	}

	err = helpers.MakeNoResponseRequest(data.RequestParams{
		Method: http.MethodDelete,
		Link:   fmt.Sprintf("%s/orchestrator_users/%d", module.Link, req.ToUserID),
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    nil,
		Query:   nil,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		return errors.Wrap(err, "failed to make delete user request")
	}

	err = r.requestTransactionsQ.FilterByIDs(transaction.ID).Delete()
	if err != nil {
		return errors.Wrap(err, "failed to delete request transaction")
	}

	return nil
}

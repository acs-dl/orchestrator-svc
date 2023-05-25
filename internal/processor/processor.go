package processor

import (
	"github.com/acs-dl/orchestrator-svc/internal/data"
	"github.com/acs-dl/orchestrator-svc/internal/sender"
	"github.com/acs-dl/orchestrator-svc/internal/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	serviceName = data.ModuleName + "-processor"

	//add needed actions for module
	SetModulesPermissionsAction = "set_modules_permissions"
	GetModulesPermissionsAction = "get_modules_permissions"
)

type Processor interface {
	HandleNewMessage(msg types.QueueOutput) error
}

type processor struct {
	log       *logan.Entry
	modulesQ  data.ModuleQ
	sender    *sender.Sender
	authTopic string
}

var handleActions = map[string]func(proc *processor, msg types.QueueOutput) error{
	GetModulesPermissionsAction: (*processor).handleGetModulesPermissions,
}

func NewProcessor(modulesQ data.ModuleQ, sender *sender.Sender, authTopic string) Processor {
	return &processor{
		log:       logan.New().WithField("service", serviceName),
		modulesQ:  modulesQ,
		sender:    sender,
		authTopic: authTopic,
	}
}

func (p *processor) HandleNewMessage(msg types.QueueOutput) error {
	p.log.Infof("handling message with id `%s`", *msg.RequestId)

	err := validation.Errors{
		"action": validation.Validate(msg.Action, validation.Required, validation.In(GetModulesPermissionsAction)),
	}.Filter()
	if err != nil {
		p.log.WithError(err).Errorf("no such action to handle for message with id `%s`", *msg.RequestId)
		return errors.Errorf("no such action `%s` to handle for message with id `%s`", *msg.Action, *msg.RequestId)
	}

	requestHandler := handleActions[*msg.Action]
	if err = requestHandler(p, msg); err != nil {
		p.log.WithError(err).Errorf("failed to handle message with id `%s`", *msg.RequestId)
		return err
	}

	p.log.Infof("finish handling message with id `%s`", *msg.RequestId)
	return nil
}

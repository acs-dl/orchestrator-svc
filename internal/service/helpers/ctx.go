package helpers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/sender"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	modulesQCtxKey
	requestsQCtxKey
	senderCtxKey
	rawDBCtxKey
	publisherCtxKey
	subscriberCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxModulesQ(entry data.ModuleQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, modulesQCtxKey, entry)
	}
}

func ModulesQ(r *http.Request) data.ModuleQ {
	return r.Context().Value(modulesQCtxKey).(data.ModuleQ)
}

func CtxRequestsQ(entry data.RequestQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, requestsQCtxKey, entry)
	}
}

func RequestsQ(r *http.Request) data.RequestQ {
	return r.Context().Value(requestsQCtxKey).(data.RequestQ)
}

func CtxSender(entry *sender.Sender) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, senderCtxKey, entry)
	}

}

func Sender(r *http.Request) *sender.Sender {
	return r.Context().Value(senderCtxKey).(*sender.Sender)
}

func CtxRawDB(entry *sql.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, rawDBCtxKey, entry)
	}

}

func RawDB(r *http.Request) *sql.DB {
	return r.Context().Value(rawDBCtxKey).(*sql.DB)
}

func CtxPublisher(entry *amqp.Publisher) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, publisherCtxKey, entry)
	}

}

func Publisher(r *http.Request) *amqp.Publisher {
	return r.Context().Value(publisherCtxKey).(*amqp.Publisher)
}

func CtxSubscriber(entry *amqp.Subscriber) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, subscriberCtxKey, entry)
	}

}

func Subscriber(r *http.Request) *amqp.Subscriber {
	return r.Context().Value(subscriberCtxKey).(*amqp.Subscriber)
}

package helpers

import (
	"context"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	modulesQCtxKey
	requestsQCtxKey
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

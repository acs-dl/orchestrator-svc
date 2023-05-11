package handlers

import (
	"net/http"

	"github.com/acs-dl/orchestrator-svc/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
)

func CheckHealthStatus(w http.ResponseWriter, r *http.Request) {
	if err := helpers.RawDB(r).Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ape.Render(w, http.StatusInternalServerError)
		return
	}

	if !helpers.Subscriber(r).IsConnected() {
		w.WriteHeader(http.StatusInternalServerError)
		ape.Render(w, http.StatusInternalServerError)
		return
	}

	if !helpers.Publisher(r).IsConnected() {
		w.WriteHeader(http.StatusInternalServerError)
		ape.Render(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	ape.Render(w, http.StatusOK)
}

package handlers

import (
	"net/http"
	"time"

	"gitlab.com/distributed_lab/ape"
)

func CheckHealthReady(w http.ResponseWriter, r *http.Request) {
	//started := r.Context().Value("started").(time.Time) //In `ape.LoganMiddleware` we have `started` var, so we can save it in ctx and use here
	started := time.Now()
	duration := time.Now().Sub(started)
	if duration.Seconds() > 10 {
		w.WriteHeader(http.StatusInternalServerError)
		ape.Render(w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	ape.Render(w, http.StatusOK)
}

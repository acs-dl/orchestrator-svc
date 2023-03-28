package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
)

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	ape.Render(w, http.StatusOK)
}

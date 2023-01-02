package handlers

import "net/http"

func UnregisterModule(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusAccepted)
}

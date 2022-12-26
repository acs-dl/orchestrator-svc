package handlers

import "net/http"

func DeleteModule(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusAccepted)
}

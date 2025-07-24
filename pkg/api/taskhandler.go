package api

import (
	"net/http"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodGet:
		GetTaskHandler(w, r)
	case http.MethodPut:
		PutTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

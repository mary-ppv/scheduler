package api

import (
	"encoding/json"
	"final/pkg/db"
	"net/http"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		SendError(w, "can not get id", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		SendError(w, "can not get task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

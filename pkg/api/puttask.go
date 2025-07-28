package api

import (
	"encoding/json"
	"final/pkg/db"
	"fmt"
	"io"
	"net/http"
)

func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var task db.Task
	if err := json.Unmarshal(body, &task); err != nil {
		SendError(w, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		SendError(w, "ID is required", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		SendError(w, "title is required", http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		SendError(w, "can not update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}

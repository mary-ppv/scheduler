package api

import (
	"encoding/json"
	"final/pkg/db"
	"log"
	"net/http"
	"strconv"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		SendError(w, "invalid JSON", http.StatusBadRequest)
		log.Printf("request structure mismatch: %v", err)
		return
	}

	if task.Title == "" {
		SendError(w, "title is required", http.StatusBadRequest)
		return
	}

	now := time.Now()

	if err := db.CheckDate(&task, now); err != nil {
		SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		SendError(w, "failed to add task", http.StatusInternalServerError)
		log.Printf("database error: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": strconv.FormatInt(id, 10)})
}

func SendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

package api

import (
	"encoding/json"
	"final/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		SendError(w, "Invalid JSON")
		return
	}

	if task.Title == "" {
		SendError(w, "title is empty")
		return
	}

	now := time.Now()

	if task.Date == "today" {
		task.Date = now.Format("20060102")
	}

	err = db.CheckDate(&task, now)
	if err != nil {
		SendError(w, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		SendError(w, err.Error())
		return
	}

	idStr := strconv.Itoa(int(id))
	response := map[string]string{"id": idStr}
	json.NewEncoder(w).Encode(response)
}

func SendError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

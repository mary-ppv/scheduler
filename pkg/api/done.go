package api

import (
	"encoding/json"
	"final/pkg/db"
	"net/http"
	"time"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		SendError(w, "field id is empty")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		SendError(w, err.Error())
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			SendError(w, err.Error())
			return
		}
	} else {
		now, err := time.Parse("20060102", task.Date)
		if err != nil {
			SendError(w, "Неверный формат даты")
			return
		}

		nextDateStr, err := db.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			SendError(w, err.Error())
			return
		}

		err = db.UpdateDate(nextDateStr, id)
		if err != nil {
			SendError(w, err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}

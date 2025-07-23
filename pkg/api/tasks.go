package api

import (
	"encoding/json"
	"final/pkg/db"
	"fmt"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

type TaskResp struct {
	Task *db.Task `json:"task"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	queryParams := r.URL.Query()
	searchFilter := queryParams.Get("search")

	limit := 50
	tasks, err := db.Tasks(limit, searchFilter)
	if err != nil {
		SendError(w, fmt.Sprintf("failed to get tasks: %v", err))
		return
	}

	resp := TasksResp{
		Tasks: tasks,
	}

	json.NewEncoder(w).Encode(resp)
}

package app

import (
	"final/pkg/api"
	"net/http"

	"final/pkg/db"
)

func Init() {
	http.HandleFunc("/api/nextdate", db.NextDateHandler)
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/api/tasks", api.TasksHandler)
	http.HandleFunc("/api/task/done", api.DoneTaskHandler)
	http.HandleFunc("/api/task/delete", api.DeleteTaskHandler)
	http.HandleFunc("/api/signin", api.SigninHandler)
}

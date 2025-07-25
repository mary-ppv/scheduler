package app

import (
	"final/pkg/api"
	"final/pkg/db"
	"final/pkg/server"

	"net/http"
)

func Init() {
	http.HandleFunc("/api/task/done", server.AuthMiddleware((api.DoneTaskHandler)))
	http.HandleFunc("/api/nextdate", db.NextDateHandler)

	http.HandleFunc("/api/task", server.AuthMiddleware((api.TaskHandler)))
	http.HandleFunc("/api/tasks", server.AuthMiddleware((api.TasksHandler)))

	http.HandleFunc("/api/signin", server.SignInHandler)
}

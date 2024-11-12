package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type task struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
}

var tasks []task
var taskId int

func Add(w http.ResponseWriter, r *http.Request) {
	var newtask task
	err := json.NewDecoder(r.Body).Decode(&newtask)
	if err != nil || newtask.Desc == "" {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	taskId++
	newtask.Id = taskId
	tasks = append(tasks, newtask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newtask)
}
func List(w http.ResponseWriter, r *http.Request) {
	if len(tasks) == 0 {
		http.Error(w, "No tasks found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(tasks)
	json.NewEncoder(w).Encode(map[string]interface{}{"count": len(tasks), "tasks": tasks, "message": "Success"})
}
func Update(w http.ResponseWriter, r *http.Request) {
	var modtask task
	err := json.NewDecoder(r.Body).Decode(&modtask)
	if err != nil || modtask.Desc == "" || modtask.Id == 0 {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	for i, t := range tasks {
		if t.Id == modtask.Id {
			tasks[i].Desc = modtask.Desc
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(modtask)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}
func Delete(w http.ResponseWriter, r *http.Request) {
	var deltask task
	err := json.NewDecoder(r.Body).Decode(&deltask)
	if err != nil || deltask.Id == 0 {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	for i, t := range tasks {
		if t.Id == deltask.Id {
			tmp1 := tasks[:i]
			tmp2 := tasks[i+1:]
			tasks = append(tmp1, tmp2...)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", Add)
		r.Get("/", List)
		r.Put("/", Update)
		r.Delete("/", Delete)
	})
	fmt.Println("Server running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))

}

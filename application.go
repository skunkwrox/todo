package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"./data"
)

func main() {

	data.Init()
	log.Println("Configuring web server")

	r := mux.NewRouter()
	r.HandleFunc("/person/", data.GetListPersons).Methods("GET")
	r.HandleFunc("/person/", data.AddPersonDetails).Methods("POST")
	r.HandleFunc("/person/{person_id:[0-9]+}/", data.GetPersonDetails).Methods("GET")
	r.HandleFunc("/person/{person_id:[0-9]+}/", data.UpdatePersonDetails).Methods("PUT")
	r.HandleFunc("/person/{person_id:[0-9]+}/", data.DeletePersonDetails).Methods("DELETE")

	r.HandleFunc("/task/", data.GetListTasks).Methods("GET")
	r.HandleFunc("/task/", data.AddTaskDetails).Methods("POST")
	r.HandleFunc("/task/{task_id:[0-9]+}/", data.GetTaskDetails).Methods("GET")
	r.HandleFunc("/task/{task_id:[0-9]+}/", data.UpdateTaskDetails).Methods("PUT")
	r.HandleFunc("/task/{task_id:[0-9]+}/", data.DeleteTaskDetails).Methods("DELETE")

	log.Println("Starting web server")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}

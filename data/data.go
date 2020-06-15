package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS person (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS task (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
	priority INTEGER NULL,
	created timestamp NOT NULL,
	last_updated timestamp NOT NULL,
	assigned_to INTEGER NULL REFERENCES person (id),
	due_by timestamp without time zone NULL
);

CREATE TABLE IF NOT EXISTS init (init bool);
`

type PersonDetails struct {
	Name  string `db:"name"`
	Email string `db:"email"`
}
type Person struct {
	Id int32 `db:"id"`
	PersonDetails
}

type TaskDb struct {
	Id          int32         `db:"id"`
	Title       string        `db:"title"`
	Description string        `db:"description"`
	DueBy       sql.NullTime  `db:"due_by"`
	Created     time.Time     `db:"created"`
	LastUpdated time.Time     `db:"last_updated"`
	AssignedTo  sql.NullInt32 `db:"assigned_to"`
	Priority    sql.NullInt32 `db:"priority"`
}

type TaskDetails struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueBy       *time.Time `json:"due_by,omitempty"`
	Created     time.Time  `json:"created"`
	LastUpdated time.Time  `json:"last_updated"`
	AssignedTo  *int32     `json:"assigned_to,omitempty"`
	Priority    *int32     `json:"priority,omitempty"`
}

type Task struct {
	Id int32 `db:"id"`
	TaskDetails
}

var db *sqlx.DB

func Init() {
	db = connectToDatabase()
	createSchema(db)
	addData(db)
}

func connectToDatabase() (*sqlx.DB) {
	log.Println("Connecting to the database")
	db, err := sqlx.Connect("postgres", "user=todo_app dbname=todo sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
func createSchema(db *sqlx.DB) sql.Result {
	log.Println("Adding table definitions (if need be)")
	return db.MustExec(schema)
}

func addData(db *sqlx.DB) {
	count := 0
	err := db.Get(&count, "SELECT COUNT(*) FROM init")
	if err != nil {
		log.Fatalln(err)
	}
	if count < 1 {
		log.Println("Adding sample data to DB")
		insertData(db)
	} else {
		log.Println("DB already initialized")
	}
}

func insertData(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO person (name, email) VALUES ($1, $2)", "Greg Sample", "greg.sample@todo.net")
	tx.MustExec("INSERT INTO person (name, email) VALUES ($1, $2)", "Jeff Sample", "jeff.sample@todo.net")
	tx.MustExec("INSERT INTO person (name, email) VALUES ($1, $2)", "Mike Sample", "mike.sample@todo.net")

	tx.MustExec("INSERT INTO task (title, description, priority, created, last_updated) VALUES ($1, $2, $3, $4, $5)",
		"Create DB", "Create database to hold to-dos data", 1, time.Now(), time.Now())

	tx.MustExec("INSERT INTO task (title, description, priority, created, last_updated) VALUES ($1, $2, $3, $4, $5)",
		"Add DB initialization", "Add code to check if the DB is initialized, if not add sample data to the DB", 2, time.Now(), time.Now())

	tx.MustExec("INSERT INTO task (title, description, priority, created, last_updated) VALUES ($1, $2, $3, $4, $5)",
		"Add access methods for Person", "Add access methods to the Person data", 3, time.Now(), time.Now())

	tx.MustExec("INSERT INTO task (title, description, priority, created, last_updated) VALUES ($1, $2, $3, $4, $5)",
		"Add access methods for Task", "Add access methods to the Task data", 3, time.Now(), time.Now())

	tx.MustExec("INSERT INTO task (title, description, created, last_updated) VALUES ($1, $2, $3, $4)",
		"Add NULL handling for tasks", "Tasks have field that can be null add code to deal with it", time.Now(), time.Now())

	tx.MustExec("INSERT INTO init (init) VALUES (True)")
	tx.Commit()
}

func GetListPersons(w http.ResponseWriter, req *http.Request) {
	var persons []Person
	err := db.Select(&persons, "SELECT * FROM person ORDER BY name ASC")
	if err != nil {
		log.Println("Error retrieving persons data")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error finding persons data\n")
		return
	}
	personsJson, err := json.Marshal(persons)
	if err != nil {
		log.Println("Error marshaling persons data")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error finding persons data\n")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(personsJson)
}

func GetPersonDetails(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	personId := vars["person_id"]

	var persons []Person
	err := db.Select(&persons, "SELECT * FROM person WHERE id = $1", personId)
	if err != nil {
		log.Println("Error retrieving person details")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error finding person details\n")
		return
	}
	if len(persons) != 1 {
		log.Println("Error retrieving person details got ", len(persons), " records back for ID ", personId)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error finding person details\n")
		return
	}
	personsJson, err := json.Marshal(persons[0])
	if err != nil {
		log.Println("Error marshaling person details")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error finding person details\n")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(personsJson)
}

func getPersonDetails(w http.ResponseWriter, req *http.Request, operation string) (error, PersonDetails) {
	var personDetails PersonDetails
	err := json.NewDecoder(req.Body).Decode(&personDetails)
	if err != nil {
		log.Println("Error decoding person details")
		http.Error(w, fmt.Sprintf("error %s person", operation), http.StatusBadRequest)
	}
	return err, personDetails
}

func personToJson(w http.ResponseWriter, person Person, operation string) (error, []byte) {
	personsJson, err := json.Marshal(person)
	if err != nil {
		log.Println("Error marshaling person details")
		http.Error(w, fmt.Sprintf("error %s person", operation), http.StatusBadRequest)
	}
	return err, personsJson
}

func sendPersonDetails(w http.ResponseWriter, statusCode int, personJson []byte) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(personJson)
}

func AddPersonDetails(w http.ResponseWriter, req *http.Request) {
	err, person := getPersonDetails(w, req, "adding")
	if err != nil {
		return
	}
	tx := db.MustBegin()
	defer tx.Rollback()

	rows, err := tx.NamedQuery("INSERT INTO person (name, email) VALUES (:name, :email) RETURNING id", &person)
	if err != nil {
		http.Error(w, "Error adding user", http.StatusBadRequest)
		return
	}

	var personId int32
	if rows.Next() {
		rows.Scan(&personId)
	}
	tx.Commit()

	newPerson := Person{personId, person}
	err, personJson := personToJson(w, newPerson, "adding")
	if err != nil {
		return
	}

	sendPersonDetails(w, http.StatusCreated, personJson)
}

func getPersonIdFromMuxVars(w http.ResponseWriter, req *http.Request) (error, int32) {
	vars := mux.Vars(req)
	personIdRaw := vars["person_id"]
	personId, err := strconv.ParseInt(personIdRaw, 10, 32)
	if err != nil {
		http.Error(w, "error updating person", http.StatusBadRequest)
		log.Println("Error updating person person ID ", personIdRaw, " not valid")
	}
	return err, int32(personId)
}

func UpdatePersonDetails(w http.ResponseWriter, req *http.Request) {
	err, personId := getPersonIdFromMuxVars(w, req)
	if err != nil {
		return
	}
	err, personDetails := getPersonDetails(w, req, "updating")
	if err != nil {
		return
	}
	person := Person{int32(personId), personDetails}
	tx := db.MustBegin()
	_, err = tx.NamedExec("UPDATE person SET name=:name, email=:email WHERE id=:id", &person)
	if err != nil {
		http.Error(w, "error updating person", http.StatusBadRequest)
		log.Println("Error updating person person ID ", personId, err)
	}

	tx.Commit()

	err, personJson := personToJson(w, person, "updating")
	if err != nil {
		return
	}

	sendPersonDetails(w, http.StatusOK, personJson)
}

func DeletePersonDetails(w http.ResponseWriter, req *http.Request) {
	err, personId := getPersonIdFromMuxVars(w, req)
	if err != nil {
		return
	}

	tx := db.MustBegin()
	person := Person{Id: personId}
	_, err = tx.NamedExec("DELETE FROM person WHERE id=:id", &person)
	if err != nil {
		http.Error(w, "error deleting person", http.StatusBadRequest)
		log.Println("Error deleting person person ID ", personId, err)
	}
	tx.Commit()
}

func GetListTasks(w http.ResponseWriter, req *http.Request) {
	var tasks []TaskDb
	err := db.Select(&tasks, "SELECT * FROM task ORDER BY id ASC")
	if err != nil {
		log.Println("Error retrieving tasks data", err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error finding tasks data\n")
		return
	}
	tasksJson, err := json.Marshal(convertFromTasksDb(tasks))
	if err != nil {
		log.Println("Error marshaling tasks data")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error finding tasks data\n")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(tasksJson)
}

func convertFromTaskDb(taskDb TaskDb) Task {
	taskDetails := TaskDetails{Title: taskDb.Title,
		Description: taskDb.Description,
		Created:     taskDb.Created,
		LastUpdated: taskDb.LastUpdated}

	if taskDb.DueBy.Valid {
		taskDetails.DueBy = &taskDb.DueBy.Time
	}

	if taskDb.AssignedTo.Valid {
		taskDetails.AssignedTo = &taskDb.AssignedTo.Int32
	}

	if taskDb.Priority.Valid {
		taskDetails.Priority = &taskDb.Priority.Int32
	}
	return Task{Id: taskDb.Id, TaskDetails: taskDetails}
}

func convertFromTasksDb(tasksDb []TaskDb) []Task {
	var tasks []Task

	for _, taskDb := range tasksDb {
		tasks = append(tasks, convertFromTaskDb(taskDb))
	}

	return tasks
}

func GetTaskDetails(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	taskId := vars["task_id"]

	var tasksDb []TaskDb
	err := db.Select(&tasksDb, "SELECT * FROM task WHERE id = $1", taskId)
	if err != nil {
		log.Println("Error retrieving task details")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error finding task details\n")
		return
	}
	if len(tasksDb) != 1 {
		log.Println("Error retrieving task details got ", len(tasksDb), " records back for ID ", taskId)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error finding task details\n")
		return
	}
	task := convertFromTaskDb(tasksDb[0])
	tasksJson, err := json.Marshal(task)
	if err != nil {
		log.Println("Error marshaling task details")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error finding task details\n")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(tasksJson)
}

func getTaskDetails(w http.ResponseWriter, req *http.Request, operation string) (error, TaskDetails) {
	var taskDetails TaskDetails
	err := json.NewDecoder(req.Body).Decode(&taskDetails)
	if err != nil {
		log.Println("Error decoding task details")
		http.Error(w, fmt.Sprintf("error %s task", operation), http.StatusBadRequest)
	}
	return err, taskDetails
}

func convertToTaskDb(task TaskDetails) TaskDb {
	taskDb := TaskDb{Title: task.Title,
		Description: task.Description,
		Created:     task.Created,
		LastUpdated: task.LastUpdated,
	}

	if task.DueBy != nil {
		taskDb.DueBy = sql.NullTime{Time: *task.DueBy, Valid: true}
	}

	if task.AssignedTo != nil {
		taskDb.AssignedTo = sql.NullInt32{Int32: *task.AssignedTo, Valid: true}
	}

	if task.Priority != nil {
		taskDb.Priority = sql.NullInt32{Int32: *task.Priority, Valid: true}
	}

	return taskDb
}

func taskToJson(w http.ResponseWriter, task Task, operation string) (error, []byte) {
	personsJson, err := json.Marshal(task)
	if err != nil {
		log.Println("Error marshaling task details")
		http.Error(w, fmt.Sprintf("error %s task", operation), http.StatusBadRequest)
	}
	return err, personsJson
}

func sendTaskDetails(w http.ResponseWriter, statusCode int, taskJson []byte) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(taskJson)
}

func AddTaskDetails(w http.ResponseWriter, req *http.Request) {
	err, task := getTaskDetails(w, req, "adding")
	if err != nil {
		return
	}

	taskDb := convertToTaskDb(task)
	tx := db.MustBegin()
	defer tx.Rollback()

	rows, err := tx.NamedQuery("INSERT INTO task (title, description, priority, created, last_updated, assigned_to, due_by)" +
		"VALUES (:title, :description, :priority, :created, :last_updated, :assigned_to, :due_by) RETURNING id", &taskDb)
	if err != nil {
		http.Error(w, "Error adding task", http.StatusBadRequest)
		return
	}

	var taskId int32
	if rows.Next() {
		rows.Scan(&taskId)
	}
	tx.Commit()

	newTask := Task{taskId, task}
	err, taskJson := taskToJson(w, newTask, "adding")
	if err != nil {
		return
	}

	sendTaskDetails(w, http.StatusCreated, taskJson)
}

func getTaskIdFromMuxVars(w http.ResponseWriter, req *http.Request) (error, int32) {
	vars := mux.Vars(req)
	taskIdRaw := vars["task_id"]
	taskId, err := strconv.ParseInt(taskIdRaw, 10, 32)
	if err != nil {
		http.Error(w, "error updating task", http.StatusBadRequest)
		log.Println("Error updating task task ID ", taskIdRaw, " not valid")
	}
	return err, int32(taskId)
}

func UpdateTaskDetails(w http.ResponseWriter, req *http.Request) {
	err, taskId := getTaskIdFromMuxVars(w, req)
	if err != nil {
		return
	}
	err, taskDetails := getTaskDetails(w, req, "updating")
	if err != nil {
		return
	}
	taskDb := convertToTaskDb(taskDetails)
	taskDb.Id = taskId
	tx := db.MustBegin()
	_, err = tx.NamedExec("UPDATE task " +
		"SET title=:title, " +
		"description=:description, " +
		"priority=:priority, " +
		"last_updated=NOW(), " +
		"assigned_to=:assigned_to, " +
		"due_by=:due_by " +
		"WHERE id=:id", &taskDb)


	if err != nil {
		http.Error(w, "error updating task", http.StatusBadRequest)
		log.Println("Error updating task task ID ", taskId, err)
	}

	tx.Commit()

	err, taskJson := taskToJson(w, convertFromTaskDb(taskDb), "updating")
	if err != nil {
		return
	}

	sendTaskDetails(w, http.StatusOK, taskJson)
}

func DeleteTaskDetails(w http.ResponseWriter, req *http.Request) {
	err, taskId := getTaskIdFromMuxVars(w, req)
	if err != nil {
		return
	}

	tx := db.MustBegin()
	task := Task{Id: taskId}
	_, err = tx.NamedExec("DELETE FROM task WHERE id=:id", &task)
	if err != nil {
		http.Error(w, "error deleting task", http.StatusBadRequest)
		log.Println("Error deleting task task ID ", taskId, err)
	}
	tx.Commit()
}

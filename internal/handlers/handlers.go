package handlers

import (
	"encoding/json"
	"errors"
	"github.com/1set/todotxt"
	"github.com/chrishaynes21/apichallenge/pkg/trace"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// ListTodos outputs the list of todos.  This function should accept
// query params that allow parameterization of the search
func ListTodos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// setup traced logging and context
	ctx := trace.Ctx()
	fields := log.Fields{"traceID": ctx.Value(trace.TIDKey), "func": "ListTodos"}
	log.WithFields(fields).Debug("begin")

	// attempt to load todos file
	todos, err := todotxt.LoadFromPath("todo.txt")
	if err != nil {
		log.WithFields(fields).WithError(err).Error("failed to load todo.txt")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// filter todos
	filtered := filterTodos(ctx, todos, r.URL.Query())

	// sort filtered todos
	sorted, err := sortTodos(ctx, filtered, r.URL.Query())
	if err != nil {
		if errors.Is(err, ErrUnknownSort) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set content type and attempt to encode response
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(sorted); err != nil {
		log.WithFields(fields).WithError(err).Error("failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.WithFields(fields).Debug("success")
}

// GetTodo gets a specific todo
func GetTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// setup traced logging and context
	ctx := trace.Ctx()
	fields := log.Fields{"traceID": ctx.Value(trace.TIDKey), "func": "GetTodo"}
	log.WithFields(fields).Debug("begin")

	// parse id param and attempt string to int conversion
	rawTID := ps.ByName("id")
	todoID, err := strconv.Atoi(rawTID)
	if err != nil {
		log.WithFields(fields).WithField("todoID", rawTID).WithError(err).Error("param conversion failure")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// attempt to load todo file
	todos, err := todotxt.LoadFromPath("todo.txt")
	if err != nil {
		log.WithFields(fields).WithError(err).Error("failed to load todo.txt")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// attempt to get requested task by ID
	task, err := todos.GetTask(todoID)
	if err != nil { // TODO: investigate potential GetTask failure reasons
		log.WithFields(fields).WithField("todoID", todoID).WithError(err).Error("task not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// set content type and attempt to encode to response
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(task); err != nil {
		log.WithFields(fields).WithError(err).Error("failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.WithFields(fields).Debug("success")
}

// UpdateTodo takes the body of the request and updates the todo in todo.txt
func UpdateTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// setup traced logging and context
	ctx := trace.Ctx()
	fields := log.Fields{"traceID": ctx.Value(trace.TIDKey), "func": "UpdateTodo"}
	log.WithFields(fields).Debug("begin")

	// parse id param and attempt string to int conversion
	rawTID := ps.ByName("id")
	todoID, err := strconv.Atoi(rawTID)
	if err != nil {
		log.WithFields(fields).WithField("todoID", rawTID).WithError(err).Error("param conversion error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// attempt to load todo file
	todos, err := todotxt.LoadFromPath("todo.txt")
	if err != nil {
		log.WithFields(fields).WithError(err).Error("failed to load todo.txt")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// attempt to get requested task by ID
	task, err := todos.GetTask(todoID)
	if err != nil { // TODO: investigate potential GetTask failure reasons
		log.WithFields(fields).WithField("todoID", todoID).WithError(err).Error("task not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// attempt to decode request body into task
	if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.WithFields(fields).WithError(err).Error("failed to decode request body into task")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// attempt bulk write all todos to file
	if err = todotxt.WriteToPath(&todos, "todo.txt"); err != nil {
		log.WithFields(fields).WithError(err).Error("failed to write todos to file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set content type and attempt to encode response
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(task); err != nil {
		log.WithFields(fields).WithError(err).Error("failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.WithFields(fields).Debug("success")
}

// CreateTodo will create a new todo in todo.txt.
// TODO: verify the newly generated ID is unique.
func CreateTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// setup traced logging and context
	ctx := trace.Ctx()
	fields := log.Fields{"traceID": ctx.Value(trace.TIDKey), "func": "CreateTodo"}
	log.WithFields(fields).Debug("begin")

	// attempt to load todo file
	todos, err := todotxt.LoadFromPath("todo.txt")
	if err != nil {
		log.WithFields(fields).WithError(err).Error("failed to load todo.txt")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// attempt to decode request body into new task
	task := &todotxt.Task{}
	if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.WithFields(fields).WithError(err).Error("failed to decode request body into task")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// set created date to now
	task.CreatedDate = time.Now()

	// add new task to todos
	todos.AddTask(task)

	// attempt to bulk write todos to file
	if err = todotxt.WriteToPath(&todos, "todo.txt"); err != nil {
		log.WithFields(fields).WithError(err).Error("failed to write todos to file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set content type and attempt to encode response
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(task); err != nil {
		log.WithFields(fields).WithError(err).Error("failed to encode task")
	}

	log.WithFields(fields).WithField("todoID", task.ID).Debug("success")
}

// DeleteTodo will delete a todo from the task list
func DeleteTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// setup traced logging and context
	ctx := trace.Ctx()
	fields := log.Fields{"traceID": ctx.Value(trace.TIDKey), "func": "DeleteTodo"}
	log.WithFields(fields).Debug("begin")

	// parse id param and attempt string to int conversion
	rawTID := ps.ByName("id")
	todoID, err := strconv.Atoi(rawTID)
	if err != nil {
		log.WithFields(fields).WithField("todoID", rawTID).WithError(err).Error("param conversion error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// attempt to load todo file
	todos, err := todotxt.LoadFromPath("todo.txt")
	if err != nil {
		log.WithFields(fields).WithError(err).Error("failed to load todo.txt")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// attempt to find the task
	todo, err := todos.GetTask(todoID)
	if err != nil {
		log.WithFields(fields).WithField("id", todoID).WithError(err).Error("failed to get task")
		w.WriteHeader(http.StatusNotFound)
	}

	// attempt to remove the task
	if err := todos.RemoveTask(*todo); err != nil {
		log.WithFields(fields).WithField("id", todo.ID).WithError(err).Error("failed to remove task")
		w.WriteHeader(http.StatusNotFound)
	}

	// attempt to bulk write todos to file
	if err = todotxt.WriteToPath(&todos, "todo.txt"); err != nil {
		log.WithFields(fields).WithError(err).Error("failed to write todos to file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.WithFields(fields).WithField("id", todoID).Debug("success")
}

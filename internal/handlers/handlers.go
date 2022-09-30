package handlers

import (
	"context"
	"encoding/json"
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
func ListTodos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	todos = filterTodos(ctx, todos, p)

	// set content type and attempt to encode response
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(todos); err != nil {
		log.WithFields(fields).WithError(err).Error("failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write success code
	w.WriteHeader(http.StatusOK)
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

	// write success code
	w.WriteHeader(http.StatusOK)
	log.WithFields(fields).Debug("success")
}

// UpdateTodo takes the body of the request and updates the todo in todo.txt
func UpdateTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// setup traced logging and context
	ctx := trace.Ctx()
	fields := log.Fields{"traceID": ctx.Value(trace.TIDKey), "func": "UpdateTodo"}
	log.WithFields(fields).Debug("begin")

	// parse id param and attempt string to int conversion
	// TODO: common helper
	rawTID := ps.ByName("id")
	todoID, err := strconv.Atoi(rawTID)
	if err != nil {
		log.WithFields(fields).WithField("todoID", rawTID).WithError(err).Error("param conversion error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// attempt to load todo file
	// TODO: common helper
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
	// TODO: investigate performance issues with bulk write
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

	// write success header
	w.WriteHeader(http.StatusOK)
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

	// TODO: handle error
	w.WriteHeader(http.StatusOK)
	log.WithFields(fields).WithField("todoID", task.ID).Debug("success")
}

func filterTodos(ctx context.Context, todos todotxt.TaskList, params httprouter.Params) todotxt.TaskList {
	var preds []todotxt.Predicate
	for _, param := range params {
		switch param.Key {
		case "context":
			preds = append(preds, todotxt.FilterByContext(param.Value))
		case "priority":
			preds = append(preds, todotxt.FilterByPriority(param.Value))
		case "project":
			preds = append(preds, todotxt.FilterByProject(param.Value))
		case "after":
			preds = append(preds, filterByAfter(ctx, param.Value))
		case "before":
			preds = append(preds, filterByBefore(ctx, param.Value))
		}
	}

	if len(preds) > 0 {
		if len(preds) == 1 {
			todos = todos.Filter(preds[0])
		} else {
			todos = todos.Filter(preds[0], preds[1:]...)
		}
	}

	return todos
}

func filterByAfter(ctx context.Context, dateStr string) todotxt.Predicate {
	return func(task todotxt.Task) bool {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.WithFields(log.Fields{"traceID": ctx.Value(trace.TIDKey), "func": "filterByAfter", "date": dateStr}).
				WithError(err).Error("failed to parse date query")
			return false
		}

		return task.DueDate.After(date)
	}
}

func filterByBefore(ctx context.Context, dateStr string) todotxt.Predicate {
	return func(task todotxt.Task) bool {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.WithFields(log.Fields{"traceID": ctx.Value(trace.TIDKey), "func": "filterByBefore", "date": dateStr}).
				WithError(err).Error("failed to parse date query")
			return false
		}

		return task.DueDate.Before(date)
	}
}

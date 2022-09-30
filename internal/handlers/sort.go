package handlers

import (
	"context"
	"errors"
	"github.com/1set/todotxt"
	"github.com/chrishaynes21/apichallenge/pkg/trace"
	log "github.com/sirupsen/logrus"
	"net/url"
)

var (
	UnknownSortError = errors.New("unknown sort type")
)

var sortTable = map[string]todotxt.TaskSortByType{
	"SortTaskIDAsc":       todotxt.SortTaskIDAsc,
	"SortTaskIDDesc":      todotxt.SortTaskIDDesc,
	"SortTodoTextAsc":     todotxt.SortTodoTextAsc,
	"SortPriorityAsc":     todotxt.SortPriorityAsc,
	"SortPriorityDesc":    todotxt.SortPriorityDesc,
	"SortCreatedDateAsc":  todotxt.SortCreatedDateAsc,
	"SortCreatedDateDesc": todotxt.SortCreatedDateDesc,
	"SortDueDateAsc":      todotxt.SortDueDateAsc,
	"SortDueDateDesc":     todotxt.SortDueDateDesc,
	"SortContextAsc":      todotxt.SortContextAsc,
	"SortContextDesc":     todotxt.SortContextDesc,
	"SortProjectAsc":      todotxt.SortProjectAsc,
	"SortProjectDesc":     todotxt.SortProjectDesc,
}

func sortTodos(ctx context.Context, todos todotxt.TaskList, queries url.Values) (todotxt.TaskList, error) {
	if !queries.Has("order") {
		return todos, nil
	}

	for _, sortType := range queries["order"] {
		sort, found := sortTable[sortType]
		if !found {
			log.WithFields(log.Fields{"traceID": ctx.Value(trace.TIDKey), "func": "sortTodos", "sortType": sortType}).
				WithError(UnknownSortError).Warn("unknown sort type encountered")
			return nil, UnknownSortError
		}

		if err := todos.Sort(sort); err != nil {
			log.WithFields(log.Fields{"traceID": ctx.Value(trace.TIDKey), "func": "sortTodos", "sortType": sortType}).
				WithError(err).Error("failed to sort todos")
			return nil, err
		}
	}

	return todos, nil
}

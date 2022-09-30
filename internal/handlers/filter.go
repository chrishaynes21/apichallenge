package handlers

import (
	"context"
	"github.com/1set/todotxt"
	"github.com/chrishaynes21/apichallenge/pkg/trace"
	log "github.com/sirupsen/logrus"
	"net/url"
	"time"
)

func filterTodos(ctx context.Context, todos todotxt.TaskList, queries url.Values) todotxt.TaskList {
	var preds []todotxt.Predicate
	output := todos
	for key, query := range queries {
		for _, value := range query {
			switch key {
			case "after":
				preds = append(preds, filterByAfter(ctx, value))
			case "before":
				preds = append(preds, filterByBefore(ctx, value))
			case "context":
				preds = append(preds, todotxt.FilterByContext(value))
			case "priority":
				preds = append(preds, todotxt.FilterByPriority(value))
			case "project":
				preds = append(preds, todotxt.FilterByProject(value))
			}
		}
	}

	for _, pred := range preds {
		output = output.Filter(pred)
	}

	return output
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

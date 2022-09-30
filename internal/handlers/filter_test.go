package handlers

import (
	"github.com/1set/todotxt"
	"github.com/chrishaynes21/apichallenge/pkg/trace"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func Test_filterTodos(t *testing.T) {
	type args struct {
		todos   todotxt.TaskList
		queries url.Values
	}
	tests := []struct {
		name string
		args args
		want todotxt.TaskList
	}{
		{
			name: "should filter by context",
			args: args{
				todos: testTodos,
				queries: url.Values{
					"context": []string{"test"},
				},
			},
			want: []todotxt.Task{testTodos[0], testTodos[1]},
		},
		{
			name: "should double filter by context",
			args: args{
				todos: testTodos,
				queries: url.Values{
					"context": []string{"test", "another test"},
				},
			},
			want: []todotxt.Task{testTodos[0], testTodos[1], testTodos[2]},
		},
		{
			name: "should filter by priority",
			args: args{
				todos: testTodos,
				queries: url.Values{
					"priority": []string{"A"},
				},
			},
			want: []todotxt.Task{testTodos[0]},
		},
		{
			name: "should double filter by priority",
			args: args{
				todos: testTodos,
				queries: url.Values{
					"priority": []string{"A", "B"},
				},
			},
			want: []todotxt.Task{testTodos[0], testTodos[1]},
		},
		{
			name: "should filter by project",
			args: args{
				todos: testTodos,
				queries: url.Values{
					"project": []string{"B proj"},
				},
			},
			want: []todotxt.Task{testTodos[0], testTodos[1]},
		},
		{
			name: "should double filter by project",
			args: args{
				todos: testTodos,
				queries: url.Values{
					"project": []string{"B proj", "C proj"},
				},
			},
			want: []todotxt.Task{testTodos[0], testTodos[1], testTodos[2]},
		},
		{
			name: "should filter by before date",
			args: args{
				todos: testTodos,
				queries: url.Values{
					"before": []string{"2023-01-01"},
				},
			},
			want: []todotxt.Task{testTodos[1], testTodos[2]},
		},
		{
			name: "should filter by after date",
			args: args{
				todos: testTodos,
				queries: url.Values{
					"after": []string{"2023-01-01"},
				},
			},
			want: []todotxt.Task{testTodos[0]},
		},
		{
			name: "should exclusive filter by after and before date",
			args: args{
				todos: testTodos,
				queries: url.Values{
					"after":  []string{"2022-01-01"},
					"before": []string{"2022-02-01"},
				},
			},
			want: []todotxt.Task{testTodos[1]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterTodos(trace.Ctx(), tt.args.todos, tt.args.queries)
			assert.ElementsMatchf(t, tt.want, got, "task mismatch")
		})
	}
}

package handlers

import (
	"github.com/1set/todotxt"
	"github.com/chrishaynes21/apichallenge/pkg/trace"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testTodos = []todotxt.Task{
	{
		ID:       1,
		Original: "completed todo 1",
		Todo:     "todo 1",
		Priority: "A",
		Projects: []string{"A proj", "B proj"},
		Contexts: []string{"test"},
		AdditionalTags: map[string]string{
			"tag1": "meta1",
			"tag2": "meta2",
		},
		CreatedDate:   time.Date(2020, 1, 2, 3, 0, 0, 0, time.UTC),
		DueDate:       time.Date(2023, 1, 2, 3, 0, 0, 0, time.UTC),
		CompletedDate: time.Date(2022, 9, 19, 0, 0, 0, 0, time.UTC),
		Completed:     true,
	},
	{
		ID:            2,
		Original:      "unfinished todo 2",
		Todo:          "todo 2",
		Priority:      "B",
		Projects:      []string{"B proj"},
		Contexts:      []string{"test"},
		CreatedDate:   time.Date(2020, 2, 2, 3, 0, 0, 0, time.UTC),
		DueDate:       time.Date(2022, 1, 2, 3, 0, 0, 0, time.UTC),
		CompletedDate: time.Time{},
		Completed:     false,
	},
	{
		ID:            3,
		Original:      "previously finished todo 3",
		Todo:          "todo 3",
		Priority:      "C",
		Projects:      []string{"C proj"},
		Contexts:      []string{"another test"},
		CreatedDate:   time.Date(2020, 3, 2, 3, 0, 0, 0, time.UTC),
		DueDate:       time.Date(2022, 1, 2, 3, 0, 0, 0, time.UTC),
		CompletedDate: time.Date(2022, 3, 2, 3, 0, 0, 0, time.UTC),
		Completed:     true,
	},
}

func Test_filterTodos(t *testing.T) {
	type args struct {
		todos  todotxt.TaskList
		params httprouter.Params
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
				params: []httprouter.Param{
					{Key: "context", Value: "test"},
				},
			},
			want: []todotxt.Task{testTodos[0], testTodos[1]},
		},
		{
			name: "should double filter by context",
			args: args{
				todos: testTodos,
				params: []httprouter.Param{
					{Key: "context", Value: "test"},
					{Key: "context", Value: "another test"},
				},
			},
			want: []todotxt.Task{testTodos[0], testTodos[1], testTodos[2]},
		},
		{
			name: "should filter by priority",
			args: args{
				todos: testTodos,
				params: []httprouter.Param{
					{Key: "priority", Value: "A"},
				},
			},
			want: []todotxt.Task{testTodos[0]},
		},
		{
			name: "should double filter by priority",
			args: args{
				todos: testTodos,
				params: []httprouter.Param{
					{Key: "priority", Value: "A"},
					{Key: "priority", Value: "B"},
				},
			},
			want: []todotxt.Task{testTodos[0], testTodos[1]},
		},
		{
			name: "should filter by project",
			args: args{
				todos: testTodos,
				params: []httprouter.Param{
					{Key: "project", Value: "B proj"},
				},
			},
			want: []todotxt.Task{testTodos[0], testTodos[1]},
		},
		{
			name: "should double filter by project",
			args: args{
				todos: testTodos,
				params: []httprouter.Param{
					{Key: "project", Value: "B proj"},
					{Key: "project", Value: "C proj"},
				},
			},
			want: []todotxt.Task{testTodos[0], testTodos[1], testTodos[2]},
		},
		{
			name: "should filter by before date",
			args: args{
				todos: testTodos,
				params: []httprouter.Param{
					{Key: "before", Value: "2023-01-01"},
				},
			},
			want: []todotxt.Task{testTodos[1], testTodos[2]},
		},
		{
			name: "should filter by after date",
			args: args{
				todos: testTodos,
				params: []httprouter.Param{
					{Key: "after", Value: "2023-01-01"},
				},
			},
			want: []todotxt.Task{testTodos[0]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterTodos(trace.Ctx(), tt.args.todos, tt.args.params)
			assert.ElementsMatchf(t, tt.want, got, "task mismatch")
		})
	}
}

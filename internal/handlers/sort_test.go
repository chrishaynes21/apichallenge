package handlers

import (
	"github.com/1set/todotxt"
	"github.com/chrishaynes21/apichallenge/pkg/trace"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func Test_sortTodos(t *testing.T) {
	type args struct {
		todos   todotxt.TaskList
		queries url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    todotxt.TaskList
		wantErr error
	}{
		{
			name: "should handle no sort type",
			args: args{
				todos:   testTodos,
				queries: map[string][]string{},
			},
			want:    todotxt.TaskList{testTodos[0], testTodos[1], testTodos[2]},
			wantErr: nil,
		},
		{
			name: "should apply a sort type",
			args: args{
				todos: testTodos,
				queries: map[string][]string{
					"order": {"SortTaskIDDesc"},
				},
			},
			want:    todotxt.TaskList{testTodos[2], testTodos[1], testTodos[0]},
			wantErr: nil,
		},
		{
			name: "should apply multiple sort types",
			args: args{
				todos: testTodos,
				queries: map[string][]string{
					"order": {"SortContextAsc", "SortProjectAsc"},
				},
			},
			want:    todotxt.TaskList{testTodos[2], testTodos[0], testTodos[1]},
			wantErr: nil,
		},
		{
			name: "should throw error for illegal sort type",
			args: args{
				todos: testTodos,
				queries: map[string][]string{
					"order": {"SortContextAsc", "Oopsie"},
				},
			},
			want:    nil,
			wantErr: ErrUnknownSort,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sortTodos(trace.Ctx(), tt.args.todos, tt.args.queries)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error(), "expected error mismatch")
				assert.Nil(t, got, "expected task list to be nil after error")
			} else {
				assert.NoError(t, err, "unexpected error")
				assert.Equal(t, got, tt.want, "elements or order does not match")
			}
		})
	}
}

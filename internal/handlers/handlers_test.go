package handlers

import (
	"github.com/1set/todotxt"
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
		DueDate:       time.Date(2022, 3, 3, 3, 0, 0, 0, time.UTC),
		CompletedDate: time.Date(2022, 3, 2, 3, 0, 0, 0, time.UTC),
		Completed:     true,
	},
}

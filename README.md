# TODO API Challenge

## API

### Model

Task

```json
{
  "ID": 1,
  "Original": "(A) 2022-09-30 todo 1 new text @new @test +apichallenge tag1:meta1 tag2:metachange due:2023-01-02",
  "Todo": "todo 1 new text",
  "Priority": "A",
  "Projects": [
    "apichallenge"
  ],
  "Contexts": [
    "new",
    "test"
  ],
  "AdditionalTags": {
    "tag1": "meta1",
    "tag2": "meta2"
  },
  "CreatedDate": "2022-09-30T00:00:00-07:00",
  "DueDate": "2023-01-02T00:00:00-08:00",
  "CompletedDate": "0001-01-01T00:00:00Z",
  "Completed": false
}
```

### Paths

#### `GET /todos`

Lists all tasks. Has optional filtering and sorting parameters:

- Parameters
    - `after`:  input in `YYYY-MM-DD` format (e.g. "2021-01-01") should filter out any tasks that have due dates after
      the specified time.
    - `before`: input in `YYYY-MM-DD` format (e.g. "2021-01-01") should filter out any tasks that have due dates before
      the specified time.
    - `context`: output should only include tasks with that context. If more than one "context tag" param is included,
      the filter should accept any todo that has one or more of the given context tags and reject any todo lacking the
      specified context tags.
    - `priority`: output should include only tasks that have the associated priority. If more than one priority tag is
      included in the params, the output should include all and only tasks that have *any* of the specified priority
      tags.
    - `project`: output should only include tasks associated with that project. If more than one "project tag" param is
      included, the filter should accept any todo that has one or more of the given project tags and reject any todo
      lacking the specified project tags.
    - `order`:  follows all orders defined by the `TaskSortByType` type in the codebase (
      defined [here](https://pkg.go.dev/github.com/1set/todotxt#TaskSortByType) ).

#### `GET /todos/<id>`

Gets the task with the `id`.

#### `PUT /todos/<id>`

Update the task at `id` with the payload.

#### `DELETE /todos/<id>`

Remove the task at `id` from the task list.

#### `POST /todos`

Create a task with the payload.

## Other Requirements

1. By default - i.e. in the absence of any other filters - output should be sorted primarily by priority, secondarily by
   due date, and tertiarily by created date. Completed tasks should *always* be listed separately (when they are listed
   at all)

## Setup - Running apichallenge.go

The basic framework has been done for you. The file `apichallenge.go` implements a basic server that can handle the web
connection for you (you don't need to implement authentication or user management for this task). To run this setup you
will need to:

1. Install `todotxt`: `go get github.com/1set/todotxt`
2. Install `httprouter`: `go get github.com/julienschmidt/httprouter`
3. Compile: `go build ./cmd`
4. Run it! : `./cmd.exe` (or other OS runnable)

Now if you navigate to `http://localhost:8080/mainpage.html` you should see the main page.

## Integration

You will see that the basic REST routes have been implemented for you. All you need to do is modify the relevant
functions to satisfy the tasks below.

## Static Assets

`apichallenge.go` will serve up any file that is stored in the `static/` folder as a direct path. For example,
navigating to `http://localhost:8080/mainpage.html` serves the `mainpage.html` file. Navigating
to `http://localhost:8080/stylesheets/main.css` serves the main stylesheet. ETC. You can put any html or css or
javascript libraries you need for the interface in the `static` folder.

## todo.txt

The application will load tasks from and save tasks into the provided `todo.txt` file. You do not need to integrate with
any kind of database. Just using this file for storage is OK.

## Tasks

Be sure that you are familiar with the [todo.txt file format](http://todotxt.org/) before starting this task. Don't
worry, it's easy to learn!

Your web application should support the following actions:

1. List all the tasks. This part of the project should accept query parameters to filter the list of todos.  **The
   filtering must be done on the backend.**  The point of this part of the task is to demonstrate that you can use an
   existing Go codebase. So, you **must** implement the list filters on the backend in Go (and not on the frontend in
   Javascript)!!!  The following query parameters should be accepted:
    1. `projects` - Any projects included here should filter the output to include *only* those tasks that are
       associated with one or more of the given projects.
    2. `priority` - Any priorities included here should filter the output to include *only* those tasks that have one of
       the priorities in question
    3. `context` - Any contexts included here should filter the output to include *only* those tasks that have one of
       the contexts in question
    4. `order` - If the order param is set, the tasks should come back in the order specified. You should support all
       the orders given in the [TaskSortByType](https://pkg.go.dev/github.com/1set/todotxt#TaskSortByType)) struct.
    5. `duebefore` - this should accept a string representing a datetime and *only* return tasks that (a) have a
       duedate (b) which is before the date specified
    6. `dueafter` - this should accept a string representing a datetime and *only* return tasks that (a) have a dueate (
       b) which is after the date specified
2. Accept input from a user and add a new todo to the list
3. Update any aspect of a task. This includes:
    1. Adding (or removing) a project
    2. Adding (or removing) a context
    3. Setting (or changing) the priority
    4. Setting (or changing) the duedate
4. Mark a task as complete
5. (Optional) Delete a task
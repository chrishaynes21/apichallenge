package cmd

import (
	"github.com/chrishaynes21/apichallenge/internal/handlers"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {

	router := httprouter.New()
	router.GET("/todos", handlers.ListTodos)
	router.GET("/todo/:id", handlers.GetTodo)
	router.PUT("/todo/:id", handlers.UpdateTodo)
	router.POST("/todo", handlers.CreateTodo)

	router.NotFound = http.FileServer(http.Dir("./static"))

	log.Fatal(http.ListenAndServe(":8080", router))
}

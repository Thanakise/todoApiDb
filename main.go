package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/thanakize/todoApiDB/database"
)
type Todo struct {
	ID    int `json:"id"`
	Title string `json:"title"`
	Status  string `json:"status"`
}

func main() {
  	ctx, cancle := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	db := database.ConnectDB()
  	r := gin.Default()

	r.GET("/api/v1/todos", getTodos(db))
	r.GET("/api/v1/todos/:id", getTodoById(db))
	r.POST("/api/v1/todos", insertTodo(db))
	r.PUT("/api/v1/todos/:id", updateTodo(db))
	r.DELETE("/api/v1/todos/:id", deleteTodos(db))
	r.PATCH("/api/v1/todos/:id/actions/status",updateStatus(db))
	r.PATCH("/api/v1/todos/:id/actions/title",updateTitle(db))


	defer cancle()
	srv := http.Server{
		Addr: ":" + os.Getenv("PORT"),
		Handler:  r,
	}
	go func ()  {
		<-ctx.Done()
		fmt.Println("Shutting down...")
		ctx, cancle := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancle()

		if err := srv.Shutdown(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed){
			log.Println(err)
		}
		}

	}()

	if err := srv.ListenAndServe(); err != nil{
		log.Println(err)
	}
}


func getTodos(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var todos []Todo
	rows, err := db.Query("SELECT id, title, status FROM todos")
	if err != nil {
		log.Fatal("can't query all todos", err)
	}
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Status)
		if err != nil {
			log.Fatal("can't Scan row into variable", err)
		}
		todos = append(todos, todo)
		fmt.Println(todo.ID, todo.Title, todo.Status)
	}
		ctx.JSON(http.StatusOK,todos)
	}
}
func getTodoById(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
	id := ctx.Param("id")
	var todo Todo
	row := db.QueryRow("SELECT id, title, status FROM todos where id=$1",id)
	err := row.Scan(&todo.ID, &todo.Title, &todo.Status)
	if err != nil {
		ctx.JSON(400, gin.H{
			"err": err.Error(),
			"statusCode": 400,
		})
		return
	}
	fmt.Println(todo.ID, todo.ID, todo.Status)
	ctx.JSON(http.StatusOK, todo)
	}
	
}

func deleteTodos(db *sql.DB) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		_, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
		if err != nil {
			ctx.JSON(409, gin.H{
				"err": err.Error(),
				"statusCode": 409,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	 }
	
	}

func insertTodo(db *sql.DB) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var todo Todo
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		q := "INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id"
		row := db.QueryRow(q, todo.Title, todo.Status)
		var id int
		err := row.Scan(&id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
				"statusCode": http.StatusInternalServerError,
			})
			return
		}
		ctx.JSON(http.StatusAccepted, gin.H{
			"id" : id,
			"title" : todo.Title,
			"status" : todo.Status,
		})
	}
}

func updateTodo(db *sql.DB) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		var todo Todo
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		q := "update todos set title = $1, status = $2 where id = $3 RETURNING id, title, status"
		// q := "INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id"
		row := db.QueryRow(q, todo.Title, todo.Status, id)
		var newTodo Todo 
		err := row.Scan(&newTodo.ID, &newTodo.Title, &newTodo.Status)
		fmt.Println(newTodo)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
				"statusCode": http.StatusInternalServerError,
			})
			return
		}
		ctx.JSON(http.StatusAccepted, gin.H{
			"id" : newTodo.ID,
			"title" : newTodo.Title,
			"status" : newTodo.Status,
		})
	}
}

func updateStatus(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		var todo Todo
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		q := "update todos set status = $1 where id = $2 RETURNING id, title, status"
		row := db.QueryRow(q, todo.Status, id) 
		var newTodo Todo 
		err := row.Scan(&newTodo.ID, &newTodo.Title, &newTodo.Status)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
				"statusCode": http.StatusInternalServerError,
			})
			return
		}
		ctx.JSON(http.StatusAccepted, gin.H{
			"id" : newTodo.ID,
			"title" : newTodo.Title,
			"status" : newTodo.Status,
		})
	}
}


func updateTitle(db *sql.DB) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		var todo Todo
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		q := "update todos set title = $1 where id = $2 RETURNING id, title, status"
		// q := "INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id"
		row := db.QueryRow(q, todo.Title, id)
		var newTodo Todo 
		err := row.Scan(&newTodo.ID, &newTodo.Title, &newTodo.Status)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
				"statusCode": http.StatusInternalServerError,
			})
			return
		}
		ctx.JSON(http.StatusAccepted, gin.H{
			"id" : newTodo.ID,
			"title" : newTodo.Title,
			"status" : newTodo.Status,
		})
	}
}




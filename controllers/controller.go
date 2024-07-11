package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanakize/todoApiDB/sharedinterface"
	"github.com/thanakize/todoApiDB/views"
)

func GetController(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) { 
		todos, err := views.GetTodos(db)
		if err != nil {
			ctx.JSON(400, gin.H{
				"err": err.Error(),
				"statusCode": 400,
			})
			return
		}
		ctx.JSON(200, todos)
	}
}

func GetIdController(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) { 
		id := ctx.Param("id")
		if id == ""{
			ctx.JSON(400, gin.H{
				"err": "id is required",
				"statusCode": 400,
			})
			
			return
		}
		todos, err := views.GetTodoById(db, id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"err": err.Error(),
				"statusCode": 400,
			})
			return
		}
		ctx.JSON(http.StatusAccepted, gin.H{
			"id" : todos.ID,
			"title" : todos.Title,
			"status" : todos.Status,
		})
	}
}
func PostController(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) { 
		var todo sharedinterface.Todo
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newTodo, err := views.InsertTodo(db, todo)
		if err != nil {
			ctx.JSON(400, gin.H{
				"err": err.Error(),
				"statusCode": 400,
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
func PutController(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) { 
		id := ctx.Param("id")
		var todo sharedinterface.Todo
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newTodo, err := views.UpdateTodo(db, todo, id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"err": err.Error(),
				"statusCode": 400,
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
func DeleteController(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) { 
		id := ctx.Param("id")
		err := views.DeleteTodo(db, id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"err": err.Error(),
				"statusCode": 400,
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status" : "success",
		})
	}
}
func PatchStatusController(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) { 
		id := ctx.Param("id")
		var todo sharedinterface.Todo
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newTodo, err := views.PatchStatusTodo(db, todo, id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"err": err.Error(),
				"statusCode": 400,
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
func PatchTitleController(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) { 
		id := ctx.Param("id")
		var todo sharedinterface.Todo
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newTodo, err := views.PatchTitleTodo(db, todo, id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"err": err.Error(),
				"statusCode": 400,
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
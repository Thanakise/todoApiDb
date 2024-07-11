package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/thanakize/todoApiDB/controllers"
)

func InitRoute(r *gin.Engine, db *sql.DB) *gin.RouterGroup{

	router := r.Group("/api/v1/todos")
	
	router.GET("/", controllers.GetController(db))
	router.GET("/:id", controllers.GetIdController(db))
	router.POST("/", controllers.PostController(db))
	router.PUT("/:id", controllers.PutController(db))
	router.DELETE("/:id", controllers.DeleteController(db))
	router.PATCH("/:id/actions/status",controllers.PatchStatusController(db))
	router.PATCH("/:id/actions/title",controllers.PatchTitleController(db))
	return router
}
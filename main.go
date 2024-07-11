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

	"github.com/gin-gonic/gin"
	"github.com/thanakize/todoApiDB/database"
	"github.com/thanakize/todoApiDB/router"
)


func main() {
  	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	
	db := database.ConnectDB()
	defer db.Close()

  	r := gin.Default()
	router.InitRoute(r, db)

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

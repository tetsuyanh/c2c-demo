package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tetsuyanh/c2c-demo/cmd/api/v1"
)

const (
	exitOK int = iota
	exitError
)

var (
	version string
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	g := gin.Default()

	// for healthe-check
	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, version)
		return
	})
	// endpoints
	v1.Router(g)

	// run server
	srv := &http.Server{
		Addr:    ":8000",
		Handler: g,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("finish to listen: %s\n", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %s\n", err)
		return exitError
	}

	return exitOK
}

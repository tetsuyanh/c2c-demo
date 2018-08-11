package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/tetsuyanh/c2c-demo/cmd/api/v1"
	"github.com/tetsuyanh/c2c-demo/conf"
	"github.com/tetsuyanh/c2c-demo/repo"
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
	c, eConf := loadConfig()
	if eConf != nil {
		log.Printf("load conf error: %s\n", eConf)
		return exitError
	}

	if eRepo := repo.Setup(c.Postgres); eRepo != nil {
		log.Printf("repo setup error: %s\n", eRepo)
		return exitError
	}
	defer repo.TearDown()

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

func loadConfig() (*conf.Config, error) {
	path := os.Getenv("C2C_DEMO_CONF_PATH")
	name := os.Getenv("C2C_DEMO_CONF_NAME")
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	if errRead := viper.ReadInConfig(); errRead != nil {
		fmt.Println("errRead:", errRead)
		return nil, errRead
	}
	c := &conf.Config{}
	if errUnmarshal := viper.Unmarshal(c); errUnmarshal != nil {
		fmt.Println("errUnmarshal:", errUnmarshal)
		return nil, errUnmarshal
	}
	return c, nil
}

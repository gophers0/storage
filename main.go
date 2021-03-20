package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gophers0/users/internal/model"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gophers0/users/internal/config"
	"github.com/gophers0/users/internal/repository/postgres"
	"github.com/gophers0/users/internal/service/httpsrv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/zergu1ar/Gaarx"
	"github.com/zergu1ar/logrus-filename"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		return
	}

	flag.Parse()
	var stop = make(chan os.Signal)
	var done = make(chan bool, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Read config && Init Logger
	cfg, err := config.Init(os.Getenv("ENV"))
	if err != nil {
		panic(err)
	}

	ctx, finish := context.WithCancel(context.Background())
	var application = &gaarx.App{}
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	if err := application.LoadConfig(cfg, "json", &config.Config{}); err != nil {
		panic(err)
	}
	configs := application.Config().(*config.Config)
	if err := application.InitializeLogger(configs.GetLogWay(), configs.GetLogDestination(), &logrus.TextFormatter{DisableColors: true}); err != nil {
		panic(err)
	}
	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	application.GetLog().AddHook(filenameHook)
	application.GetLog().SetLevel(logrus.Level(configs.System.Log.Level))
	application.Initialize(
		postgres.WithDatabase(configs.GetConnString(), configs.System.DB.Dialect, model.ModelsList...),
		gaarx.WithContext(ctx),
		gaarx.WithServices(
			httpsrv.New(application.GetLog()),
		),
	)
	go func() {
		sig := <-stop
		time.Sleep(2 * time.Second)
		finish()
		fmt.Printf("caught sig: %+v\n", sig)
		done <- true
	}()
	application.Work()
	<-ctx.Done()
	os.Exit(0)
}

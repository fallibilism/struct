package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"v/pkg/config"
	"v/pkg/handlers"

	"github.com/urfave/cli/v2"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s\n", c.App.Version)
	}

	app := &cli.App{
		Name:        "struct",
		Usage:       "video conference system",
		Description: "xxx",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "Configuration file",
				DefaultText: "config.yaml",
				Value:       "config.yaml",
			},
		},
		Action:  runServer,
		Version: Version,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func runServer(c *cli.Context) error {

	config.SetConfig(c.String("config"))

	router := handlers.Handler()

	println("config.Postgres")
	db, err := config.NewDbConnection(&config.Postgres)
	println("config.Postgres")
	if err != nil {
		log.Panicln("could not connect to database", "error", err)
	}

	// redis, err := config.NewRedisConnection(&config.Redis)

	appConf := &config.AppConfig{
		DB: db,
		// Redis: redis,
	}

	config.TestConfig = appConf // a hack for testing

	if err != nil {
		log.Panicln("could not connect to redis", "error", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigChan
		log.Panicln("exit requested, shutting down", "signal", sig)
		_ = router.Shutdown()
	}()
	return router.Listen()
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"v/pkg/config"
	"v/pkg/handlers"
	"v/pkg/models"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s\n", c.App.Version)
	}

	app := &cli.App{
		Name:        "Komrade",
		Usage:       "Komrade API",
		Description: "xxx",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
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

	conf := config.SetConfig(c.String("config"))

	router := handlers.Handler()

	err := setupConnections(conf)
	// println("ignore setup connection for now")

	if err != nil {
		log.Fatalln(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigChan
		log.Panicln("exit requested, shutting down", "signal", sig)
		_ = router.Shutdown()
		//
	}()

	return router.Listen()
}

// redis and postgres connection setup
func setupConnections(conf *config.Config) error {

	db, err := config.NewDbConnection(&conf.Db)
	if err != nil {
		err := fmt.Errorf("could not connect to database: %v", err)
		return err
	}

	migrations(db)
	// redis, err := NewRedisConnection(&conf.Redis)

	// if err != nil {
	// 	err := fmt.Errorf("could not connect to redis: %v", err)
	// 	return err
	// }

	appConf := &config.AppConfig{
		DB: db,
		// Redis: redis,
	}

	// config.TestConfig = appConf // a hack for testing
	config.App = appConf

	return nil

}

// model migration here
func migrations(db *gorm.DB) {

	db.AutoMigrate(&models.Room{})

}

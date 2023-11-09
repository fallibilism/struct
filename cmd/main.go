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

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func main() {
	var configFile string

	if err := godotenv.Load(); err != nil {
		panic("config: " + err.Error())
	}

	switch os.Getenv("ENV") {
		case "production":
			configFile = "config.prod.yaml"
		case "testing":
			configFile = "config.test.yaml"
		case "development":
			fallthrough
		default:
			configFile = "config.dev.yaml"
	}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s\n", c.App.Version)
	}

	app := &cli.App{
		Name:        "Komrade",
		Usage:       "Komrade API",
		Description: "xxx",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "migrations",
				Aliases: []string{"m"},
				Usage:       "Start with a Database migration",
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases: []string{"c"},
				Usage:       "Configuration file",
				DefaultText: configFile,
				Value:       configFile,
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

	setupConnections(conf, c.Bool("migrations"))
	// println("ignore setup connection for now")


	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		defer handlePanic(router)
		sig := <-sigChan
		log.Panicln("exit requested, shutting down", "signal", sig)
	}()
	return router.Listen()
}

// redis and postgres connection setup
func setupConnections(conf *config.Config, withMigration bool) {

	db, err := config.NewDbConnection(&conf.Postgres)
	if err != nil {
		err := fmt.Errorf("could not connect to database: %v", err)
		panic(err)
	}

	if withMigration {
		println("Running migrations")
		migrations(db)
	}

	appConf := &config.AppConfig{
		DB: db,
	}

	config.App = appConf

}

// model migration here
func migrations(db *gorm.DB) {
	if err := db.AutoMigrate(&models.Room{}); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.Token{}); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic(err)
	}
}

func handlePanic(router *handlers.Routes) {
    if r := recover(); r != nil {
        // Handle the panic, log it, or take other appropriate actions
        fmt.Println("Panic:", r)
	_ = router.Shutdown()
	os.Exit(1)
    }
}

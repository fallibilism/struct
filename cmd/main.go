package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"v/pkg/config"
	"v/pkg/handlers"
	"v/pkg/models"

	stt "cloud.google.com/go/speech/apiv1"
	tts "cloud.google.com/go/texttospeech/apiv1"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/option"
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
			&cli.StringFlag{
				Name:     "gcp-credentials-path",
				Aliases: []string{"gcp-cred-path"},
				Usage:       "Credential file",
				DefaultText: os.Getenv("GCP_CREDENTIALS_PATH"),
			},
			&cli.StringFlag{
				Name:     "gcp-credentials-body",
				Aliases: []string{"gcp-cred-body"},
				Usage:       "Credential body",
			},
		},
		Action:  runServer,
		Version: Version,
	}

	err := app.Run(os.Args)
	if err != nil {
		println(err)
		os.Exit(1)
	}
}

func runServer(c *cli.Context) error {
	var router *handlers.Routes
	defer handlePanic(router)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		defer handlePanic(router)
		sig := <-sigChan
		panic("exit requested, shutting down signal: " + sig.String())
	}()
	gcpFile := c.String("gcp-credentials-path")
	if gcpFile == "" {
		gcpFile = os.Getenv("GCP_CREDENTIALS_PATH")
	}


	if err := setupConnections(config.SetConfig(c.String("config")), c.Bool("migrations"), gcpFile, c.String("gcp-credentials-body")); err != nil {
		panic(err)
	}


	router = handlers.Handler()

	// println("ignore setup connection for now")



	return router.Listen()
}

// redis and postgres connection setup
func setupConnections(conf *config.Config, withMigration bool, gcpFile, gcpBody string) error {
	var OpenaiClient *openai.Client
	var TextClient *tts.Client
	var SpeechClient *stt.Client
	var Redis *redis.Client

	ctx := context.Background()

	DB, err := config.NewDbConnection(&conf.Postgres)
	if err != nil {
		err := fmt.Errorf("could not connect to database: %v", err)
		return err
	}

	OpenaiClient = openai.NewClient(conf.Openai.Token)

	gcpCred := option.WithCredentialsFile(gcpFile)

	if gcpBody != "" {
		gcpCred = option.WithCredentialsJSON([]byte(gcpBody))
	}

	SpeechClient, err = stt.NewClient(ctx, gcpCred)
	if err != nil {
		return err
	}

	TextClient, err = tts.NewClient(ctx, gcpCred)
	if err != nil {
		return err
	}


	// --- Validation ---
	if DB == nil {
		return errors.New("db is null")
	}
	if Redis == nil {
		return errors.New("redis is null")
	}
	if TextClient == nil {
		return errors.New("Text client is null")
	}
	if SpeechClient == nil {
		return errors.New("Speech client is null")
	}
	if OpenaiClient == nil {
		return errors.New("openai client is null")
	}

	appConf := &config.AppConfig{
		DB: DB,
		Redis: Redis,
		SpeechClient: SpeechClient,
		TextClient: TextClient,
		OpenaiClient: OpenaiClient,
	}


	if withMigration {
		println("Running migrations")
		if err := migrations(DB); err != nil {
			return err
		}
	}

	config.App = appConf
	return nil
}

// model migration here
func migrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Room{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Token{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	return nil
}

func handlePanic(router *handlers.Routes) {
    if r := recover(); r != nil {
        // Handle the panic, log it, or take other appropriate actions
        fmt.Println("Panic:", r)
	if router == nil {
		os.Exit(1)
	}
	_ = router.Shutdown()
	os.Exit(1)
    }
}

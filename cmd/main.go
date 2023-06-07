package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"v/pkg/handlers"

	"github.com/urfave/cli/v2"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s\n", c.App.Version)
	}

	app := &cli.App{
		Name:        "plugnmeet-server",
		Usage:       "Scalable, Open source web conference system",
		Description: "without option will start server",
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

	router := handlers.Handler()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigChan
		log.Panicln("exit requested, shutting down", "signal", sig)
		_ = router.Shutdown()
	}()
	return router.Listen()
}

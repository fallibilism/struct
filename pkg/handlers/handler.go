package handlers

import (
	"log"
	"v/pkg/config"
	"v/pkg/controllers"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/template/html"
)

type Routes struct {
	App *fiber.App
}

func Handler() *Routes {
	templateEngine := html.New(config.ViewPath, config.ViewExt)
	conf := fiber.Config{
		Views:       templateEngine,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	}
	app := fiber.New(conf)

	router := &Routes{
		App: app,
	}
	router.middleware()

	router.routes()

	return router
}

func (r *Routes) Listen() error {
	log.Printf("Server listening on port %s", config.Server.Port)
	return r.App.Listen(config.Server.Host + ":" + config.Server.Port)
}
func (r *Routes) Shutdown() error {

	log.Printf("Gracefully shutting down server...")

	// close database connection
	// close redis connection
	r.App.Shutdown()

	log.Printf("Server shutdown successful")

	return nil
}

func (r *Routes) middleware() {

	if config.Developement {
		r.App.Use(logger.New())
	}

	r.App.Use(recover.New())
	r.App.Use(cors.New(cors.Config{
		AllowMethods: "POST,GET,OPTIONS",
	}))

	// prometheus
	if config.Prometheus.Enabled {
		prometheus := fiberprometheus.New(config.Prometheus.Namespace)
		prometheus.RegisterAt(r.App, config.Prometheus.MetricsPath)
		r.App.Use(prometheus.Middleware)
	}
}

func (r *Routes) routes() {
	app := r.App

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	app.Get("/info", func(c *fiber.Ctx) error {
		return c.Render("info", nil)
	})

	app.Get("/health-check", controllers.HandleHealthCheck)
	app.Post("/webhook", controllers.HandleWebhook)
	app.Get("/download/uploadedFile/:sid/*", controllers.HandleDownloadUploadedFile)
	app.Get("/download/recording/:token", controllers.HandleDownloadRecording)

	//livekit
	lti_v1 := app.Group("/lti/v1", controllers.HandleV1HeaderToken)
	lti_v1.Post("/room/join", controllers.HandleLTIV1JoinRoom)
	lti_v1.Get("/room/check", controllers.HandleLTIV1CheckRoom)

	// all auth group routes require auth header (API key and secret key)
	auth := app.Group("/auth", controllers.HandleAuthHeaderCheck)
	auth.Post("/get-client-files", controllers.HandleGetClientFiles)

	room := auth.Group("/room")
	room.Post("/create", controllers.HandleRoomCreate)
	room.Get("/generate-join-token", controllers.HandleGenerateJoinToken)
	room.Get("/room-activity", controllers.HandleRoomActivity)
	room.Get("/active-room-info", controllers.HandleActiveRoomInfo)
	room.Get("/active-rooms-info", controllers.HandleActiveRoomsInfo)
	room.Post("/end", controllers.HandleEndRoom)

	recording := auth.Group("/recording")
	recording.Post("/fetch", controllers.HandleFetchRecordings)
	recording.Post("/delete", controllers.HandleDeleteRecording)
	recording.Get("/generate-download-token", controllers.HandleDownloadToken)

	api := app.Group("/token", controllers.HandleVerifyHeaderToken)
	api.Post("/verify", controllers.HandleVerifyToken)
	api.Post("/renew", controllers.HandleRenewToken)
	api.Post("/revoke", controllers.HandleRevokeToken)
}

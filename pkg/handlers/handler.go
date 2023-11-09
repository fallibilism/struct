package handlers

import (
	"log"
	"net"
	"time"
	"v/pkg/config"
	"v/pkg/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/goccy/go-json"
	jwtware "github.com/gofiber/contrib/jwt"
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

	// changed to TCP listener to support Railway deployment

	addr, err := net.Listen("tcp", ":"+config.Server.Port)

	if err != nil {
		return err
	}

	return r.App.Listener(addr)
}
func (r *Routes) Shutdown() error {

	log.Printf("Gracefully shutting down server...")

	// close database connection
	// close redis connection
	if err := r.App.Shutdown();err != nil {
		return err
	}

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

	app.Post("/login", func(c *fiber.Ctx) error {
				// Create the Claims
		queries := c.Queries()
		user_id, ok := queries["user_id"]
		if !ok {
			return fiber.ErrUnauthorized
		}
		admin := queries["admin"]
		claims := jwt.MapClaims{
			"user_id":  user_id,
			"admin": admin == "true",
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString(config.RsaPrivateKey)
		if err != nil {
			log.Printf("token.SignedString: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t})
	})

	app.Get("/info", func(c *fiber.Ctx) error {
		return c.Render("info", nil)
	})

	app.Get("/health-check", controllers.HandleHealthCheck)
	app.Post("/webhook", controllers.HandleWebhook)
	app.Get("/download/uploadedFile/:sid/*", controllers.HandleDownloadUploadedFile)
	app.Get("/download/recording/:token", controllers.HandleDownloadRecording)

	// all auth group routes require auth header (API key and secret key)
//	auth := app.Group("/auth", controllers.HandleAuthHeaderCheck)
//	auth.Post("/get-client-files", controllers.HandleGetClientFiles)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    config.RsaPrivateKey.Public(),
		},
	}))

	room := app.Group("/room")
	room.Post("/create", controllers.HandleRoomCreate)
	room.Get("/generate-token", controllers.HandleGenerateJoinToken)
//	room.Get("/room-activity", controllers.HandleRoomActivity)
//	room.Get("/active-room-info", controllers.HandleActiveRoomInfo)
//	room.Post("/end", controllers.HandleEndRoom)

	recording := app.Group("/recording")
	recording.Post("/fetch", controllers.HandleFetchRecordings)
	recording.Post("/delete", controllers.HandleDeleteRecording)
	recording.Get("/generate-download-token", controllers.HandleDownloadToken)

	api := app.Group("/token", controllers.HandleVerifyHeaderToken)
	api.Post("/verify", controllers.HandleVerifyToken)
	api.Post("/renew", controllers.HandleRenewToken)
	api.Post("/revoke", controllers.HandleRevokeToken)


	// redirect unknown pages to 404 page

	/*
	app.Get("*", func(c *fiber.Ctx) error {
		return c.Render("404", nil)

	})

	app.Post("*", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"error": "404",
			"msg":   "Page not found",
		})
	})
	*/

}



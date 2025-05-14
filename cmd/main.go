package main

import (
	"stegano-webapp/steagano-webapp/config"
	"stegano-webapp/steagano-webapp/internal/encrypt"
	"stegano-webapp/steagano-webapp/internal/home"
	"stegano-webapp/steagano-webapp/internal/profile"
	"stegano-webapp/steagano-webapp/internal/transcribe"
	"stegano-webapp/steagano-webapp/pkg/database"
	"stegano-webapp/steagano-webapp/pkg/logger"
	"stegano-webapp/steagano-webapp/pkg/middleware"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
)

func main() {

	//Configs
	config.Init()
	logConfig := config.NewLogConfig()
	databaseConfig := config.NewDBConfig()
	// authConfig := config.NewAuthConfig()

	//Logger
	customLogger := logger.NewLogger(logConfig)

	//App
	app := fiber.New(fiber.Config{
		Prefork:           true,
		StreamRequestBody: true,
	})

	//Middlewares
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: customLogger,
	}))
	app.Use(recover.New())
	app.Static("/public", "./public")

	dbpool := database.CreateDataBasePool(databaseConfig, customLogger)
	defer dbpool.Close()

	//Sessions
	storage := postgres.New(postgres.Config{
		DB:         dbpool,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	})

	store := session.New(session.Config{
		Storage: storage,
	})

	app.Static("/static", "./static")

	app.Use(middleware.AuthMiddleware(store))

	//Repositories
	homeRepository := home.NewUsersRepository(dbpool, customLogger)
	encryptRepository := encrypt.NewEncryptRepository(dbpool, customLogger)
	transcribeRepository := transcribe.NewEncryptRepository(dbpool, customLogger)
	profileRepository := profile.NewProfileRepository(dbpool, customLogger)

	//Handlers
	home.NewHandler(app, customLogger, homeRepository, store)
	encrypt.NewEncryptHandler(app, customLogger, encryptRepository, store)
	transcribe.NewTranscribeHandler(app, customLogger, transcribeRepository, store)
	profile.NewProfileHandler(app, customLogger, profileRepository, store)



	app.Listen(":3001")
}

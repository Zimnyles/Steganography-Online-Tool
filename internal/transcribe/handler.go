package transcribe

import (
	"fmt"
	"net/http"
	filenamegenerator "stegano-webapp/steagano-webapp/pkg/filenameGenerator"
	"stegano-webapp/steagano-webapp/pkg/middleware"
	"stegano-webapp/steagano-webapp/pkg/tadapter"
	"stegano-webapp/steagano-webapp/pkg/transcribe"
	"stegano-webapp/steagano-webapp/views"
	"stegano-webapp/steagano-webapp/views/components"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type TranscribeHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *TranscribeRepository
	store        *session.Store
}

func NewTranscribeHandler(router fiber.Router, customLogger *zerolog.Logger, repository *TranscribeRepository, store *session.Store) {
	h := &TranscribeHandler{
		router:       router,
		customLogger: customLogger,
		repository:   repository,
		store:        store,
	}

	//GET
	h.router.Get("/transcribe", middleware.AuthMiddleware(store), h.transcribe)
	//POST
	h.router.Post("/api/createtranscribe", middleware.AuthMiddleware(store), h.apiCreateTranscribe)
}

func (h *TranscribeHandler) transcribe(c *fiber.Ctx) error {
	component := views.TranscribePage()

	return tadapter.Render(c, component, http.StatusOK)

}

func (h *TranscribeHandler) apiCreateTranscribe(c *fiber.Ctx) error {

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	login := sess.Get("login")
	if login == nil {
		c.Response().Header.Add("Hx-Redirect", "/login")
		return c.Redirect("/login", http.StatusOK)
	}
	authedLogin := sess.Get("login").(string)
	h.customLogger.Info().Msg(authedLogin)

	image, err := c.FormFile("image")

	if err != nil || image == nil {
		component := components.Notification("Выберите изображение", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	uniqueFilenameCode := filenamegenerator.GenerateFilename()
	uniqueFilename := "image_" + uniqueFilenameCode + strconv.FormatInt((time.Now().Unix()), 10) + ".jpg"
	filepath := "static/images/transimages/" + uniqueFilename

	if err := c.SaveFile(image, filepath); err != nil {
		component := components.Notification("Ошибка сервера при сохранении изображения, попробуйте позже", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusInternalServerError)
	}

	message, err := transcribe.TranscribeImage(filepath)
	if err != nil {
		fmt.Println(err)
	}

	if message != "" {
		err = h.repository.AllUsersCounterPlus()
		if err != nil {
			h.customLogger.Info().Msg("failed to update all users action counter, handler level")
		}
		err = h.repository.userTranscribeCounterPlus(authedLogin)
		if err != nil {
			h.customLogger.Info().Msg("failed to update transcribed counter, handler level")
		}
	}

	component := components.TranscribeResult(message)
	return tadapter.Render(c, component, http.StatusOK)
}

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
	fmt.Println(message)

	component := components.TranscribeResult(message)
	return tadapter.Render(c, component, http.StatusOK)
}

package encrypt

import (
	"fmt"
	"net/http"
	"os"
	"stegano-webapp/steagano-webapp/pkg/encrypt"
	filenamegenerator "stegano-webapp/steagano-webapp/pkg/filenameGenerator"
	"stegano-webapp/steagano-webapp/pkg/tadapter"
	"stegano-webapp/steagano-webapp/views"
	"stegano-webapp/steagano-webapp/views/components"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type EncryptHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *EncryptRepository
	store        *session.Store
}

func NewEncryptHandler(router fiber.Router, customLogger *zerolog.Logger, repository *EncryptRepository, store *session.Store) {
	handler := &EncryptHandler{
		router:       router,
		customLogger: customLogger,
		repository:   repository,
		store:        store,
	}

	//Get
	handler.router.Get("/encrypt", handler.authMiddleware, handler.encrypt)

	//POST
	handler.router.Post("/api/createencrypt", handler.apiCreateEncrypt)
}

func (h *EncryptHandler) apiCreateEncrypt(c *fiber.Ctx) error {
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

	textToEncrypt := c.FormValue("text")
	h.customLogger.Info().Msg(textToEncrypt)

	image, err := c.FormFile("image")

	if image == nil && textToEncrypt == "" {
		component := components.Notification("Введите текст и выберите изображение", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}
	if textToEncrypt == "" {
		component := components.Notification("Введите текст", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}
	if err != nil || image == nil {
		component := components.Notification("Выберите изображение", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	uniqueFilenameCode := filenamegenerator.GenerateFilename()
	uniqueFilename := "image_" + uniqueFilenameCode + strconv.FormatInt((time.Now().Unix()), 10) + ".png"
	filepath := "static/images/encryimages/" + uniqueFilename

	if err := c.SaveFile(image, filepath); err != nil {
		component := components.Notification("Ошибка сервера при сохранении изображения, попробуйте позже", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusInternalServerError)
	}

	isJPEG, _ := isJPEG(filepath)
	if isJPEG {
		component := components.Notification("Файл должен быть .png", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusInternalServerError)
	}

	encryptedImageFilepath, encryptedImageFilename, _ := encrypt.EncryptImage(filepath, textToEncrypt)
	if encryptedImageFilepath != "" {
		err = h.repository.AllUsersCounterPlus()
		if err != nil {
			h.customLogger.Info().Msg("failed to update all users action counter, handler level")
		}
		err = h.repository.userEncryptCounterPlus(authedLogin)
		if err != nil {
			h.customLogger.Info().Msg("failed to update encrypted counter, handler level")
		}
		component := components.EncryptResult(encryptedImageFilepath, encryptedImageFilename)
		// component := components.Notification("Успех! Ваше послание зашифровано, скачивание начнется автоматически", components.NotificationSuccess)
		return tadapter.Render(c, component, http.StatusOK)
	}

	fmt.Println(encryptedImageFilepath + "test")

	//Counters

	//Return

	component := views.ErrorPage(200, "успех")
	return tadapter.Render(c, component, 200)

}

func (h *EncryptHandler) encrypt(c *fiber.Ctx) error {
	component := views.EncryptPage()

	return tadapter.Render(c, component, http.StatusOK)
}

func (h *EncryptHandler) authMiddleware(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	loginVal := sess.Get("login")
	if loginVal == nil {
		c.Response().Header.Add("Hx-Redirect", "/login")
		return c.Redirect("/login", fiber.StatusSeeOther)
	}
	login, ok := loginVal.(string)
	if !ok || login == "" {

		c.Response().Header.Add("Hx-Redirect", "/login")
		return c.Redirect("/login", fiber.StatusSeeOther)
	}

	return c.Next()
}

func isJPEG(filepath string) (bool, error) {
	file, _ := os.Open(filepath)
	buf := make([]byte, 2)
	_, err := file.Read(buf)
	if err != nil {
		return false, err
	}
	defer file.Close()
	return buf[0] == 0xFF && buf[1] == 0xD8, nil
}

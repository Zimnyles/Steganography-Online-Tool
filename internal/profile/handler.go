package profile

import (
	"fmt"
	"net/http"
	"stegano-webapp/steagano-webapp/pkg/middleware"
	"stegano-webapp/steagano-webapp/pkg/tadapter"
	"stegano-webapp/steagano-webapp/views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type ProfileHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *ProfileRepository
	store        *session.Store
}

func NewProfileHandler(router fiber.Router, customLogger *zerolog.Logger, repository *ProfileRepository, store *session.Store) {
	h := &ProfileHandler{
		router:       router,
		customLogger: customLogger,
		repository:   repository,
		store:        store,
	}

	//GET
	h.router.Get("/profile/:username", middleware.AuthMiddleware(store), h.profile)

}

func (h ProfileHandler) profile(c *fiber.Ctx) error {
	username := c.Params("username")
	isLoginExists, _ := h.repository.IsLoginExistsForString(username, h.customLogger)
	if !isLoginExists {
		component := views.ErrorPage(http.StatusNotFound, "Страница не найдена")
		return tadapter.Render(c, component, http.StatusNotFound)
	}
	userData, err := h.repository.GetUserDataFromLogin(username, h.customLogger)
	if err != nil {
		component := views.ErrorPage(http.StatusInternalServerError, "cannot get userdata, try later")
		return tadapter.Render(c, component, http.StatusInternalServerError)
	}
	fmt.Println(userData)
	component := views.ProfilePage(*userData)
	return tadapter.Render(c, component, http.StatusOK)

}

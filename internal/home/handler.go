package home

import (
	"fmt"
	"net/http"
	"stegano-webapp/steagano-webapp/internal/models"
	"stegano-webapp/steagano-webapp/pkg/tadapter"
	"stegano-webapp/steagano-webapp/pkg/validator"
	"stegano-webapp/steagano-webapp/views"
	"stegano-webapp/steagano-webapp/views/components"

	"strings"

	"github.com/a-h/templ"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type HomeHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
	repository   *UsersRepository
	store        *session.Store
}

func NewHandler(router fiber.Router, customLogger *zerolog.Logger, repository *UsersRepository, store *session.Store) {
	handler := &HomeHandler{
		router:       router,
		customLogger: customLogger,
		repository:   repository,
		store:        store,
	}
	//GET

	handler.router.Get("/", handler.authMiddleware, handler.home)
	handler.router.Get("/login", handler.login)
	handler.router.Get("/registration", handler.registation)
	handler.router.Get("api/logout", handler.apiLogout)

	//POST
	handler.router.Post("api/login", handler.apiLogin)
	handler.router.Post("api/registration", handler.apiRegistration)
}

func (handler *HomeHandler) home(c *fiber.Ctx) error {
	sess, err := handler.store.Get(c)
	if err != nil {
		panic(err)
	}
	login := sess.Get("login").(string)
	userData, err := handler.repository.getUserData(login, handler.customLogger)
	if err != nil {
		panic(err)
	}

	encryptedCounter, err := handler.repository.getEncryptedCounter(login, handler.customLogger)
	if err != nil {
		handler.customLogger.Info().Msg("cannot get encrypted counter(user)")
	}

	transcribedCounter, err := handler.repository.getTranscribedCounter(login, handler.customLogger)
	if err != nil {
		handler.customLogger.Info().Msg("cannot get transcribed counter(user)")
	}

	allUserActions, err := handler.repository.getUserActionsCounter(login, handler.customLogger)
	if err != nil {
		fmt.Println(err)
	}

	var counters models.Counters = models.Counters{
		UserEncrypted: *encryptedCounter,
		UserTranscribed: *transcribedCounter,
		AllUsersActions: *allUserActions,
	}

	component := views.Home(*userData, counters)
	return tadapter.Render(c, component, http.StatusOK)
}

func (h *HomeHandler) apiLogout(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Delete("login")
	if err := sess.Save(); err != nil {
		panic(err)
	}
	c.Response().Header.Add("Hx-Redirect", "/login")
	return c.Redirect("/login", http.StatusOK)
}

func (h *HomeHandler) login(c *fiber.Ctx) error {
	component := views.Login()
	// a
	return tadapter.Render(c, component, http.StatusOK)
}

func (h *HomeHandler) registation(c *fiber.Ctx) error {
	component := views.Registration()

	return tadapter.Render(c, component, http.StatusOK)
}

func (h *HomeHandler) apiLogin(c *fiber.Ctx) error {

	form := LoginForm{
		Login:    c.FormValue("login"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	emailIsExists, _ := h.repository.IsEmailExistsForLogin(form, h.customLogger)

	if !emailIsExists {
		component := components.Notification("Пользователся с такой почтой не существует", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	UserCredentials, _ := h.repository.GetPasswordByEmail(form, h.customLogger)
	if UserCredentials == nil {
		h.customLogger.Info().Msg("ошбика сервера 1")
		component := components.Notification("Ошибка сервера, попробуйте позже", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	err := bcrypt.CompareHashAndPassword([]byte(UserCredentials.PasswordHash), []byte(form.Password))
	if err != nil {
		h.customLogger.Info().Msg("ошбика сервера 2")
		component := components.Notification("Неверный пароль", components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Set("login", strings.ToLower(form.Login))
	if err := sess.Save(); err != nil {
		panic(err)
	}
	c.Response().Header.Add("Hx-Redirect", "/")
	return c.Redirect("/", http.StatusOK)

}

func (h *HomeHandler) authMiddleware(c *fiber.Ctx) error {
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

func (h *HomeHandler) apiRegistration(c *fiber.Ctx) error {

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	form := UserCreateForm{
		Login:    c.FormValue("login"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	errors := validate.Validate(
		&validators.EmailIsPresent{Name: "Email", Field: form.Email, Message: "Email не задан или задан неверно"},
		&validators.StringIsPresent{Name: "Password", Field: form.Password, Message: "Пароль не задан или задан неверно"},
		&validators.StringIsPresent{Name: "Login", Field: form.Login, Message: "Логин не задан или задан неверно"},
	)

	var component templ.Component
	if len(errors.Error()) > 0 {
		component = components.Notification(validator.FormatErrors(errors), components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	msg, err := h.repository.addUser(form, h.customLogger)
	if err != nil {
		h.customLogger.Error().Msg(err.Error())
		component = components.Notification(msg, components.NotificationFail)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}

	sess.Set("login", strings.ToLower(form.Login))
	if err := sess.Save(); err != nil {
		panic(err)
	}

	c.Response().Header.Add("Hx-Redirect", "/")
	return c.Redirect("/", http.StatusOK)

}

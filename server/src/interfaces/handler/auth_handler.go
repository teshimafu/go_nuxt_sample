package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/teshimafu/lazyPM/server/src/domain/factory"
	"github.com/teshimafu/lazyPM/server/src/domain/valueobject"
	"github.com/teshimafu/lazyPM/server/src/interfaces/presenter"
	"github.com/teshimafu/lazyPM/server/src/usecase/service"
)

type SignupForm struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SigninForm struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthHandler struct {
	authService   *service.AuthService
	userPresenter *presenter.UserPresenter
	userFactory   factory.UserFactory
}

func NewAuthHandler(authService *service.AuthService, userPresenter *presenter.UserPresenter, userFactory factory.UserFactory) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		userPresenter: userPresenter,
		userFactory:   userFactory,
	}
}

func (a *AuthHandler) PostSignup(c echo.Context) error {
	userCmd := &SignupForm{}
	if err := c.Bind(userCmd); err != nil {
		return err
	}
	user, err := a.userFactory.CreateUser(userCmd.Name, userCmd.Email, userCmd.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	createdUser, err := a.authService.Signup(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, a.userPresenter.ToJson(createdUser))
}

func (a *AuthHandler) PostSignin(c echo.Context) error {
	req := &SigninForm{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	password, err := valueobject.NewPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	token, err := a.authService.Signin(email, password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zxcghoulhunter/InnoTaxi/internal/model"
	"github.com/zxcghoulhunter/InnoTaxi/internal/service/AuthService"
)

type Handler struct {
	service AuthService.Service
}

func NewHandler(service AuthService.Service) Handler {
	return Handler{service: service}
}

//@Summary Register
//@Tags auth
//@Description create user
//@ID create-user
//@Accept json
//@Produce json
//@Param input body model.Login true "account info"
//@Router /register [post]
//@Success 200 {string}
func (h *Handler) Register(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.handleError(err, c)
		return
	}
	err = h.service.Save(c.Request.Context(), user)
	if err != nil {
		h.handleError(err, c)
		return
	}
	err = h.SignIn(c, model.Login{Phone: user.Phone, Password: user.Password})
	if err != nil {
		h.handleError(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user.Phone})
}

//@Summary Login
//@Tags auth
//@Description login user
//@ID login-user
//@Accept json
//@Produce json
//@Param input body model.Login true "account info"
//@Router /login [post]
//@Success 200 {string}
func (h *Handler) Login(c *gin.Context) {
	var user model.Login
	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.handleError(err, c)
		return
	}
	err = h.service.Authorize(c.Request.Context(), user)
	if err != nil {
		h.handleError(err, c)
		return
	}
	err = h.SignIn(c, user)
	if err != nil {
		h.handleError(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user.Phone})
}

//@Summary Logout
//@Tags auth
//@Description logout user
//@ID logout-user
//@Accept json
//@Produce json
//@Router /logout [get]
//@Success 200  {string}
func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"Status": "logout"})
}

func (h *Handler) handleError(err error, c *gin.Context) {
	switch {
	case errors.Is(err, AuthService.ErrWrongPassword):
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong password"})
	case errors.Is(err, AuthService.ErrUserAlreadyExists):
		c.JSON(http.StatusBadRequest, gin.H{"message": "User already exists"})
	case errors.Is(err, AuthService.ErrUserNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

	}

}

func (h *Handler) SignIn(c *gin.Context, user model.Login) error {
	tokenString, expirationTime, err := AuthService.SignIn(c.Request.Context(), user)
	if err != nil {
		return err
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	return nil
}

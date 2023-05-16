package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/shamank/eduTour-backend/app/internal/service"
	"net/http"
)

type userSignUpInput struct {
	Name     string `json:"name" validate:"required,min=2,max=64"`
	Email    string `json:"email" validate:"required,email,max=64"`
	Phone    string `json:"phone" validate:"required,phone,max=13"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type userSignInInput struct {
	Name     string `json:"name" validate:"required,min=2,max=64"`
	Email    string `json:"email" validate:"required,email,max=64"`
	Phone    string `json:"phone" validate:"required,phone,max=13"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpireIn     int    `json:"expire_in"`
}

func (h *Handler) initAuthRouter(api *gin.RouterGroup) {
	auth := api.Group("auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.userRefresh)
	}
}

func (h *Handler) signUp(c *gin.Context) {
	var input userSignUpInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (h *Handler) signIn(c *gin.Context) {

	var input userSignInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.services.Users.SignIn(c.Request.Context(), service.UserSignInInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpireIn:     int(res.ExpireIn.Seconds()),
	})
}

type refreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) userRefresh(c *gin.Context) {
	var input refreshInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid body")
		return
	}

	res, err := h.services.Users.RefreshTokens(c.Request.Context(), input.RefreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpireIn:     int(res.ExpireIn.Seconds()),
	})
}

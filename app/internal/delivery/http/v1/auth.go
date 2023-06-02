package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shamank/eduTour-backend/app/internal/domain"
	"github.com/shamank/eduTour-backend/app/internal/service"
	"net/http"
)

type userSignUpInput struct {
	UserName string `json:"username" binding:"required,min=4,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type userSignInInput struct {
	//UserName string `json:"username" validate:"required,min=2,max=64"`
	Email string `json:"email" binding:"required,email,max=64"`
	//Phone    string `json:"phone" validate:"required,phone,max=13"`
	Password string `json:"password" binding:"required,min=8,max=64"`
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

		auth.GET("/me", h.userIdentity, h.userPing)
		auth.GET("/verify", h.userIdentity, h.verifyToken)
	}
}

// @Summary User SignUp
// @Tags auth
// @Description create user account
// @ModuleID authSignUp
// @Accept  json
// @Produce  json
// @Param input body userSignUpInput true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input userSignUpInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.Authorization.SignUp(c.Request.Context(), service.UserSignUpInput{
		UserName: input.UserName,
		Email:    input.Email,
		Password: input.Password,
	}); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary User SignIn
// @Tags auth
// @Description user sign in
// @ModuleID authSignIn
// @Accept  json
// @Produce  json
// @Param input body userSignInInput true "sign in info"
// @Success 200 {object} tokenResponse
// @Failure 400,401,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {

	var input userSignInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.services.Authorization.SignIn(c.Request.Context(), service.UserSignInInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
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
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// @Summary Refresh Token
// @Tags auth
// @Description user refresh token
// @ModuleID authRefreshToken
// @Accept  json
// @Produce  json
// @Param input body refreshInput true "refresh token input"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/refresh [post]
func (h *Handler) userRefresh(c *gin.Context) {
	var input refreshInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid body")
		return
	}

	res, err := h.services.Authorization.RefreshToken(c.Request.Context(), input.RefreshToken)
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

// @Summary User check token
// @Tags auth
// @Description user check access token
// @ModuleID authPing
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/me [get]
func (h *Handler) userPing(c *gin.Context) {
	_, err := h.parseAuthHeader(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary Verify token for other apps
// @Tags backend
// @Description verify token for other apps
// @ModuleID authVerify
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/verify [get]
func (h *Handler) verifyToken(c *gin.Context) {
	usr, err := h.parseAuthHeader(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res, err := h.services.Authorization.GetFullUserInfo(c.Request.Context(), usr.userID)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)

}

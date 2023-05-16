package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/shamank/eduTour-backend/app/internal/domain"
	"net/http"
)

func (h *Handler) initAuthRouter(api *gin.RouterGroup) {
	auth := api.Group("auth")
	{
		auth.POST("/sign-up")
		auth.POST("/sign-in")
		auth.POST("/refresh")
	}
}

func (h *Handler) signUp(c *gin.Context) {
	var input domain.UserSignUp
	if err := c.BindJSON(&input); err != nil {
		newErrorReponse(c, http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (h *Handler) signIn(c *gin.Context) {

	var input domain.UserSignIn

	if err := c.BindJSON(&input); err != nil {
		newErrorReponse(c, http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

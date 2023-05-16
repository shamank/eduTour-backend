package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/shamank/eduTour-backend/app/internal/service"
	"github.com/shamank/eduTour-backend/app/pkg/auth"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) InitAPI(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initAuthRouter(v1)
	}
}

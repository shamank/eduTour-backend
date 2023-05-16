package http

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/shamank/eduTour-backend/app/internal/delivery/http/v1"
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

func (h *Handler) InitAPI() *gin.Engine {
	router := gin.Default()

	router.Use(CORS)

	handlerV1 := v1.NewHandler(h.services, h.tokenManager)

	api := router.Group("/api")
	{
		handlerV1.InitAPI(api)
	}

	return router

}

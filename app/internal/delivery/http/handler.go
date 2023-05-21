package http

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/shamank/eduTour-backend/app/internal/delivery/http/v1"
	"github.com/shamank/eduTour-backend/app/internal/service"
	"github.com/shamank/eduTour-backend/app/pkg/auth"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/shamank/eduTour-backend/docs"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		handlerV1.InitAPI(api)
	}

	return router

}

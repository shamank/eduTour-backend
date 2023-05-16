package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shamank/eduTour-backend/app/internal/service"
	"github.com/shamank/eduTour-backend/app/pkg/auth"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
	validator    *validator.Validate
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	validate := validator.New()
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
		validator:    validate,
	}
}

func (h *Handler) InitAPI(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initAuthRouter(v1)
		h.initEventRouter(v1)
		h.initEventCategoriesRouter(v1)
	}
}

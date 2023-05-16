package v1

import "github.com/gin-gonic/gin"

func (h *Handler) initEventCategoriesRouter(api *gin.RouterGroup) {
	categories := api.Group("categories", h.userIdentity, h.adminOnly)
	{
		categories.POST("/", h.createCategory)
		categories.PUT("/:id", h.updateCategoryByID)
		categories.DELETE("/:id", h.deleteCategoryByID)
	}
}
func (h *Handler) createCategory(c *gin.Context) {

}

func (h *Handler) updateCategoryByID(c *gin.Context) {

}

func (h *Handler) deleteCategoryByID(c *gin.Context) {

}

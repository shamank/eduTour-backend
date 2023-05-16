package v1

import "github.com/gin-gonic/gin"

func (h *Handler) initEventRouter(api *gin.RouterGroup) {
	events := api.Group("events")
	{

		events.GET("/", h.getAllEvents)
		events.POST("/", h.userIdentity, h.adminOnly, h.createEvent)

		events.GET("/:id", h.getEventByID)
		events.PUT("/:id", h.userIdentity, h.adminOnly, h.updateEventByID)
		events.DELETE("/:id", h.userIdentity, h.adminOnly, h.deleteEventByID)
	}
}
func (h *Handler) getAllEvents(c *gin.Context) {

}

func (h *Handler) createEvent(c *gin.Context) {

}

func (h *Handler) getEventByID(c *gin.Context) {

}

func (h *Handler) updateEventByID(c *gin.Context) {

}

func (h *Handler) deleteEventByID(c *gin.Context) {

}

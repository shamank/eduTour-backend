package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
	userCtx             = "userID"
	roleCtx             = "role"

	adminRole = "admin"
)

func (h *Handler) parseAuthHeader(c *gin.Context) (string, string, error) {
	header := c.GetHeader(AuthorizationHeader)
	if header == "" {
		return "", "", errors.New("auth header is empty")
	}
	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", "", errors.New("auth header is invalid")
	}

	if len(headerParts[1]) == 0 {
		return "", "", errors.New("auth token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}

func (h *Handler) userIdentity(c *gin.Context) {
	id, role, err := h.parseAuthHeader(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, id)
	c.Set(roleCtx, role)
}

func (h *Handler) adminOnly(c *gin.Context) {
	role, ok := c.Get(roleCtx)
	if !ok || role != adminRole {
		newErrorResponse(c, http.StatusForbidden, "you are not admin")
		return
	}
}

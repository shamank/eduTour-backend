package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	AuthorizationHeader = "Authorization"
	userCtx             = "userID"
	roleCtx             = "role"

	adminRole = "admin"
)

type userContext struct {
	userID   int
	userName string
	Role     string
}

func (h *Handler) parseAuthHeader(c *gin.Context) (userContext, error) {
	header := c.GetHeader(AuthorizationHeader)
	if header == "" {
		return userContext{}, errors.New("auth header is empty")
	}
	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return userContext{}, errors.New("auth header is invalid")
	}

	if len(headerParts[1]) == 0 {
		return userContext{}, errors.New("auth token is empty")
	}

	res, err := h.tokenManager.Parse(headerParts[1])
	if err != nil {
		return userContext{}, err
	}
	if res.ExpireAt < time.Now().Unix() {
		return userContext{}, errors.New("token expired")
	}

	return userContext{
		userID:   res.UserID,
		userName: res.UserName,
		Role:     res.Role,
	}, nil
}

func (h *Handler) userIdentity(c *gin.Context) {
	usr, err := h.parseAuthHeader(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, usr)
}

func (h *Handler) adminOnly(c *gin.Context) {
	role, ok := c.Get(roleCtx)
	if !ok {
		newErrorResponse(c, http.StatusForbidden, "you are not login")
		return
	}

	if role != adminRole {
		newErrorResponse(c, http.StatusForbidden, "you are not admin")
		return
	}

	return

}

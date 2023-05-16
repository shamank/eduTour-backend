package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorReponse(c *gin.Context, statusCode int, msg string) {
	logrus.Error(msg)

	c.AbortWithStatusJSON(statusCode, errorResponse{Message: msg})
}

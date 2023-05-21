package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/shamank/eduTour-backend/app/internal/service"
	"net/http"
)

type userProfileOutput struct {
	UserName   string   `json:"username"`
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	MiddleName string   `json:"middle_name"`
	Avatar     string   `json:"avatar"`
	Roles      []string `json:"roles"`
}

type userProfileInput struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Avatar     string `json:"avatar"`
}

func (h *Handler) initUsersRouter(api *gin.RouterGroup) {
	users := api.Group("users")
	{
		users.GET("/profile/:user_name", h.getUserProfile)
		users.PUT("/profile/:user_name", h.userIdentity, h.updateUserProfile)
	}
}

// @Summary Get Profile
// @Tags users
// @Description get user profile
// @ModuleID userGetProfile
// @Accept  json
// @Produce  json
// @Param username path string true "username"
// @Success 200 {object} userProfileOutput
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/profile/ [get]
func (h *Handler) getUserProfile(c *gin.Context) {

	userName := c.Param("user_name")

	res, err := h.services.Users.GetUserProfile(c.Request.Context(), userName)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, userProfileOutput{
		UserName:   res.UserName,
		FirstName:  res.FirstName,
		LastName:   res.LastName,
		MiddleName: res.MiddleName,
		Avatar:     res.Avatar,
		Roles:      res.Roles,
	})

}

// @Summary Update Profile
// @Tags users
// @Description update user profile
// @ModuleID userUpdateProfile
// @Accept  json
// @Produce  json
// @Param username path string true "username"
// @Security ApiKeyAuth
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/profile/ [put]
func (h *Handler) updateUserProfile(c *gin.Context) {
	var input userProfileInput

	userName := c.Param("user_name")

	usr, err := h.parseAuthHeader(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	IsAdmin := false
	for _, role := range usr.Roles {
		if role == adminRole {
			IsAdmin = true
			break
		}
	}

	if usr.userName != userName && !IsAdmin {
		newErrorResponse(c, http.StatusForbidden, "permission denied")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.services.Users.UpdateUserProfile(c.Request.Context(), userName, service.UserProfileInput{
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		MiddleName: input.MiddleName,
		Avatar:     input.Avatar,
	})
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"ok"})
}

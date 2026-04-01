package handlers

import (
	"homework04/models"
	"homework04/services"
	"homework04/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService, jwtSecret []byte) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// 注册函数
func (h *UserHandler) Register(c *gin.Context) {

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBizError(400, err.Error()))
		return
	}

	userID, err := h.userService.Register(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    200,
		Message: "success",
		Data:    gin.H{"user_id": userID},
	})
}

// // 登录函数
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBizError(400, err.Error()))
		return
	}

	token, err := h.userService.Login(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    200,
		Message: "login success",
		Data:    gin.H{"token": token},
	})
}

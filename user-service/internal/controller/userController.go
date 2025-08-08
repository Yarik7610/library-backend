package controller

import (
	"net/http"

	"github.com/Yarik7610/library-backend/user-service/internal/dto"
	"github.com/Yarik7610/library-backend/user-service/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	SignUp(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{userService: userService}
}

func (c *userController) SignUp(ctx *gin.Context) {
	var signUpUserDTO dto.SignUpUser

	if err := ctx.ShouldBindJSON(&signUpUserDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.SignUp(&signUpUserDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

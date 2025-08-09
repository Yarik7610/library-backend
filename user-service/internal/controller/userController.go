package controller

import (
	"net/http"

	"github.com/Yarik7610/library-backend/user-service/internal/dto"
	"github.com/Yarik7610/library-backend/user-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
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
		zap.S().Error("Sign up error: ", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *userController) SignIn(ctx *gin.Context) {
	var signInUserDTO dto.SignInUser

	if err := ctx.ShouldBindJSON(&signInUserDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.userService.SignIn(&signInUserDTO)
	if err != nil {
		zap.S().Error("Sign in error: ", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

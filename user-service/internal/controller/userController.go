package controller

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/user-service/internal/dto"
	"github.com/Yarik7610/library-backend/user-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	GetMe(ctx *gin.Context)
	GetEmailsByUserIDs(ctx *gin.Context)
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
		zap.S().Errorf("Sign up error: %v\n", err)
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
		zap.S().Errorf("Sign in error: %v\n", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *userController) GetMe(ctx *gin.Context) {
	userIDString := ctx.GetHeader(sharedconstants.HEADER_USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, customErr := c.userService.GetMe(uint(userID))
	if customErr != nil {
		zap.S().Errorf("Me error: %v\n", customErr)
		ctx.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *userController) GetEmailsByUserIDs(ctx *gin.Context) {
	userIDsStrings := ctx.QueryArray("ids")

	userIDs := make([]uint, 0)
	for _, s := range userIDsStrings {
		userID, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			zap.S().Errorf("Get emails by user IDs error: %v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userIDs = append(userIDs, uint(userID))
	}

	emails, err := c.userService.GetEmailsByUserIDs(userIDs)
	if err != nil {
		zap.S().Errorf("Get emails error: %v\n", err)
		ctx.JSON(err.Code, gin.H{"error": err.Message})
		return
	}
	ctx.JSON(http.StatusOK, emails)
}

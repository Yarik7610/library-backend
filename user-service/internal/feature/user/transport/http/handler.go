package http

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/service"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/transport/http/dto"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/errs"

	httpInfrastructure "github.com/Yarik7610/library-backend/user-service/internal/infrastructure/transport/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	GetMe(c *gin.Context)
	GetEmailsByUserIDs(c *gin.Context)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}

// SignUp godoc
//
//	@Summary		Register new user
//	@Description	Creates a new user account
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.SignUpUserRequest	true	"Sign up payload"
//	@Success		201		{object}	dto.User
//	@Failure		400 {object} dto.Error "Bad request"
//	@Failure		409 {object} dto.Error "Entity already exists"
//	@Failure		500	{object} dto.Error "Internal server error"
//	@Router			/sign-up [post]
func (h *userHandler) SignUp(c *gin.Context) {
	var SignUpUserRequestDTO dto.SignUpUserRequest
	if err := c.ShouldBindJSON(&SignUpUserRequestDTO); err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	user, err := h.userService.SignUp(&SignUpUserRequestDTO)
	if err != nil {
		zap.S().Errorf("Sign up error: %v\n", err)
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// SignIn godoc
//
//	@Summary		Authorize user
//	@Description	Authorize existings user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.SignInUserRequest	true	"Sign in data"
//	@Success		200		{object}	dto.UserToken
//	@Failure		400 {object} dto.Error "Bad request"
//	@Failure		500	{object} dto.Error "Internal server error"
//	@Router			/sign-in [post]
func (h *userHandler) SignIn(c *gin.Context) {
	var SignInUserRequestDTO dto.SignInUserRequest
	if err := c.ShouldBindJSON(&SignInUserRequestDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.SignIn(&SignInUserRequestDTO)
	if err != nil {
		zap.S().Errorf("Sign in error: %v\n", err)
		c.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetMe godoc
//
//	@Summary		Get current user info
//	@Description	Returns info about the authenticated user
//	@Tags			user
//	@Produce		json
//	@Security 	BearerAuth
//	@Success		200	{object}	dto.User
//	@Failure		401	{object}	map[string]string
//	@Router			/me [get]
func (h *userHandler) GetMe(c *gin.Context) {
	userIDString := c.GetHeader(sharedconstants.HEADER_USER_ID)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, customErr := h.userService.GetMe(uint(userID))
	if customErr != nil {
		zap.S().Errorf("Me error: %v\n", customErr)
		c.JSON(customErr.Code, gin.H{"error": customErr.Message})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetEmailsByUserIDs godoc
//
//	@Summary		Get emails by user IDs
//	@Description	Returns a list of emails for given user IDs
//	@Tags			internal
//	@Produce		json
//	@Param			ids	query		[]int	true	"User IDs"
//	@Success		200	{array}		string
//	@Failure		400	{object}	map[string]string
//	@Router			/emails [get]
func (h *userHandler) GetEmailsByUserIDs(c *gin.Context) {
	userIDsStrings := c.QueryArray("ids")

	userIDs := make([]uint, 0)
	for _, s := range userIDsStrings {
		userID, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			zap.S().Errorf("Get emails by user IDs error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userIDs = append(userIDs, uint(userID))
	}

	emails, err := h.userService.GetEmailsByUserIDs(userIDs)
	if err != nil {
		zap.S().Errorf("Get emails error: %v\n", err)
		c.JSON(err.Code, gin.H{"error": err.Message})
		return
	}
	c.JSON(http.StatusOK, emails)
}

package http

import (
	"net/http"

	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/service"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/transport/http/dto"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/transport/http/mapper"
	"github.com/Yarik7610/library-backend/user-service/internal/feature/user/transport/http/query"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/errs"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/observability/logging"

	httpInfrastructure "github.com/Yarik7610/library-backend/user-service/internal/infrastructure/transport/http"
	"github.com/Yarik7610/library-backend/user-service/internal/infrastructure/transport/http/header"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	GetMe(c *gin.Context)
	GetEmailsByUserIDs(c *gin.Context)
}

type userHandler struct {
	config      *config.Config
	logger      *logging.Logger
	userService service.UserService
}

func NewUserHandler(
	config *config.Config,
	logger *logging.Logger,
	userService service.UserService) UserHandler {
	return &userHandler{
		config:      config,
		logger:      logger,
		userService: userService,
	}
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
//	@Failure		400	{object} dto.Error "Bad request"
//	@Failure		409 {object} dto.Error "Entity already exists"
//	@Failure		500	{object} dto.Error "Internal server error"
//	@Router			/sign-up [post]
func (h *userHandler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	var signUpUserRequestDTO dto.SignUpUserRequest
	if err := c.ShouldBindJSON(&signUpUserRequestDTO); err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	userDomain := mapper.SignUpUserRequestDTOToDomain(&signUpUserRequestDTO)
	if err := h.userService.SignUp(ctx, &userDomain); err != nil {
		h.logger.Error(ctx, "Sign up error", logging.Error(err))
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapper.UserDomainToDTO(&userDomain))
}

// SignIn godoc
//
//	@Summary		Authorize user
//	@Description	Authorize existings user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.SignInUserRequest	true	"Sign in payload"
//	@Success		200	{object}	dto.Token
//	@Failure		400 {object} 	dto.Error "Bad request"
//	@Failure		404 {object} 	dto.Error "Entity not found"
//	@Failure		500	{object} 	dto.Error "Internal server error"
//	@Router			/sign-in [post]
func (h *userHandler) SignIn(c *gin.Context) {
	ctx := c.Request.Context()

	var signInUserRequestDTO dto.SignInUserRequest
	if err := c.ShouldBindJSON(&signInUserRequestDTO); err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	userDomain := mapper.SignInUserRequestDTOToDomain(&signInUserRequestDTO)
	tokenDomain, err := h.userService.SignIn(ctx, &userDomain)
	if err != nil {
		h.logger.Error(ctx, "Sign in error", logging.Error(err))
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.TokenDomainToDTO(tokenDomain))
}

// GetMe godoc
//
//	@Summary		Get current user info
//	@Description	Returns info about the authorized user
//	@Tags			user
//	@Produce		json
//	@Security 	BearerAuth
//	@Success		200	{object}	dto.User
//	@Failure		401 {object} 	dto.Error "The token is missing, invalid or expired"
//	@Failure		404 {object} 	dto.Error "Entity not found"
//	@Failure		500	{object} 	dto.Error "Internal server error"
//	@Router			/me [get]
func (h *userHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := header.GetUserID(c)
	if err != nil {
		httpInfrastructure.RenderError(c, err)
		return
	}

	userDomain, err := h.userService.GetMe(ctx, uint(userID))
	if err != nil {
		h.logger.Error(ctx, "Get me error", logging.Error(err))
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.UserDomainToDTO(userDomain))
}

// GetEmailsByUserIDs godoc
//
//	@Summary		Get emails by user IDs
//	@Description	Returns a list of emails for given user IDs
//	@Tags			internal
//	@Produce		json
//	@Param	ids	query	[]int	true "User IDs" collectionFormat(multi)
//	@Success		200	{array}		string
//	@Failure		400 {object} 	dto.Error "Bad request"
//	@Failure		500	{object} 	dto.Error "Internal server error"
//	@Router			/emails [get]
func (h *userHandler) GetEmailsByUserIDs(c *gin.Context) {
	ctx := c.Request.Context()

	var getEmailsByUserIDsQuery query.GetEmailsByUserIDs
	if err := c.ShouldBindQuery(&getEmailsByUserIDsQuery); err != nil {
		httpInfrastructure.RenderError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	emails, err := h.userService.GetEmailsByUserIDs(ctx, getEmailsByUserIDsQuery.IDs)
	if err != nil {
		h.logger.Error(ctx, "Get emails by user IDs error", logging.Error(err))
		httpInfrastructure.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, emails)
}

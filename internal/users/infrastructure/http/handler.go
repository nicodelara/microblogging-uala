package http

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/nicodelara/microblogging-uala/internal/users/application"
	"github.com/nicodelara/microblogging-uala/internal/users/domain/ports"
)

type userHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) *userHandler {
	return &userHandler{
		service: service,
	}
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func validateCreateUserRequest(c *gin.Context) (*createUserRequest, error) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (h *userHandler) CreateUser(c *gin.Context) {
	req, err := validateCreateUserRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req.Username, req.Email)
	if err != nil {
		switch err {
		case application.ErrUserAlreadyExists, application.ErrEmailAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

type followRequest struct {
	FollowUsername string `json:"followUsername" binding:"required"`
}

func validateFollowRequest(c *gin.Context) (*followRequest, string, error) {
	username := c.Param("username")
	if username == "" {
		return nil, "", errors.New("username is required")
	}

	var req followRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, "", err
	}
	return &req, username, nil
}

func (h *userHandler) FollowUser(c *gin.Context) {
	req, username, err := validateFollowRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	follow, err := h.service.FollowUser(c.Request.Context(), username, req.FollowUsername)
	if err != nil {
		switch err {
		case application.ErrFollowerNotFound, application.ErrFollowingNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case application.ErrAlreadyFollowing:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, follow)
}

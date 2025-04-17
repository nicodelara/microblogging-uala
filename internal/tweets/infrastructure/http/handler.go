package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicodelara/microblogging-uala/internal/tweets/application"
	"github.com/nicodelara/microblogging-uala/internal/tweets/domain/ports"
)

type TweetHandler struct {
	service ports.TweetService
}

func NewTweetHandler(service ports.TweetService) *TweetHandler {
	return &TweetHandler{
		service: service,
	}
}

type createTweetRequest struct {
	Username string `json:"username" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

func validateCreateTweetRequest(c *gin.Context) (*createTweetRequest, error) {
	var req createTweetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	if len(req.Content) > 280 {
		return nil, application.ErrContentTooLong
	}

	return &req, nil
}

func (h *TweetHandler) CreateTweet(c *gin.Context) {
	req, err := validateCreateTweetRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweet, err := h.service.CreateTweet(c.Request.Context(), req.Username, req.Content)
	if err != nil {
		switch err {
		case application.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, tweet)
}

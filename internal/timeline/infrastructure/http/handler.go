package http

import (
	"net/http"
	"strconv"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/nicodelara/microblogging-uala/internal/timeline/application"
	"github.com/nicodelara/microblogging-uala/internal/timeline/domain/ports"
)

type TimelineHandler struct {
	service ports.TimelineService
}

func NewTimelineHandler(service ports.TimelineService) *TimelineHandler {
	return &TimelineHandler{service: service}
}

type getTimelineRequest struct {
	Username string
	Offset   int
	Limit    int
}

type TweetView struct {
	TweetID   string `json:"tweet_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type TimelineResponse struct {
	Username string      `json:"username"`
	Tweets   []TweetView `json:"tweets"`
}

func validateGetTimelineRequest(c *gin.Context) (*getTimelineRequest, error) {
	username := c.Param("username")
	if username == "" {
		return nil, errors.New("username is required")
	}

	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	offset := 0
	limit := 10

	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil {
			offset = v
		}
	}

	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			limit = v
		}
	}

	return &getTimelineRequest{
		Username: username,
		Offset:   offset,
		Limit:    limit,
	}, nil
}

func (h *TimelineHandler) GetTimeline(c *gin.Context) {
	req, err := validateGetTimelineRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweets, err := h.service.GetTimeline(c.Request.Context(), req.Username, req.Offset, req.Limit)
	if err != nil {
		switch err {
		case application.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var tweetsView []TweetView
	for _, tweet := range tweets.Tweets {
		tweetsView = append(tweetsView, TweetView{
			TweetID:   tweet.ID,
			Username:  tweet.Username,
			Content:   tweet.Content,
			CreatedAt: tweet.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, TimelineResponse{
		Username: req.Username,
		Tweets:   tweetsView,
	})
}

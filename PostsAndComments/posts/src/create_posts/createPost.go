package create_posts

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
	"posts/src"
	"posts/src/mongoDB/database"
	"strings"
	"time"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

func CreatePost(c *gin.Context) {
	var post src.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		HandleError(c, NewBadRequestError(err.Error()))
		return
	}

	if post.Title == "" || post.Content == "" {
		HandleError(c, NewBadRequestError("Title and Content are required"))
		return
	}

	post.ID = primitive.NewObjectID()
	post.PostID = uuid.New().String()
	post.CreatedAt = time.Now()

	userID, err := getUserIDFromContext(c)
	if err != nil {
		HandleError(c, NewBadRequestError("Failed to get user ID"))
		return
	}

	post.UserID = userID
	post.ID = primitive.NewObjectID()
	post.CreatedAt = time.Now()

	if err := initMediaFile(post.Image); err != nil {
		HandleError(c, NewInternalServerError(fmt.Sprintf("Failed to initialize Image: %s", err.Error())))
		return
	}

	if err := initMediaFile(post.Video); err != nil {
		HandleError(c, NewInternalServerError(fmt.Sprintf("Failed to initialize Video: %s", err.Error())))
		return
	}

	if err := initMediaFile(post.Gif); err != nil {
		HandleError(c, NewInternalServerError(fmt.Sprintf("Failed to initialize Gif: %s", err.Error())))
		return
	}

	collection := database.Client.Database("blog").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, post)
	if err != nil {
		HandleError(c, NewInternalServerError("Failed to create post: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, post)
}

func initMediaFile(media *src.MediaFile) error {
	if media == nil {
		return nil
	}
	media.ID = primitive.NewObjectID()
	media.CreatedAt = time.Now()

	if _, err := url.ParseRequestURI(media.URL); err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	validFileTypes := []string{"image/jpeg", "image/png", "video/mp4", "image/gif"}
	if !contains(validFileTypes, media.FileType) {
		return fmt.Errorf("invalid file type: %s", media.FileType)
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}

func HandleError(c *gin.Context, err error) {
	var e *AppError
	switch {
	case errors.As(err, &e):
		c.JSON(e.StatusCode, e)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
}

func getUserIDFromContext(c *gin.Context) (int64, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	userIDInt, ok := userID.(int64)
	if !ok {
		return 0, errors.New("user ID is not of type int64")
	}

	return userIDInt, nil
}

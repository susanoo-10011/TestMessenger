package src

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		tokenRecord, err := verifyToken(db, tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Сохранение данных пользователя в контексте
		c.Set("userId", tokenRecord.UserID)
		c.Next()
	}
}

func verifyToken(db *sql.DB, tokenString string) (*TokenRecord, error) {
	var record TokenRecord
	err := db.QueryRow(`
        SELECT id, user_id, token, expires_at, created_at
        FROM tokens
        WHERE token = $1 AND expires_at > NOW()
    `, tokenString).Scan(&record.ID, &record.UserID, &record.Token, &record.ExpiresAt, &record.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("token not found or expired")
		}
		return nil, err
	}

	return &record, nil
}

func CreatePost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if post.Title == "" || post.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Content are required"})
		return
	}

	post.ID = primitive.NewObjectID()
	post.CreatedAt = time.Now()

	if post.Image != nil {
		post.Image.CreatedAt = time.Now()
	}
	if post.Video != nil {
		post.Image.CreatedAt = time.Now()
	}
	if post.Gif != nil {
		post.Image.CreatedAt = time.Now()
	}
	collection := Client.Database("blog").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(http.StatusCreated, post)
}

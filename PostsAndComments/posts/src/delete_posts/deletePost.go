package delete_posts

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrPostNotFound   = errors.New("post not found")
	ErrUnauthorized   = errors.New("unauthorized to delete this post")
	ErrDatabaseError  = errors.New("database error")
	ErrFailedToDelete = errors.New("failed to delete post")
)

func DeletePostFromDB(db *sql.DB, postID string, userID string) error {
	tx, err := db.Begin()
	if err != nil {
		return ErrDatabaseError
	}
	defer tx.Rollback()
	var ownerID string
	err = tx.QueryRow("SELECT user_id FROM posts WHERE post_id = $1 FOR UPDATE", postID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPostNotFound
		}
		return ErrDatabaseError
	}

	if ownerID != userID {
		return ErrUnauthorized
	}

	result, err := tx.Exec("DELETE FROM posts WHERE post_id = $1", postID)
	if err != nil {
		return ErrDatabaseError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrDatabaseError
	}

	if rowsAffected == 0 {
		return ErrFailedToDelete
	}
	if err := tx.Commit(); err != nil {
		return ErrDatabaseError
	}

	return nil
}

func DeletePostHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("id")
		userID, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
			return
		}

		err := DeletePostFromDB(db, postID, userID.(string))
		if err != nil {
			switch err {
			case ErrPostNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case ErrUnauthorized:
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			case ErrFailedToDelete:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": ErrDatabaseError.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
	}
}

package src

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"log"
	"net/http"
	"strings"
)

const (
	ErrAuthHeaderRequired     = "Authorization header is required"
	ErrTokenNotFoundOrExpired = "token not found or expired"
	ErrTokenAlreadyExists     = "token already exists"
	ErrUserNotFound           = "user not found"
	ErrDatabaseError          = "database error"
	ErrUnknownError           = "unknown error"
)

func AuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrAuthHeaderRequired})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		tokenRecord, err := verifyToken(db, tokenString)
		if err != nil {
			switch err.(type) {
			case *TokenError:
				tokenErr := err.(*TokenError)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": tokenErr.Message})
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": ErrUnknownError})
				log.Printf("Error verifying token: %v", err)
			}
			return
		}
		c.Set("userId", tokenRecord.UserID)
		c.Next()
	}
}

func handleDatabaseError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return &TokenError{Message: ErrTokenAlreadyExists, Type: TokenDuplicateError}
		case "23503":
			return &TokenError{Message: ErrUserNotFound, Type: UserNotFoundError}
		default:
			return &TokenError{Message: ErrDatabaseError + ": " + pqErr.Message, Type: DatabaseError}
		}
	}
	return &TokenError{Message: ErrUnknownError + ": " + err.Error(), Type: UnknownError}
}

func verifyToken(db *sql.DB, tokenString string) (*TokenRecord, error) {
	var record TokenRecord
	err := db.QueryRow(`
        SELECT id, user_id, token, expires_at, created_at
        FROM tokens
        WHERE token = $1 AND expires_at > NOW()`, tokenString).Scan(&record.ID, &record.UserID, &record.Token, &record.ExpiresAt, &record.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &TokenError{Message: ErrTokenNotFoundOrExpired, Type: TokenNotFoundError}
		}
		return nil, handleDatabaseError(err)
	}
	return &record, nil
}

type TokenErrorType int

const (
	TokenNotFoundError TokenErrorType = iota
	TokenDuplicateError
	UserNotFoundError
	DatabaseError
	UnknownError
)

type TokenError struct {
	Message string
	Type    TokenErrorType
}

func (e *TokenError) Error() string {
	return e.Message
}

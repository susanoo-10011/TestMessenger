package delete_posts

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"posts/src"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	authorized := r.Group("/")
	authorized.Use(src.AuthMiddleware(db))
	{
		authorized.DELETE("/posts/:id", DeletePostHandler(db))
	}

	return r
}

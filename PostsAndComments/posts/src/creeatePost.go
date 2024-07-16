package src

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func CreatePost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.ID = primitive.NewObjectID()
	post.CreatedAt = time.Now()

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

/*
func CreatePost(c *gin.Context) error {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return fmt.Errorf("failed to bind JSON: %w", err)
	}

	if post.Title == "" || post.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Content are required"})
		return fmt.Errorf("title and content are required")
	}

	post.ID = primitive.NewObjectID()
	post.CreatedAt = time.Now()

	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	collection := Client.Database(dbName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return fmt.Errorf("failed to create post: %w", err)
	}

	c.JSON(http.StatusCreated, post)
	return nil
}

*/

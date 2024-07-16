package src

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client

type Post struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	Image     *MediaFile         `bson:"image,omitempty"`
	Video     *MediaFile         `bson:"video,omitempty"`
	Gif       *MediaFile         `bson:"gif,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
}

type MediaFile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	URL       string             `bson:"url"`
	FileType  string             `bson:"file_type"`
	CreatedAt time.Time          `bson:"created_at"`
}

func StartServer() {
	r := gin.Default()
	r.POST("/posts", CreatePost)
	if err := r.Run(":9090"); err != nil {
		panic(err)
	}
}

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

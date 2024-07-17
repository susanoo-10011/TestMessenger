package src

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var Client *mongo.Client

type Post struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"Title"`
	Content   string             `bson:"Content"`
	Image     *MediaFile         `bson:"Image,omitempty"`
	Video     *MediaFile         `bson:"Video,omitempty"`
	Gif       *MediaFile         `bson:"Gif,omitempty"`
	CreatedAt time.Time          `bson:"Created_at"`
}

type MediaFile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	URL       string             `bson:"url"`
	FileType  string             `bson:"file_type"`
	CreatedAt time.Time          `bson:"created_at"`
}

type TokenRecord struct {
	ID        int64     `db:"id"`
	UserID    string    `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

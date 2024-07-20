package src

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id"`
	PostID    string             `bson:"post_id" json:"post_id"`
	UserID    int64              `bson:"user_id" json:"user_id"`
	Title     string             `bson:"Title" validate:"requireed, min=3, mac=100"`
	Content   string             `bson:"Content" validate:"required"`
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

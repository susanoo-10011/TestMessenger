package src

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID              primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	PostID          string             `bson:"post_id" json:"post_id,omitempty"`
	UserID          int64              `bson:"user_id" json:"user_id,omitempty"`
	Title           string             `bson:"Title" validate:"requireed, min=3, mac=100" json:"title,omitempty"`
	Content         string             `bson:"Content" validate:"required" json:"content,omitempty"`
	Image           *MediaFile         `bson:"Image,omitempty" json:"image,omitempty"`
	Video           *MediaFile         `bson:"Video,omitempty" json:"video,omitempty"`
	Gif             *MediaFile         `bson:"Gif,omitempty" json:"gif,omitempty"`
	CreatedAt       time.Time          `bson:"Created_at" json:"created_at"`
	LikesCount      int64              `bson:"likes_count" json:"likes_count,omitempty"`
	RepostCount     int64              `bson:"repost_count" json:"repost_count,omitempty"`
	CommentCount    int64              `bson:"comment_count" json:"comment_count,omitempty"`
	ViewsCount      int64              `bson:"viewsCount" json:"views_count,omitempty"`
	IncludeComments bool               `bson:"include_comments" json:"include_comments,omitempty"`
}

type MediaFile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	URL       string             `bson:"url"`
	FileType  string             `bson:"file_type"`
	CreatedAt time.Time          `bson:"created_at"`
}
type Comment struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	PostID          primitive.ObjectID `bson:"post_id"`
	UserID          int64              `bson:"user_id"`
	UserName        string             `bson:"user_name"`
	UserAvatar      string             `bson:"user_avatar"`
	ContentComments string             `bson:"content"`
	CreatedAt       time.Time          `bson:"created_at"`
	LikesCount      int64              `bson:"likes_count"`
	ParentCommentID primitive.ObjectID `bson:"parent_comment_id,omitempty"`
}
type Like struct {
	ID        primitive.ObjectID `bson:"_id"`
	PostID    string             `bson:"post_id"`
	UserID    int64              `bson:"user_id"`
	CreatedAt time.Time          `bson:"created_at"`
}

type Repost struct {
	ID        primitive.ObjectID `bson:"_id"`
	PostID    string             `bson:"post_id"`
	UserID    int64              `bson:"user_id"`
	CreatedAt time.Time          `bson:"created_at"`
}

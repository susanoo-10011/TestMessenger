package edit_posts

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"posts/src"
	"time"
)




func EditPost(postID string, updatedPost src.Post, userID int64) error {
	ctx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"post_id": postID, "usser_id": userID}

	update := bson.{
		"$set" := bson.M{
			"Title": updatedPost.Title,
			"Content": updatedPost.Content,
			"Image": updatedPost.Image,
			"Video": updatedPost.Video,
			"Gif": updatedPost.Gif,
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("post not found or user doesn't have permission to edit")
	}

	return nil

}

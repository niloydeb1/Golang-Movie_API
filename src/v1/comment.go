package v1

import (
	"context"
	"errors"
	"github.com/niloydeb1/Golang-Movie_API/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const CommentCollection = "commentCollection"

type Comment struct {
	ID             string    `json:"id" bson:"id"`
	MovieId        string    `json:"movie_id" bson:"movie_id"`
	ReviewId       string    `json:"review_id" bson:"review_id"`
	CommenterId    string    `json:"commenter_id" bson:"commenter_id"`
	CommenterEmail string    `json:"email" bson:"email"`
	Comment        string    `json:"comment" bson:"comment"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
}

func (c Comment) Validate() error {
	if c.ReviewId == "" {
		return errors.New("review id is not provided")
	}
	if c.Comment == "" {
		return errors.New("comment is not provided")
	}
	return nil
}

func (c Comment) GetByID(id string) Comment {
	query := bson.M{
		"$and": []bson.M{
			{"id": id},
		},
	}
	coll := config.GetDmManager().Db.Collection(CommentCollection)
	result := coll.FindOne(config.GetDmManager().Ctx, query, nil)
	res := new(Comment)
	err := result.Decode(res)
	if err != nil {
		log.Println("[ERROR]", err)
	}
	return *res
}

func (c Comment) GetByReviewId(reviewId string, pagination Pagination) ([]Comment, int64) {
	var data []Comment
	query := bson.M{
		"$and": []bson.M{
			{"review_id": reviewId},
		},
	}
	coll := config.GetDmManager().Db.Collection(CommentCollection)
	skip := pagination.Page * pagination.Limit
	findOptions := options.FindOptions{
		Limit: &pagination.Limit,
		Skip:  &skip,
		Sort:  bson.M{"created_at": 1},
	}
	result, err := coll.Find(config.GetDmManager().Ctx, query, &findOptions)
	if err != nil {
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(Comment)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		data = append(data, *elemValue)
	}
	count, err := coll.CountDocuments(config.GetDmManager().Ctx, query)
	if err != nil {
		log.Println(err.Error())
	}
	return data, count
}

func (c Comment) Store(comment Comment) error {
	coll := config.GetDmManager().Db.Collection(CommentCollection)
	_, err := coll.InsertOne(config.GetDmManager().Ctx, comment)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
		return err
	}
	return nil
}

func (c Comment) Delete(id string) error {
	coll := config.GetDmManager().Db.Collection(CommentCollection)
	filter := bson.M{"id": id}
	data, err := coll.DeleteOne(config.GetDmManager().Ctx, filter)
	if err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	if data.DeletedCount == 0 {
		log.Println("No data found to delete!")
		return errors.New("no data found to delete")
	}
	return nil
}

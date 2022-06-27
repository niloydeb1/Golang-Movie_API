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

const ReviewCollection = "reviewCollection"

type Review struct {
	ID            string        `json:"id" bson:"id"`
	Movie         ReviewedMovie `json:"movie" bson:"movie"`
	ReviewerEmail string        `json:"email" bson:"email"`
	ReviewerId    string        `json:"reviewer_id" bson:"reviewer_id"`
	ReviewTitle   string        `json:"review_title" bson:"review_title"`
	Description   string        `json:"description" bson:"description"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
}

type ReviewedMovie struct {
	ID       string `json:"id" bson:"id"`
	Title    string `json:"Title" bson:"Title"`
	Year     string `json:"Year" bson:"Year"`
	Genre    string `json:"Genre" bson:"Genre"`
	Director string `json:"Director" bson:"Director"`
}

func (r Review) Validate() error {
	if r.Movie.ID == "" {
		return errors.New("movie id is not provided")
	}
	if r.ReviewTitle == "" {
		return errors.New("review title is not provided")
	}
	if r.Description == "" {
		return errors.New("review description is not provided")
	}
	return nil
}

func (r Review) GetByID(id string) Review {
	query := bson.M{
		"$and": []bson.M{
			{"id": id},
		},
	}
	coll := config.GetDmManager().Db.Collection(ReviewCollection)
	result := coll.FindOne(config.GetDmManager().Ctx, query, nil)
	res := new(Review)
	err := result.Decode(res)
	if err != nil {
		log.Println("[ERROR]", err)
	}
	return *res
}

func (r Review) GetByMovieTitle(title string) []Review {
	var data []Review
	query := bson.M{
		"$and": []bson.M{
			{"movie.Title": title},
		},
	}
	coll := config.GetDmManager().Db.Collection(ReviewCollection)
	result, err := coll.Find(config.GetDmManager().Ctx, query, &options.FindOptions{Sort: bson.M{"created_at": -1}})
	if err != nil {
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(Review)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		data = append(data, *elemValue)
	}
	return data
}

func (r Review) Store(review Review) error {
	coll := config.GetDmManager().Db.Collection(ReviewCollection)
	_, err := coll.InsertOne(config.GetDmManager().Ctx, review)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
		return err
	}
	return nil
}

func (r Review) Search(query bson.M, pagination Pagination) ([]Review, int64) {
	var data []Review
	coll := config.GetDmManager().Db.Collection(ReviewCollection)
	skip := pagination.Page * pagination.Limit
	findOptions := options.FindOptions{
		Limit: &pagination.Limit,
		Skip:  &skip,
		Sort:  bson.M{"created_at": -1},
	}
	result, err := coll.Find(config.GetDmManager().Ctx, query, &findOptions)
	if err != nil {
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(Review)
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

func (r Review) Delete(id string) error {
	coll := config.GetDmManager().Db.Collection(ReviewCollection)
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

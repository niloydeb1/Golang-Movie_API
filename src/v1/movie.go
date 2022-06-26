package v1

import (
	"context"
	"github.com/niloydeb1/Golang-Movie_API/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const MovieCollection = "movieCollection"

type Movie struct {
	ID         string `json:"id" bson:"id"`
	Title      string `json:"Title" bson:"Title"`
	Year       string `json:"Year" bson:"Year"`
	Released   string `json:"Released" bson:"Released"`
	Runtime    string `json:"Runtime" bson:"Runtime"`
	Genre      string `json:"Genre" bson:"Genre"`
	Director   string `json:"Director" bson:"Director"`
	Writer     string `json:"Writer" bson:"Writer"`
	Actors     string `json:"Actors" bson:"Actors"`
	Plot       string `json:"Plot" bson:"Plot"`
	Language   string `json:"Language" bson:"Language"`
	Country    string `json:"Country" bson:"Country"`
	Awards     string `json:"Awards" bson:"Awards"`
	Poster     string `json:"Poster" bson:"Poster"`
	Metascore  string `json:"Metascore" bson:"Metascore"`
	ImdbRating string `json:"imdbRating" bson:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes" bson:"imdbVotes"`
	Type       string `json:"Type" bson:"Type"`
	BoxOffice  string `json:"BoxOffice" bson:"BoxOffice"`
	Website    string `json:"Website" bson:"Website"`
}

func (m Movie) GetByID(id string) Movie {
	query := bson.M{
		"$and": []bson.M{
			{"id": id},
		},
	}
	coll := config.GetDmManager().Db.Collection(MovieCollection)
	result := coll.FindOne(config.GetDmManager().Ctx, query, nil)
	res := new(Movie)
	err := result.Decode(res)
	if err != nil {
		log.Println("[ERROR]", err)
	}
	return *res
}

func (m Movie) GetByTitle(title string) Movie {
	query := bson.M{
		"$and": []bson.M{
			{"Title": title},
		},
	}
	coll := config.GetDmManager().Db.Collection(MovieCollection)
	result := coll.FindOne(config.GetDmManager().Ctx, query, nil)
	res := new(Movie)
	err := result.Decode(res)
	if err != nil {
		log.Println("[ERROR]", err)
	}
	return *res
}

func (m Movie) Store(movie Movie) error {
	coll := config.GetDmManager().Db.Collection(MovieCollection)
	_, err := coll.InsertOne(config.GetDmManager().Ctx, movie)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
		return err
	}
	return nil
}

func (m Movie) Search(query bson.M, pagination Pagination) ([]Movie, int64) {
	var data []Movie
	coll := config.GetDmManager().Db.Collection(MovieCollection)
	skip := pagination.Page * pagination.Limit
	findOptions := options.FindOptions{
		Limit: &pagination.Limit,
		Skip:  &skip,
	}
	result, err := coll.Find(config.GetDmManager().Ctx, query, &findOptions)
	if err != nil {
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(Movie)
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

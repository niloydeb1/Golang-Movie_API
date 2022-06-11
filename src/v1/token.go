package v1

import (
	"context"
	"errors"
	"github.com/niloydeb1/Golang-Movie_API/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type TokenService struct {
}

const TokenCollection = "tokenCollection"

func (t TokenService) GetByToken(token string) Token {
	var res Token
	query := bson.M{
		"$or": []interface{}{
			bson.M{"token": token},
			bson.M{"refresh_token": token},
		},
	}
	coll := config.GetDmManager().Db.Collection(TokenCollection)
	result, err := coll.Find(config.GetDmManager().Ctx, query, nil)
	if err != nil {
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(Token)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			return Token{}
		}
		res = *elemValue
	}
	return res
}

func (t TokenService) GetByUID(uid string) Token {
	var res Token
	query := bson.M{
		"$and": []bson.M{},
	}
	and := []bson.M{{"uid": uid}}
	query["$and"] = and
	coll := config.GetDmManager().Db.Collection(TokenCollection)
	result, err := coll.Find(config.GetDmManager().Ctx, query, nil)
	if err != nil {
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(Token)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			return Token{}
		}
		res = *elemValue
	}
	return res
}

func (t TokenService) Store(token Token) error {
	coll := config.GetDmManager().Db.Collection(TokenCollection)
	_, err := coll.InsertOne(config.GetDmManager().Ctx, token)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
	}
	return nil
}

func (t TokenService) Delete(uid string) error {
	coll := config.GetDmManager().Db.Collection(TokenCollection)
	filter := bson.M{"uid": uid}
	res, err := coll.DeleteOne(config.GetDmManager().Ctx, filter)
	if err != nil {
		log.Println("[ERROR]", err)
	}
	if res.DeletedCount == 0 {
		return errors.New("[ERROR] Delete failed")
	}
	return err
}

func (t TokenService) Update(token string, refreshToken string, existingToken string) error {
	oldTokenObj := t.GetByToken(existingToken)
	if oldTokenObj.Uid == "" {
		return errors.New("[ERROR] Token does not exists")
	}
	oldTokenObj.Token = token
	oldTokenObj.RefreshToken = refreshToken

	filter := bson.M{
		"$and": []interface{}{
			bson.M{"uid": oldTokenObj.Uid},
		},
	}
	update := bson.M{
		"$set": oldTokenObj,
	}
	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	coll := config.GetDmManager().Db.Collection(TokenCollection)
	err := coll.FindOneAndUpdate(config.GetDmManager().Ctx, filter, update, &opt)
	if err != nil {
		log.Println("[ERROR]", err.Err())
	}
	return nil
}
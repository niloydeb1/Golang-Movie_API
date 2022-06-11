package v1

import (
	"context"
	"github.com/niloydeb1/Golang-Movie_API/config"
	"github.com/niloydeb1/Golang-Movie_API/enums"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

const UserCollection = "userCollection"

type User struct {
	ID                 string                 `json:"id" bson:"id"`
	FirstName          string                 `json:"first_name" bson:"first_name" `
	LastName           string                 `json:"last_name" bson:"last_name"`
	Email              string                 `json:"email" bson:"email" `
	Phone              string                 `json:"phone" bson:"phone" `
	Password           string                 `json:"password" bson:"password" `
	Status             enums.STATUS           `json:"status" bson:"status"`
	CreatedDate        time.Time              `json:"created_date" bson:"created_date"`
	UpdatedDate        time.Time              `json:"updated_date" bson:"updated_date"`
	Role  			   enums.ROLE			  `json:"role" bson:"role"`
}

func (u User) GetUsers(status enums.STATUS) []User {
	var results []User
	query := bson.M{
		"$and": []bson.M{
			{"status": status},
		},
	}
	coll := config.GetDmManager().Db.Collection(UserCollection)
	result, err := coll.Find(config.GetDmManager().Ctx, query, nil)
	if err != nil {
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(User)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		results = append(results, *elemValue)
	}
	return results
}

func (u User) UpdateStatus(id string, status enums.STATUS) error {
	user := u.GetByID(id)
	user.Status = status
	filter := bson.M{
		"$and": []bson.M{
			{"id": id},
		},
	}
	update := bson.M{
		"$set": user,
	}
	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	coll := config.GetDmManager().Db.Collection(UserCollection)
	err := coll.FindOneAndUpdate(config.GetDmManager().Ctx, filter, update, &opt)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Err())
	}
	return nil
}

func (u User) UpdatePassword(user User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
	}
	user.Password = string(hashedPassword)
	filter := bson.M{
		"$and": []interface{}{
			bson.M{"id": user.ID},
		},
	}
	update := bson.M{
		"$set": user,
	}
	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	coll := config.GetDmManager().Db.Collection(UserCollection)
	uopdateErr := coll.FindOneAndUpdate(config.GetDmManager().Ctx, filter, update, &opt)
	if err != nil {
		log.Println("[ERROR]", uopdateErr.Err())
	}
	return nil
}

func (u User) GetByEmail(email string) User {
	var res User
	query := bson.M{
		"$and": []bson.M{},
	}
	and := []bson.M{
		{"email": email},
	}
	query["$and"] = and
	coll := config.GetDmManager().Db.Collection(UserCollection)
	result, err := coll.Find(config.GetDmManager().Ctx, query, nil)
	if err != nil {
		log.Println(err.Error())
		return User{}
	}
	for result.Next(context.TODO()) {
		elemValue := new(User)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		res = *elemValue
	}
	return res
}

func (u User) Store(user User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
	}
	user.Password = string(hashedPassword)
	coll := config.GetDmManager().Db.Collection(UserCollection)
	_, err = coll.InsertOne(config.GetDmManager().Ctx, user)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
	}
	return nil
}

func (u User) Get() []User {
	var results []User
	coll := config.GetDmManager().Db.Collection(UserCollection)
	result, err := coll.Find(config.GetDmManager().Ctx, bson.D{}, nil)
	if err != nil {
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(User)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		results = append(results, *elemValue)
	}
	return results
}

func (u User) GetByID(id string) User {
	var res User
	query := bson.M{
		"$and": []bson.M{
			{"id": id},
		},
	}
	coll := config.GetDmManager().Db.Collection(UserCollection)
	result, err := coll.Find(config.GetDmManager().Ctx, query, nil)
	if err != nil {
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(User)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		res = *elemValue
	}
	return res
}

func (u User) Delete(id string) error {
	user := u.GetByID(id)
	user.Status = enums.DELETED
	filter := bson.M{
		"$and": []bson.M{
			{"id": id},
		},
	}
	update := bson.M{
		"$set": user,
	}
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	coll := config.GetDmManager().Db.Collection(UserCollection)
	err := coll.FindOneAndUpdate(config.GetDmManager().Ctx, filter, update, &opt)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Err())
	}
	return nil
}
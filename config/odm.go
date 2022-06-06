package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

type dmManager struct {
	Ctx context.Context
	Db  *mongo.Database
}

var singletonDmManager *dmManager
var onceDmManager sync.Once

// GetDmManager returns dmManager
func GetDmManager() *dmManager {
	onceDmManager.Do(func() {
		log.Println("[INFO] Starting Initializing Singleton DB Manager")
		singletonDmManager = &dmManager{}
		singletonDmManager.initConnection()
	})
	return singletonDmManager
}

func (dm *dmManager) initConnection() {
	// Base context.
	log.Println(DatabaseConnectionString)
	ctx := context.Background()
	dm.Ctx = ctx
	clientOpts := options.Client().ApplyURI(DatabaseConnectionString)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Println("[ERROR] DB Connection error:", err.Error())
		return
	}

	db := client.Database(DatabaseName)
	dm.Db = db

	log.Println("[INFO] Initialized Singleton DB Manager")
}


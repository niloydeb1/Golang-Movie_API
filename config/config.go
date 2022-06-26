package config

import (
	"github.com/joho/godotenv"
	"github.com/niloydeb1/Golang-Movie_API/enums"
	"log"
	"os"
	"strings"
)

// FirstName admin first name.
var FirstName string

// LastName admin last name.
var LastName string

// Email admin email.
var Email string

// Password admin password.
var Password string

// PhoneNumber admin phone number.
var PhoneNumber string

// ServerPort refers to server port.
var ServerPort string

// DbServer refers to database server ip.
var DbServer string

// DbPort refers to database server port.
var DbPort string

// DbUsername refers to database name.
var DbUsername string

// DbPassword refers to database password.
var DbPassword string

// DatabaseConnectionString refers to database connection string.
var DatabaseConnectionString string

// DatabaseName refers to database name.
var DatabaseName string

// Database refers to database options.
var Database string

// PrivateKey refers to rsa private key .
var PrivateKey string

// Publickey refers to rsa public key .
var Publickey string

// TokenLifetime refers to token lifetime.
var TokenLifetime string

// RunMode refers to run mode.
var RunMode string

// EnableOpenTracing set true if opentracing is needed.
var EnableOpenTracing bool

// InitEnvironmentVariables initializes environment variables
func InitEnvironmentVariables() {
	RunMode = os.Getenv("RUN_MODE")
	if RunMode == "" {
		RunMode = string(enums.DEVELOP)
	}

	if RunMode != string(enums.PRODUCTION) {
		// Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Println("ERROR:", err.Error())
			return
		}
	}
	log.Println("RUN MODE:", RunMode)

	ServerPort = os.Getenv("SERVER_PORT")
	DbServer = os.Getenv("MONGO_SERVER")
	DbPort = os.Getenv("MONGO_PORT")
	DbUsername = os.Getenv("MONGO_USERNAME")
	DbPassword = os.Getenv("MONGO_PASSWORD")
	DatabaseName = os.Getenv("DATABASE_NAME")
	Database = os.Getenv("DATABASE")
	if Database == enums.MONGO {
		if DbUsername == "" && DbPassword == "" {
			DatabaseConnectionString = "mongodb://" + DbServer + ":" + DbPort
		} else {
			DatabaseConnectionString = "mongodb://" + DbUsername + ":" + DbPassword + "@" + DbServer + ":" + DbPort
		}
	}
	PrivateKey = os.Getenv("PRIVATE_KEY")
	Publickey = os.Getenv("PUBLIC_KEY")
	TokenLifetime = os.Getenv("TOKEN_LIFETIME")

	FirstName = os.Getenv("USER_FIRST_NAME")
	LastName = os.Getenv("USER_LAST_NAME")
	Email = os.Getenv("USER_EMAIL")
	Password = os.Getenv("USER_PASSWORD")
	PhoneNumber = os.Getenv("USER_PHONE")
	if os.Getenv("ENABLE_OPENTRACING") == "" {
		EnableOpenTracing = false
	} else {
		if strings.ToLower(os.Getenv("ENABLE_OPENTRACING")) == "true" {
			EnableOpenTracing = true
		} else {
			EnableOpenTracing = false
		}
	}
}

package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4/middleware"
	"github.com/niloydeb1/Golang-Movie_API/api"
	"github.com/niloydeb1/Golang-Movie_API/config"
	"github.com/niloydeb1/Golang-Movie_API/enums"
	v1 "github.com/niloydeb1/Golang-Movie_API/src/v1"
	"log"
	"net/http"
	"time"
)

func main() {
	e := config.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	initSuperAdmin()

	api.Routes(e)
	e.Logger.Fatal(e.Start(":" + config.ServerPort))
}

func initSuperAdmin() {
	if config.Email != "" {
		user := v1.User{}.GetByEmail(config.Email)
		if user.ID == "" {
			user = v1.User{
				ID:                 uuid.New().String(),
				FirstName:          config.FirstName,
				LastName:           config.LastName,
				Email:              config.Email,
				Phone:              config.PhoneNumber,
				Password:           config.Password,
				Status:             enums.ACTIVE,
				CreatedDate:        time.Now().UTC(),
				UpdatedDate:        time.Now().UTC(),
				Role:               enums.SUPERADMIN,
			}
			err := v1.User{}.Store(user)
			if err == nil {
				log.Println(err)
			}
		}
	}
}

//swag init --parseDependency --parseInternal
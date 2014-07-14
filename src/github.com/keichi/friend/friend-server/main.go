package main

import (
	"encoding/json"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

func main() {
	api := Api{}
	api.LoadConfig()
	api.InitDB()
	api.InitSchema()

	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
	}

	err := handler.SetRoutes(
		rest.RouteObjectMethod("GET", "/users/:name", &api, "GetUser"),
		rest.RouteObjectMethod("POST", "/users", &api, "CreateUser"),
		rest.RouteObjectMethod("DELETE", "/users/:name", &api, "DeleteUser"),
		rest.RouteObjectMethod("POST", "/login", &api, "LoginUser"),
		rest.RouteObjectMethod("POST", "/logout", &api, "LogoutUser"),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8080", &handler))
}

type Config struct {
	PasswordMinLength int
	HashStretchCount  int
	SessionKeyLength  int
	SessionExpiration int
	ProhibitedNames   []string
}

type Api struct {
	DB     gorm.DB
	Config Config
}

func (api *Api) LoadConfig() {
	api.Config = Config{
		PasswordMinLength: 8,
		HashStretchCount:  1024,
		SessionKeyLength:  32,
		SessionExpiration: 30,
		ProhibitedNames:   []string{},
	}

	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	if err = json.NewDecoder(file).Decode(&api.Config); err != nil {
		log.Fatal(err)
	}
}

func (api *Api) InitDB() {
	var err error
	api.DB, err = gorm.Open("sqlite3", "./friend.db")
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	api.DB.LogMode(true)
}

func (api *Api) InitSchema() {
	api.DB.AutoMigrate(User{})
	api.DB.AutoMigrate(TrustRelation{})
	api.DB.AutoMigrate(Session{})
}

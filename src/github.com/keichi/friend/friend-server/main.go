package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

func main() {
	api := Api{}
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

type User struct {
	Id        int64 `primaryKey:"yes"`
	Name      string
	Password  string
	PublicKey string
	Sessions  []Session
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TrustRelation struct {
	Id        int64 `primaryKey:"yes"`
	TrusterId int64
	TrusteeId int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	Id        int64 `primaryKey:"yes"`
	UserId    int64
	Token     string
	Expires   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Api struct {
	DB gorm.DB
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

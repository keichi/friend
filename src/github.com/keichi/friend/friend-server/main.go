package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strings"
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

func (api *Api) GetUser(w rest.ResponseWriter, r *rest.Request) {
	name := r.PathParam("name")
	token := r.Header.Get("X-Friend-Session-Token")
	user := User{}
	if api.DB.Where("name = ?", name).First(&user).RecordNotFound() {
		rest.Error(w, "User not found", 400)
		return
	}

	user.Password = ""
	if api.AuthenticateUser(name, token) {
		api.DB.Model(&user).Related(&user.Sessions)
	}

	w.WriteJson(&user)
}

func (api *Api) DeleteUser(w rest.ResponseWriter, r *rest.Request) {
}

func GetPasswordHash(name string, password string) (hash []byte) {
	hasher := sha256.New()
	hash = []byte{}

	for i := 0; i < 1000; i++ {
		hasher.Write(hash)
		hasher.Write([]byte(name))
		hasher.Write([]byte(password))
		hash = hasher.Sum(nil)
	}

	return
}

func (api *Api) CreateUser(w rest.ResponseWriter, r *rest.Request) {
	user := User{}
	r.DecodeJsonPayload(&user)

	if strings.TrimSpace(user.Name) == "" {
		rest.Error(w, "Username is empty", 400)
		return
	}
	if len(strings.TrimSpace(user.Password)) <= 8 {
		rest.Error(w, "Password is too short", 400)
		return
	}

	if api.DB.Where("name = ?", user.Name).First(&user).RecordNotFound() {
		user.Id = 0
		hash := GetPasswordHash(user.Name, user.Password)
		user.Password = hex.EncodeToString(hash)

		api.DB.Save(&user)

		user.Password = ""
		w.WriteJson(user)
		return
	}

	rest.Error(w, "User with the same name already exists", 400)
}

func (api *Api) LoginUser(w rest.ResponseWriter, r *rest.Request) {
	user := User{}
	r.DecodeJsonPayload(&user)

	if strings.TrimSpace(user.Name) == "" {
		rest.Error(w, "Username is empty", 400)
		return
	}
	if strings.TrimSpace(user.Password) == "" {
		rest.Error(w, "Password is empty", 400)
		return
	}

	dbUser := User{}
	if api.DB.Where("name = ?", user.Name).First(&dbUser).RecordNotFound() {
		rest.Error(w, "User not found", 400)
		return
	}

	if dbUser.Password != hex.EncodeToString(GetPasswordHash(user.Name, user.Password)) {
		rest.Error(w, "Password is wrong", 400)
		return
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		rest.Error(w, "Failed to generate session key", 500)
		return
	}
	token := hex.EncodeToString(buf)
	session := Session{
		Token:   token,
		Expires: time.Now().AddDate(0, 0, 30),
	}

	dbUser.Sessions = append(dbUser.Sessions, session)
	api.DB.Save(&dbUser)
	w.WriteJson(&session)
}

func (api *Api) AuthenticateUser(name string, token string) (succeeded bool) {
	user := User{}
	session := Session{}

	api.DB.Where("name = ?", name).First(&user)
	return !api.DB.Where("user_id = ? and token = ?", user.Id, token).First(&session).RecordNotFound()
}

func (api *Api) LogoutUser(w rest.ResponseWriter, r *rest.Request) {
	token := r.Header.Get("X-Friend-Session-Token")
	if api.DB.Where("token = ?", token).Delete(Session{}).RecordNotFound() {
		rest.Error(w, "Session token is not valid", 400)
	}
}

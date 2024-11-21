package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type User struct {
	Username string
	Password string
	Token    string
}

var users = map[string]User{
	"admin": User{
		Username: "admin",
		Password: "password",
	},
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/check-token", checkTokenHandler)

	http.ListenAndServe(":8080", mux)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	token := r.FormValue("token")

	user, ok := users[username]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if token != user.Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func checkTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")

	_, ok := users[token]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func generateToken(username string) string {
	token := fmt.Sprintf("%s-%d", username, time.Now().Unix())
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hash)
}

func init() {
	for username, user := range users {
		user.Token = generateToken(username)
	}
}

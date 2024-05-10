package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"Shawty/models"

	"github.com/dgrijalva/jwt-go"
)

var users = map[string]models.Credentials{
	"user1": {UserID: 1, Username: "user1", Password: "password1", Type: 1},
	"user2": {UserID: 2, Username: "user2", Password: "password2", Type: 2},
}

var jwtKey = []byte("your-secret-key")

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedCreds, ok := users[creds.Username]
	if !ok || expectedCreds.Password != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   creds.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

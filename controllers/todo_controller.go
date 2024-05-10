package controllers

import (
	"Shawty/models"
	"Shawty/services"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func ParseToken(tokenString string) (string, error) {
	claims := jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Subject, nil
}

var todoService = services.NewToDoService()

func ListLists(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username, err := ParseToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, ok := users[username]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var lists []*models.ToDoList
	if user.Type == 1 {
		lists = todoService.GetListsByUserID(user.UserID)
	} else if user.Type == 2 {
		lists = todoService.GetAllLists()
	}

	// Silinmemiş listeleri filtrele
	var nonDeletedLists []*models.ToDoList
	for _, list := range lists {
		if !list.Deleted {
			nonDeletedLists = append(nonDeletedLists, list)
		}
	}

	// İstek gövdesini oku ve ToDoList ID'sini al
	var reqBody struct {
		ID uint `json:"id,omitempty"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if reqBody.ID != 0 {
		// Eğer ID varsa, sadece bu ID'ye ait listeyi döndür
		list := todoService.GetListByID(reqBody.ID, user.UserID, user.Type)
		if list == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		nonDeletedLists = list
	}
	// Dönüştürülen JSON'u gönder
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nonDeletedLists)
}
func CreateList(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username, err := ParseToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, ok := users[username]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var reqBody struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	list := todoService.CreateList(user.UserID, reqBody.Title)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
func DeleteList(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username, err := ParseToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, ok := users[username]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// İstek gövdesini oku ve ToDoList ID'sini al
	var reqBody struct {
		ID uint `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ToDoList'i sil
	err = todoService.DeleteList(user.UserID, reqBody.ID, user.Type)
	if err != nil {
		if err.Error() == "list not found" {
			w.WriteHeader(http.StatusNotFound)
		} else if err.Error() == "unauthorized" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
func UpdateListTitle(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username, err := ParseToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, ok := users[username]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// İstek gövdesini oku ve ToDoList ID'sini ve yeni başlığı al
	var reqBody struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Başlık güncelle
	var updateErr error = todoService.UpdateListTitle(user.UserID, reqBody.ID, reqBody.Title, user.Type)
	if updateErr != nil {
		if updateErr.Error() == "list not found" {
			w.WriteHeader(http.StatusNotFound)
		} else if updateErr.Error() == "unauthorized" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
func GetItems(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username, err := ParseToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, ok := users[username]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// İstek gövdesini oku ve ToDoItem ID'sini ve list ID'sini al
	var reqBody struct {
		ID     uint `json:"id"`
		ListID uint `json:"list_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var items []*models.ToDoItem
	if reqBody.ID != 0 {
		items = todoService.GetItemsByTaskID(reqBody.ID, user.UserID, user.Type)
	} else if reqBody.ListID != 0 {
		items = todoService.GetItemsByListID(user.UserID, reqBody.ListID, user.Type)
	}

	// Dönüştürülen JSON'u gönder
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
func CreateTask(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username, err := ParseToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, ok := users[username]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// İstek gövdesini oku ve ToDoItem bilgilerini al
	var reqBody struct {
		ListID uint   `json:"list_id"`
		Task   string `json:"task"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ToDoItem oluştur
	item, err := todoService.CreateTask(user.UserID, reqBody.ListID, reqBody.Task, user.Type)
	if err != nil {
		if err.Error() == "list not found" {
			w.WriteHeader(http.StatusNotFound)
		} else if err.Error() == "unauthorized" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Başarı durumunda oluşturulan ToDoItem'i JSON olarak döndür
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

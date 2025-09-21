package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/zollidan/fasadowo/models"
	"github.com/zollidan/fasadowo/utils"
	"gorm.io/gorm"
)

type UserData struct {
	Name string `json:"name"`
	Surname string `json:"surname"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}


type AuthHandler struct {
	DB *gorm.DB
}

func ( h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user UserData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Error decoding json body.")
		return
	}
	log.Printf("Received user: %+v\n", user)

	err = gorm.G[models.User](h.DB).Create(context.Background(), &models.User{
		Name: user.Name,
		Surname: user.Surname,
		Phone: user.Phone,
		Email: user.Email,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Error decoding json body.")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
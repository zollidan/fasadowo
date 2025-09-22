package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zollidan/fasadowo/models"
	"github.com/zollidan/fasadowo/utils"
	"gorm.io/gorm"
)

type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterData struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	DB *gorm.DB
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userJson UserLoginData
	err := json.NewDecoder(r.Body).Decode(&userJson)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Error decoding json body.")
		return
	}
	existUser, err := gorm.G[models.User](h.DB).Where("email = ?", userJson.Email).First(context.Background())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.WriteError(w, http.StatusNotFound, "Invalid credentials!")
		return
	}
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal Server Error!")
		return
	}

	if existUser.Password != userJson.Password {
		utils.WriteError(w, http.StatusBadRequest, "Invalid credentials!")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": existUser.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// TODO: provide cfg with secret
	tokenString, err := token.SignedString([]byte("dev-secret"))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal Server Error!")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokenString); err != nil {
		log.Printf("Failed to write error respponse: %v", err)
	}

}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user UserRegisterData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Error decoding json body.")
		return
	}

	err = gorm.G[models.User](h.DB).Create(context.Background(), &models.User{
		Name:    user.Name,
		Surname: user.Surname,
		Phone:   user.Phone,
		Email:   user.Email,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Error decoding json body.")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

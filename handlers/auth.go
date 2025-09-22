package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/zollidan/fasadowo/models"
	"github.com/zollidan/fasadowo/utils"
	"golang.org/x/crypto/bcrypt"
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
	DB        *gorm.DB
	TokenAuth *jwtauth.JWTAuth
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

	err = bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(userJson.Password))
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Invalid credentials!")
		return
	}

	_, tokenString, _ := h.TokenAuth.Encode(map[string]interface{}{
		"user_id": existUser.ID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour).Unix(),
		"role":    existUser.Role,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user UserRegisterData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Error decoding json body.")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	err = gorm.G[models.User](h.DB).Create(context.Background(), &models.User{
		Name:     user.Name,
		Surname:  user.Surname,
		Phone:    user.Phone,
		Email:    user.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Error decoding json body.")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

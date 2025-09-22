package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type ErrorResponse struct {
	Error   string
	Message string
}

type SuccessResponse struct {
	Message string
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to write error respponse: %v", err)
	}
}

// func writeSuccess(w http.ResponseWriter, statusCode int, message string) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	resp := SuccessResponse{
// 		Message: message,
// 	}
// 	if err := json.NewEncoder(w).Encode(resp); err != nil {
// 		log.Printf("Failed to write error respponse: %v", err)
// 	}
// }

// func GetParam(param string, w http.ResponseWriter, r *http.Request) (*int, error) {
// 	idStr := chi.URLParam(r, "subcategoryID")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		WriteError(w, http.StatusBadRequest, "Invalid subcategory ID")
// 		return nil, err
// 	}

// 	return id, nil
// }

// Обобщённый метод выборки по ID
func GetByID[T any](w http.ResponseWriter, db *gorm.DB, id int, preload ...string) (*T, error) {
	var entity T
	query := db.Model(&entity)

	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			WriteError(w, http.StatusNotFound, "record not found")
			return nil, err
		}
		WriteError(w, http.StatusInternalServerError, "db error")
		return nil, err
	}

	return &entity, nil
}

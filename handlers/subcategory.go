package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/zollidan/fasadowo/models"
	"github.com/zollidan/fasadowo/utils"
	"gorm.io/gorm"
)

type SubcategoryHandler struct {
	DB *gorm.DB
}
func (h *SubcategoryHandler) ListSubcategory(w http.ResponseWriter, r *http.Request) {
	var subcategories []models.Subcategory
	if err := h.DB.Find(&subcategories).Error; err != nil {
		http.Error(w, "failed to fetch subcategory", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(subcategories)
	if err != nil {
		http.Error(w, "failed to serialize subcategories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *SubcategoryHandler) GetSubcategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "subcategoryID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid subcategory ID")
		return
	}

	var subcategory models.Subcategory
	subcategory, err = gorm.G[models.Subcategory](h.DB).Preload("Collections", nil).Where("id = ?", id).First(context.Background())
	if err == gorm.ErrRecordNotFound {
		utils.WriteError(w, http.StatusNotFound, "subcategory not found")
		return
	} 
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error fetchig subcategory")
		return
	}
	
	json.NewEncoder(w).Encode(subcategory)
}
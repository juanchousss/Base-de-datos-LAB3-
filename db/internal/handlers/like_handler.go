package handlers

import (
	"encoding/json"
	"net/http"

	"db/internal/middleware"
	"db/internal/repository"
)

type LikeHandler struct {
	Repo *repository.LikeRepository
}

// NewLikeHandler es el constructor
func NewLikeHandler(repo *repository.LikeRepository) *LikeHandler {
	return &LikeHandler{Repo: repo}
}

// AddLike maneja la petición POST /api/likes
func (h *LikeHandler) AddLike(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Acceso no autorizado", http.StatusUnauthorized)
		return
	}

	var req struct {
		IDPost int `json:"id_post"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IDPost == 0 {
		http.Error(w, "JSON inválido o falta id_post", http.StatusBadRequest)
		return
	}

	if err := h.Repo.AddLike(userID, req.IDPost); err != nil {
		http.Error(w, "Error interno al dar like", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Like agregado"}`))
}

// RemoveLike maneja la petición DELETE /api/likes
func (h *LikeHandler) RemoveLike(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Acceso no autorizado", http.StatusUnauthorized)
		return
	}

	var req struct {
		IDPost int `json:"id_post"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IDPost == 0 {
		http.Error(w, "JSON inválido o falta id_post", http.StatusBadRequest)
		return
	}

	if err := h.Repo.RemoveLike(userID, req.IDPost); err != nil {
		http.Error(w, "Error interno al quitar like", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Like removido"}`))
}

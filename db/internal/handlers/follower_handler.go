package handlers

import (
	"encoding/json"
	"net/http"

	"db/internal/middleware"
	"db/internal/repository"
)

type FollowerHandler struct {
	Repo *repository.FollowerRepository
}

func NewFollowerHandler(repo *repository.FollowerRepository) *FollowerHandler {
	return &FollowerHandler{Repo: repo}
}

func (h *FollowerHandler) Follow(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del usuario que hace la acción (el seguidor)
	idFollower, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Acceso no autorizado", http.StatusUnauthorized)
		return
	}

	// 2. Leer a quién quiere seguir desde el JSON
	var req struct {
		IDFollow int `json:"id_follow"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IDFollow == 0 {
		http.Error(w, "JSON inválido o falta id_follow", http.StatusBadRequest)
		return
	}

	// Regla de negocio: No puedes seguirte a ti mismo
	if idFollower == req.IDFollow {
		http.Error(w, "No puedes seguirte a ti mismo", http.StatusBadRequest)
		return
	}

	// 3. Guardar en la Base de Datos
	err := h.Repo.FollowUser(idFollower, req.IDFollow)
	if err != nil {
		http.Error(w, "Error interno al seguir al usuario", http.StatusInternalServerError)
		return
	}

	// 4. Éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Usuario seguido con éxito"}`))
}

// Unfollow maneja la petición DELETE para dejar de seguir a un usuario
func (h *FollowerHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer ID del usuario que hace la acción
	idFollower, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Acceso no autorizado", http.StatusUnauthorized)
		return
	}

	// 2. Leer a quién quiere dejar de seguir desde el JSON
	var req struct {
		IDFollow int `json:"id_follow"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IDFollow == 0 {
		http.Error(w, "JSON inválido o falta id_follow", http.StatusBadRequest)
		return
	}

	// 3. Ejecutar el borrado en la Base de Datos
	err := h.Repo.UnfollowUser(idFollower, req.IDFollow)
	if err != nil {
		http.Error(w, "Error interno al dejar de seguir al usuario", http.StatusInternalServerError)
		return
	}

	// 4. Éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Has dejado de seguir al usuario"}`))
}

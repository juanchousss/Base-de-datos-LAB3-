package handlers

import (
	"encoding/json"
	"net/http"

	"db/internal/models"
	"db/internal/repository"
)

type SearchHandler struct {
	UserRepo *repository.UserRepository
	PostRepo *repository.PostRepository
}

// NewSearchHandler es el constructor
func NewSearchHandler(userRepo *repository.UserRepository, postRepo *repository.PostRepository) *SearchHandler {
	return &SearchHandler{
		UserRepo: userRepo,
		PostRepo: postRepo,
	}
}

// GlobalSearch ejecuta la búsqueda en paralelo para usuarios y posts
func (h *SearchHandler) GlobalSearch(w http.ResponseWriter, r *http.Request) {
	termino := r.URL.Query().Get("q")

	if termino == "" {
		http.Error(w, "Falta el parámetro 'q' de búsqueda en la URL", http.StatusBadRequest)
		return
	}

	// 1. Buscar Usuarios
	usuarios, err := h.UserRepo.SearchUsers(termino)
	if err != nil {
		http.Error(w, "Error al buscar usuarios", http.StatusInternalServerError)
		return
	}

	// 2. Buscar Posts
	posts, err := h.PostRepo.SearchPosts(termino)
	if err != nil {
		http.Error(w, "Error al buscar posts", http.StatusInternalServerError)
		return
	}

	// 3. Empaquetar y enviar
	respuesta := models.SearchResponse{
		Users: usuarios,
		Posts: posts,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuesta)
}

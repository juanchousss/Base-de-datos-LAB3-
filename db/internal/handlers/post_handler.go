package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"db/internal/middleware"
	"db/internal/models"
	"db/internal/repository"
)

// PostHandler encapsula el repositorio
type PostHandler struct {
	Repo *repository.PostRepository
}

func NewPostHandler(repo *repository.PostRepository) *PostHandler {
	return &PostHandler{Repo: repo}
}

// CreatePost maneja la petición POST /api/posts
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// 1. Extraemos el ID del usuario real desde el Contexto (validado por el JWT)
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Acceso no autorizado", http.StatusUnauthorized)
		return
	}

	var post models.Post

	// Decodificamos el JSON que envía la app móvil
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 2. BLINDAJE DE SEGURIDAD:
	// Sobrescribimos el IDUser con el ID real del token.
	// Así evitamos que un usuario malicioso mande el ID de otra persona en el JSON.
	post.IDUser = userID

	// Llamamos a nuestro repositorio para guardar en la BD
	id, err := h.Repo.CreatePost(&post)
	if err != nil {
		http.Error(w, "Error interno al crear el post", http.StatusInternalServerError)
		return
	}

	post.ID = id // Asignamos el ID generado por PostgreSQL

	// Devolvemos el post creado en formato JSON con código 201 (Created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// DeletePost maneja la petición DELETE /api/posts/{id}
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Acceso no autorizado", http.StatusUnauthorized)
		return
	}

	// Go 1.22+: Extraemos el {id} directamente de la ruta de la URL
	idPostStr := r.PathValue("id")
	idPost, err := strconv.Atoi(idPostStr)
	if err != nil {
		http.Error(w, "ID de post inválido en la URL", http.StatusBadRequest)
		return
	}

	err = h.Repo.DeletePost(idPost, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post no encontrado o no tienes permiso para borrarlo", http.StatusForbidden)
			return
		}
		http.Error(w, "Error interno al borrar el post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Post eliminado con éxito"}`))
}

// GetFeed maneja la petición GET /api/feed?cursor=0&limit=20
func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	// 1. REEMPLAZO: Ya no simulamos con el Header. Extraemos el ID real del Contexto
	idUser, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Acceso no autorizado", http.StatusUnauthorized)
		return
	}

	// Leemos los parámetros de paginación de la URL
	cursorStr := r.URL.Query().Get("cursor")
	cursor, _ := strconv.Atoi(cursorStr)

	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 20 // Límite por defecto para móviles
	}

	feed, err := h.Repo.GetFeed(idUser, cursor, limit)
	if err != nil {
		http.Error(w, "Error al obtener el feed", http.StatusInternalServerError)
		return
	}

	// Si no hay posts, devolvemos un arreglo vacío en lugar de null
	if feed == nil {
		feed = []models.Post{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)

}

// GetThread maneja la petición GET /api/posts/{id}/replies
func (h *PostHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	idPost, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de post inválido", http.StatusBadRequest)
		return
	}

	replies, err := h.Repo.GetReplies(idPost)
	if err != nil {
		http.Error(w, "Error interno al obtener el hilo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(replies)
}

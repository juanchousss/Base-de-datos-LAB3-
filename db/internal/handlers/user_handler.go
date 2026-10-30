package handlers

import (
	"database/sql"
	"db/internal/models"
	"db/internal/repository"
	"db/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

// Register recibe la petición POST desde el teléfono para crear una cuenta
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Leer el JSON que manda la app
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validar que no vengan campos vacíos
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	// Mandar a guardar a la base de datos
	id, err := h.Repo.RegisterUser(&user)
	if err != nil {
		http.Error(w, "Error al crear el usuario en BD", http.StatusInternalServerError)
		return
	}

	// Ocultar el password antes de responder
	user.Password = ""
	user.ID = id

	// Responder a la app
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login procesa la autenticación y devuelve un JWT
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// 1. Leer el JSON entrante
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 2. Buscar al usuario por correo en la Base de Datos
	user, err := h.Repo.GetUserByEmail(creds.Email)
	if err != nil {

		log.Printf("Error buscando usuario en DB: %v", err)

		// Por seguridad, si el correo no existe, no damos pistas. Error genérico.
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}

	// 3. Comparar la contraseña enviada con el Hash guardado
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {

		log.Printf("Error de contraseña (Bcrypt): %v", err)

		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}

	// 4. Si la contraseña es correcta, generar el Token JWT
	tokenString, err := services.GenerateToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Error al generar el token de acceso", http.StatusInternalServerError)
		return
	}

	// 5. Responder a la app móvil enviando el token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// GetProfile maneja la petición GET /api/users/{id}
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Extraemos el {id} directamente de la ruta de la URL
	idStr := r.PathValue("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Llamamos al repositorio
	profile, err := h.Repo.GetUserProfile(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Retornamos el perfil completo en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

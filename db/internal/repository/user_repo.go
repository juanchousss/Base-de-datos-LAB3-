package repository

import (
	"database/sql"
	"db/internal/models"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// RegistrarUsuario encripta la contraseña y guarda el registro en PostgreSQL
func (r *UserRepository) RegisterUser(user *models.User) (int, error) {
	// 1. Encriptar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error al encriptar: %v", err)
		return 0, err
	}

	// 2. Insertar en la BD usando el hash
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id_user, created_at`

	err = r.DB.QueryRow(query, user.Username, user.Email, string(hashedPassword)).
		Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		log.Printf("Error en BD: %v", err)
		return 0, err
	}

	return user.ID, nil
}

// GetUserByEmail busca a un usuario en la base de datos por su correo electrónico para el Login
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	var hash string

	// Fíjate que ya estamos usando el "created_at" que arreglamos en la base de datos
	query := `
		SELECT id_user, username, email, password_hash, created_at 
		FROM users 
		WHERE email = $1`

	err := r.DB.QueryRow(query, email).
		Scan(&user.ID, &user.Username, &user.Email, &hash, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Guardamos temporalmente el hash en el struct para validar el login en el handler
	user.Password = hash

	return &user, nil
}

// GetUserProfile obtiene los datos, estadísticas y publicaciones de un usuario
func (r *UserRepository) GetUserProfile(userID int) (*models.UserProfile, error) {
	var profile models.UserProfile

	// 1. Obtener datos del usuario y contar seguidores/seguidos usando Subconsultas SQL
	queryUser := `
		SELECT 
			id_user, 
			email,
			(SELECT COUNT(*) FROM followers WHERE id_follow = $1) AS followers_count,
			(SELECT COUNT(*) FROM followers WHERE id_follower = $1) AS following_count
		FROM users 
		WHERE id_user = $1`

	err := r.DB.QueryRow(queryUser, userID).Scan(
		&profile.IDUser,
		&profile.Email,
		&profile.FollowersCount,
		&profile.FollowingCount,
	)
	if err != nil {
		return nil, err // Devolverá sql.ErrNoRows si el usuario no existe
	}

	// 2. Obtener el "Muro" (solo las publicaciones de este usuario específico)
	queryPosts := `
		SELECT id_post, id_user, content, id_reposts_of, created_at
		FROM posts
		WHERE id_user = $1
		ORDER BY created_at DESC`

	rows, err := r.DB.Query(queryPosts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profile.Posts = []models.Post{} // Inicializar para que el JSON devuelva [] y no null
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.IDUser, &p.Content, &p.IDRepostsOf, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		profile.Posts = append(profile.Posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &profile, nil

}

// SearchUsers busca usuarios por username ignorando mayúsculas y minúsculas
func (r *UserRepository) SearchUsers(termino string) ([]models.UserSearch, error) {
	// Limpiamos el término por si el usuario escribió "@usuario"
	terminoLimpio := strings.TrimPrefix(termino, "@")
	busqueda := "%" + terminoLimpio + "%"

	query := `
		SELECT id_user, username
		FROM users
		WHERE username ILIKE $1
		LIMIT 10`

	rows, err := r.DB.Query(query, busqueda)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []models.UserSearch{}
	for rows.Next() {
		var u models.UserSearch
		if err := rows.Scan(&u.IDUser, &u.Username); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usuarios, nil
}

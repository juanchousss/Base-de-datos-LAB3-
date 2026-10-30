package models

import "time"

// User mapea la tabla 'users'
type User struct {
	ID       int    `json:"id_user"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// Omitimos el password_hash en el JSON por seguridad
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

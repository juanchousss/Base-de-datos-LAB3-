package models

import "time"

// Post mapea la tabla 'posts' y se prepara para ser enviado como JSON
type Post struct {
	ID     int `json:"id_post"`
	IDUser int `json:"id_user"`
	// Usamos punteros para los campos que pueden ser NULL en la base de datos
	Content     *string   `json:"content,omitempty"`
	IDRepostsOf *int      `json:"id_reposts_of,omitempty"`
	IDReplyTo   *int      `json:"id_reply_to,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

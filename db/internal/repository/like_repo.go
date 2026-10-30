package repository

import (
	"database/sql"
	"log"
)

type LikeRepository struct {
	DB *sql.DB
}

// NewLikeRepository es el constructor
func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{DB: db}
}

// AddLike inserta el "me gusta" en la base de datos
func (r *LikeRepository) AddLike(idUser, idPost int) error {
	// ON CONFLICT DO NOTHING evita que la base de datos lance error si el usuario
	// manda la petición dos veces por accidente (requiere que id_user e id_post sean Primary Key compuesta)
	query := `
		INSERT INTO likes (id_user, id_post)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`

	_, err := r.DB.Exec(query, idUser, idPost)
	if err != nil {
		log.Printf("Error al dar like en DB: %v", err)
		return err
	}
	return nil
}

// RemoveLike elimina el "me gusta" de la base de datos
func (r *LikeRepository) RemoveLike(idUser, idPost int) error {
	query := `
		DELETE FROM likes 
		WHERE id_user = $1 AND id_post = $2`

	_, err := r.DB.Exec(query, idUser, idPost)
	if err != nil {
		log.Printf("Error al quitar like en DB: %v", err)
		return err
	}
	return nil
}

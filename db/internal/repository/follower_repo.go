package repository

import (
	"database/sql"
	"log"
)

type FollowerRepository struct {
	DB *sql.DB
}

// NewFollowerRepository es un constructor para inicializar el repositorio
func NewFollowerRepository(db *sql.DB) *FollowerRepository {
	return &FollowerRepository{DB: db}
}

// FollowUser inserta la relación entre el seguidor y el seguido
func (r *FollowerRepository) FollowUser(idFollower, idFollow int) error {
	query := `
		INSERT INTO followers (id_follower, id_follow)
		VALUES ($1, $2)`

	_, err := r.DB.Exec(query, idFollower, idFollow)
	if err != nil {
		log.Printf("Error al insertar seguidor en DB: %v", err)
		return err
	}

	return nil
}

// UnfollowUser elimina la relación entre el seguidor y el seguido
func (r *FollowerRepository) UnfollowUser(idFollower, idFollow int) error {
	query := `
		DELETE FROM followers 
		WHERE id_follower = $1 AND id_follow = $2`

	_, err := r.DB.Exec(query, idFollower, idFollow)
	if err != nil {
		log.Printf("Error al dejar de seguir en DB: %v", err)
		return err
	}

	return nil
}

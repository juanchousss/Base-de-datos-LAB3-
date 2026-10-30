package repository

import (
	"database/sql"
	"db/internal/models"
	"log"
)

// PostRepository maneja todas las operaciones de BD para los posts
type PostRepository struct {
	DB *sql.DB
}

// NewPostRepository es un constructor
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// 1. CrearPost inserta un nuevo post en la BD y retorna su ID
func (r *PostRepository) CreatePost(post *models.Post) (int, error) {
	query := `
		INSERT INTO posts (id_user, content, id_reposts_of, id_reply_to)
		VALUES ($1, $2, $3, $4)
		RETURNING id_post, created_at`

	err := r.DB.QueryRow(query, post.IDUser, post.Content, post.IDRepostsOf, post.IDReplyTo).
		Scan(&post.ID, &post.CreatedAt)

	if err != nil {
		log.Printf("Error insertando post: %v", err)
		return 0, err
	}

	return post.ID, nil
}

// 2. ObtenerFeed devuelve los posts de los usuarios que sigue alguien, usando Paginación por Cursor
func (r *PostRepository) GetFeed(idUser int, cursorIdPost int, limit int) ([]models.Post, error) {
	// Query altamente optimizado usando JOIN con la tabla seguidores
	// Buscamos posts cuyo ID sea MENOR al cursor para ir hacia atrás en el tiempo
	query := `
		SELECT p.id_post, p.id_user, p.content, p.id_reposts_of, p.id_reply_to, p.created_at 
		FROM posts p
		INNER JOIN followers s ON p.id_user = s.id_follow
		WHERE s.id_follower = $1 AND p.id_post < $2
		ORDER BY p.id_post DESC
		LIMIT $3`

	// Si es la primera vez que carga (no hay cursor), usamos un número gigante
	if cursorIdPost == 0 {
		cursorIdPost = 999999999
	}

	rows, err := r.DB.Query(query, idUser, cursorIdPost, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []models.Post

	// Iteramos sobre los resultados
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.IDUser, &p.Content, &p.IDRepostsOf, &p.IDReplyTo, &p.CreatedAt)
		if err != nil {
			log.Printf("Error escaneando fila: %v", err)
			continue
		}
		feed = append(feed, p)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error de iteración en la base de datos: %v", err)
		return nil, err
	}

	return feed, nil
}

// DeletePost elimina un post (o repost). Exigmos el idUser como medida de seguridad extrema.
func (r *PostRepository) DeletePost(idPost, idUser int) error {
	// Solo borra si el post existe Y además pertenece al usuario que hace la petición
	query := `DELETE FROM posts WHERE id_post = $1 AND id_user = $2`

	res, err := r.DB.Exec(query, idPost, idUser)
	if err != nil {
		return err
	}

	// Comprobamos si realmente se eliminó alguna fila
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // Usamos este error estándar para indicar "No encontrado o no autorizado"
	}

	return nil
}

// SearchPosts busca coincidencias de texto dentro del contenido de las publicaciones
func (r *PostRepository) SearchPosts(termino string) ([]models.Post, error) {
	busqueda := "%" + termino + "%"

	query := `
		SELECT id_post, id_user, content, id_reposts_of, id_reply_to, created_at
		FROM posts
		WHERE content ILIKE $1
		ORDER BY created_at DESC
		LIMIT 20`

	rows, err := r.DB.Query(query, busqueda)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.IDUser, &p.Content, &p.IDRepostsOf, &p.IDReplyTo, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetReplies obtiene todos los posts que son respuesta a un post específico (El Hilo)
func (r *PostRepository) GetReplies(idPost int) ([]models.Post, error) {
	// Buscamos los posts cuyo "padre" sea el ID que estamos consultando.
	// Usamos ASC (Ascendente) para que los comentarios más viejos salgan arriba, como en un chat.
	query := `
		SELECT id_post, id_user, content, id_reposts_of, id_reply_to, created_at
		FROM posts
		WHERE id_reply_to = $1
		ORDER BY created_at ASC`

	rows, err := r.DB.Query(query, idPost)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	replies := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.IDUser, &p.Content, &p.IDRepostsOf, &p.IDReplyTo, &p.CreatedAt); err != nil {
			return nil, err
		}
		replies = append(replies, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return replies, nil
}

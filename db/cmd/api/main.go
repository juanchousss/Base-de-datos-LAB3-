package main

import (
	"log"
	"net/http"

	"db/internal/config"
	"db/internal/handlers"
	"db/internal/repository"
	"db/internal/routes"
)

func main() {
	// 1. Inicializar la conexión a PostgreSQL
	// Reemplaza "usuario", "password" y "dbname" con los datos reales de tu BD local
	db := config.InitDB("postgres", "Siza22**", "twitter_clone", "localhost", "5432")
	defer db.Close()

	// 2. Inicializar los Repositorios (Capa de Base de Datos)
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	followerRepo := repository.NewFollowerRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	// 3. Inicializar los Handlers (Capa HTTP)
	userHandler := handlers.NewUserHandler(userRepo)
	postHandler := handlers.NewPostHandler(postRepo)
	followerHandler := handlers.NewFollowerHandler(followerRepo)
	likeHandler := handlers.NewLikeHandler(likeRepo)
	searchHandler := handlers.NewSearchHandler(userRepo, postRepo)
	// 4. Inicializar el Enrutador
	mux := routes.ConfigureRoutes(userHandler, postHandler, followerHandler, likeHandler, searchHandler)

	// 5. Levantar el Servidor HTTP
	log.Println("Servidor API corriendo en http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Fallo crítico en el servidor: %v", err)
	}
}

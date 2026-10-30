package routes

import (
	"net/http"

	"db/internal/handlers"
	"db/internal/middleware"
)

// ConfigureRoutes inicializa el multiplexor (enrutador) con todos los endpoints
func ConfigureRoutes(
	userHandler *handlers.UserHandler,
	postHandler *handlers.PostHandler,
	followerHandler *handlers.FollowerHandler,
	likeHandler *handlers.LikeHandler,
	searchHandler *handlers.SearchHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	// --- Rutas Públicas (No requieren Token) ---
	mux.HandleFunc("POST /api/users/register", userHandler.Register)
	mux.HandleFunc("POST /api/users/login", userHandler.Login)

	// --- Rutas Privadas (Protegidas por el Middleware) ---
	mux.HandleFunc("POST /api/posts", middleware.Authentication(postHandler.CreatePost))
	mux.HandleFunc("GET /api/feed", middleware.Authentication(postHandler.GetFeed))
	mux.HandleFunc("POST /api/users/follow", middleware.Authentication(followerHandler.Follow))
	mux.HandleFunc("DELETE /api/users/follow", middleware.Authentication(followerHandler.Unfollow))
	mux.HandleFunc("POST /api/likes", middleware.Authentication(likeHandler.AddLike))
	mux.HandleFunc("DELETE /api/likes", middleware.Authentication(likeHandler.RemoveLike))
	mux.HandleFunc("GET /api/search", middleware.Authentication(searchHandler.GlobalSearch))
	mux.HandleFunc("DELETE /api/posts/{id}", middleware.Authentication(postHandler.DeletePost))
	mux.HandleFunc("GET /api/users/{id}", middleware.Authentication(userHandler.GetProfile))
	mux.HandleFunc("GET /api/posts/{id}/replies", middleware.Authentication(postHandler.GetThread))
	return mux
}

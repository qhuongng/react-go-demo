package server

import (
	"net/http"

	"chi-mysql-boilerplate/internal/controllers"
	"chi-mysql-boilerplate/internal/server/middleware"
	"chi-mysql-boilerplate/internal/utils/helpers"
	"chi-mysql-boilerplate/internal/utils/validators"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	validator := validators.NewValidator(helpers.MessageLogs)

	postHandler := controllers.NewPostHandler(s.db, validator)
	authHandler := controllers.NewAuthHandler(s.db, validator)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.Cors())

	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Get("/posts", postHandler.GetAllPosts)
		v1.Get("/posts/{id}", postHandler.GetPostById)

		v1.Post("/auth/register", authHandler.Register)
		v1.Post("/auth/login", authHandler.Login)

		// protected routes
		v1.Group(func(v1 chi.Router) {
			v1.Use(middleware.VerifyAccessToken)
			v1.Get("/posts/by-user/{userId}", postHandler.GetPostsByUserId)
			v1.Post("/posts", postHandler.CreatePost)
			v1.Put("/posts/{id}", postHandler.UpdatePostById)
			v1.Delete("/posts/{id}", postHandler.DeletePostById)
		})

		// routes that need the refresh token
		v1.Group(func(v1 chi.Router) {
			v1.Use(middleware.ExtractRefreshToken)
			v1.Post("/auth/refresh", authHandler.RefreshAccessToken)
			v1.Post("/auth/logout", authHandler.Logout)
		})
	})

	return r
}

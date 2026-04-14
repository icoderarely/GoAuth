package router

import (
	"net/http"

	"github.com/icoderarely/GoAuth/internal/domain"
	"github.com/icoderarely/GoAuth/internal/handler"
	"github.com/icoderarely/GoAuth/internal/middleware"
	"github.com/icoderarely/GoAuth/internal/service"
)

func NewRouter(authSvc service.AuthService) http.Handler {
	mux := http.NewServeMux()

	// handler dependecies
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler()

	// middleware
	authMw := middleware.AuthMiddleware(authSvc)
	adminMw := middleware.RequireRole(domain.RoleAdmin)

	// auth routes
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	// protected routes
	mux.Handle("GET /me", authMw(http.HandlerFunc(userHandler.Me)))
	mux.Handle("GET /dashboard", authMw(http.HandlerFunc(userHandler.Dashboard)))

	// admin route
	mux.Handle("GET /admin", authMw(adminMw(http.HandlerFunc(userHandler.Admin))))

	return mux
}

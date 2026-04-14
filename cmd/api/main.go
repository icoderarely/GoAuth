package main

import (
	"context"
	"log"
	"net/http"

	"github.com/icoderarely/GoAuth/config"
	"github.com/icoderarely/GoAuth/internal/domain"
	"github.com/icoderarely/GoAuth/internal/repository/inmemory"
	"github.com/icoderarely/GoAuth/internal/router"
	"github.com/icoderarely/GoAuth/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.LoadConfig()

	repo := inmemory.NewStore(nil)

	authSvc := service.NewAuthService(
		repo,
		cfg.JWTSecret,
		cfg.TokenTTLHours,
	)

	ctx := context.Background()

	admin, err := authSvc.Register(ctx, "admin", "pass")
	if err != nil {
		log.Fatal("failed to create admin:", err)
	}

	if err := repo.SetRole(ctx, admin.ID, domain.RoleAdmin); err != nil {
		log.Fatal("failed to set admin role:", err)
	}

	r := router.NewRouter(authSvc)
	http.ListenAndServe(":"+cfg.Port, r)
}

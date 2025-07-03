package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/oThinas/bid/internal/api"
	"github.com/oThinas/bid/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	api := api.Api{
		Router:      chi.NewMux(),
		UserService: services.NewUserService(pool),
	}

	api.BindRoutes()

	fmt.Println("Server started on port :8080")
	if err := http.ListenAndServe("localhost:8080", api.Router); err != nil {
		panic(err)
	}
}

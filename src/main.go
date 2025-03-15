package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github-repo-tracker/src/controller"
	"github-repo-tracker/src/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := initDB(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close(ctx)

	repoRepository := repository.NewPostgresRepoRepository(db)

	r := gin.Default()

	controller.RegisterRoutes(r, repoRepository)

	srv := startServer(r)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

func initDB(ctx context.Context) (*pgx.Conn, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Connected to database")
	return conn, nil
}

func startServer(r *gin.Engine) *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	return srv
}
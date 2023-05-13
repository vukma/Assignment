package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"

	productapi "product/product-api"
)

var (
	listenAddr  string // Port on which the service will listen on
	postgresURI string // Postgres connection string
)

func init() {

	loc, err := time.LoadLocation("Europe/Belgrade") // Load the correct timezone

	if err != nil {
		log.Fatalf("Error loading timezone: %s\n", err)
	}

	time.Local = loc

	// Load environment file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}

	// Map loaded environment variables to variables declared in var() block
	listenAddr = os.Getenv("PORT")
	postgresURI = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASS"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL"),
	)
}

func main() {
	logger := log.New(os.Stdout, "[product-api] ", log.Ldate|log.Ltime) // Create a custom log for the API

	db, err := sql.Open("postgres", postgresURI) // Connect to Postgres database

	if err != nil {
		logger.Fatalf("Error opening Postgres connection: %s\n", err)
	}

	defer db.Close() // Close the connection before main() func exists

	// Ping the database to check connection
	if err := db.Ping(); err != nil {
		logger.Fatalf("Error contacting Postgres database: %s\n", err)
	}

	logger.Println("Successfully connected to Postgres")

	p := productapi.New(db, logger) // Create a new Product API

	// Run postgres migrations
	if err := p.Migrate(); err != nil {
		logger.Fatalf("Error running migrations: %s\n", err)
	}

	r := httprouter.New() // Create a new router

	// Declare API routes
	r.GET("/healthz", p.Logger(p.Healthz))
	r.GET("/product", p.Logger(p.GetAll))
	r.GET("/product/:id", p.Logger(p.GetByID))
	r.POST("/product", p.Logger(p.PostOne))
	r.PUT("/product/:id", p.Logger(p.UpdateOne))
	r.DELETE("/product/:id", p.Logger(p.DeleteOne))

	// Configure HTTP server
	s := &http.Server{
		Addr:         listenAddr,
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Run the server in a separate go routine
	go func() {
		logger.Printf("Server listening on: %s\n", s.Addr)
		logger.Fatal(s.ListenAndServe()) // Exit the process if server fails
	}()

	c := make(chan os.Signal, 1) // Create channel that will capture OS signals

	// Notify the channel for the following signals
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	sig := <-c // Capure the received signal

	logger.Printf("Signal received: %s; Attempt graceful shutdown\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Create a context with a 30 second timeout

	defer cancel() // Cancel the context before the main() func exits

	// Run graceful shutdown
	if err := s.Shutdown(ctx); err != nil {
		logger.Fatalf("Error shutting down the s: %s\n", err)
	}
}

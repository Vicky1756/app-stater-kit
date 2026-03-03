package routes

import (
	"app-starter-kit/backend/packages/db"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var e *echo.Echo

func StartServer() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	e = echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database Connection
	conn, err := db.ConnectDB()
	if err != nil {
		fmt.Printf("Db connection error occurred: %s\n", err.Error())
		return
	}
	defer conn.Close()

	// Migration Logic
	runMigration := os.Getenv("RUN_MIGRATION")
	dbName := os.Getenv("DATABASE_NAME")
	if runMigration == "true" && conn != nil {
		if err := db.Migrate(conn, dbName); err != nil {
			fmt.Printf("Db migration failed: %s\n", err.Error())
			return
		}
	}

	// CORS Setup
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173", os.Getenv("CLIENT_URL")},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	// Register Handlers
	Handlers(e, conn)

	// Port Formatting (Conscientious fix)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	address := port

	fmt.Printf("Server up at http://localhost%s\n", address)

	// Start using Echo's built-in server management
	e.Logger.Fatal(e.Start(address))
}

func StopServer() {
	if e != nil {
		// Use a timeout context for a graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			fmt.Printf("Shutdown server error: %s\n", err.Error())
		}
	}
}

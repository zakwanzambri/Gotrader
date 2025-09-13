package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zakwanzambri/Gotrader/internal/database"
	"github.com/zakwanzambri/Gotrader/internal/handlers"
	"github.com/zakwanzambri/Gotrader/internal/models"
	"github.com/zakwanzambri/Gotrader/internal/telegram"
	"github.com/zakwanzambri/Gotrader/internal/websocket"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	dbPath := getEnv("DB_PATH", "./signals.db")
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize handlers
	handler := handlers.NewHandler(db)

	// Initialize Telegram bot if token is provided
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken != "" {
		chatIDsStr := os.Getenv("TELEGRAM_CHAT_IDS")
		var chatIDs []int64

		if chatIDsStr != "" {
			for _, idStr := range strings.Split(chatIDsStr, ",") {
				if id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64); err == nil {
					chatIDs = append(chatIDs, id)
				}
			}
		}

		if len(chatIDs) > 0 {
			bot, err := telegram.NewBot(telegramToken, chatIDs, func(signal *models.TradingSignal) {
				// Save signal to database
				if err := db.InsertSignal(signal); err != nil {
					log.Printf("Failed to save signal: %v", err)
					return
				}

				// Broadcast to WebSocket clients
				hub.BroadcastSignal(signal)
				log.Printf("New signal: %s %s @ %.4f", signal.Action, signal.Symbol, signal.Price)
			})

			if err != nil {
				log.Printf("Failed to initialize Telegram bot: %v", err)
			} else {
				go bot.Start()
				log.Println("Telegram bot started")
			}
		} else {
			log.Println("No Telegram chat IDs configured")
		}
	} else {
		log.Println("No Telegram bot token configured")
	}

	// Setup Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	config.AllowOrigins = strings.Split(allowedOrigins, ",")
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// API routes
	api := router.Group("/api/v1")
	{
		api.GET("/health", handler.HealthCheck)
		api.GET("/signals", handler.GetSignals)
		api.PUT("/signals/:id/status", handler.UpdateSignalStatus)
		api.GET("/stats", handler.GetStats)
	}

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		hub.HandleWebSocket(c.Writer, c.Request)
	})

	// Serve static files for frontend
	router.Static("/static", "./frontend/build/static")
	router.StaticFile("/", "./frontend/build/index.html")
	router.StaticFile("/favicon.ico", "./frontend/build/favicon.ico")

	// Catch-all route for React Router
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
		} else {
			c.File("./frontend/build/index.html")
		}
	})

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
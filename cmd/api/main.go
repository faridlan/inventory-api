package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	// Sesuaikan import path ini dengan struktur folder Anda
	"github.com/faridlan/inventory-api/docs"
	_ "github.com/faridlan/inventory-api/docs"
	"github.com/faridlan/inventory-api/internal/config"
	myHttp "github.com/faridlan/inventory-api/internal/delivery/http" // Alias untuk folder http
	"github.com/faridlan/inventory-api/internal/repository/postgres"
	"github.com/faridlan/inventory-api/internal/usecase"
)

// @title Inventory API
// @version 1.0
// @description API untuk manajemen inventori produk pakaian. Dibuat sebagai bahan latihan integrasi Frontend.
// @host localhost:8080
// @BasePath /
func main() {
	// ==========================================
	// 0. LOAD ENV VARS
	// ==========================================
	err := godotenv.Load()
	if err != nil {
		slog.Warn("File .env tidak ditemukan, menggunakan environment variable dari sistem")
	}

	// ==========================================
	// 1. INISIASI DATABASE
	// ==========================================
	// (Pastikan database.ConnectDB() milik Anda sudah membaca dari os.Getenv jika diperlukan)
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	db := config.InitDB(dbUser, dbPassword, dbHost, dbPort, dbName)

	// ==========================================
	// 2. INISIASI REPOSITORY & USECASE
	// ==========================================
	productRepo := postgres.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)

	// ==========================================
	// 3. INISIASI HANDLER & REGISTRY
	// ==========================================
	productHandler := myHttp.NewProductHandler(productUsecase)

	handlers := myHttp.AppHandlers{
		Product: productHandler,
	}

	dbURL := os.Getenv("DB_URL")
	config.RunDBMigration(dbURL)

	swaggerHost := os.Getenv("SWAGGER_HOST")
	if swaggerHost != "" {
		docs.SwaggerInfo.Host = swaggerHost
	}

	// ==========================================
	// 4. SETUP FIBER APP & MIDDLEWARE
	// ==========================================
	app := fiber.New(fiber.Config{
		// Custom global error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":  "Terjadi kesalahan pada sistem internal",
				"detail": err.Error(),
			})
		},
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	// Setup routes API Anda
	myHttp.SetupRoutes(app, handlers)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "*" // Fallback aman untuk keperluan trainee di localhost
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     frontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: false, // Diubah ke true jika nanti menggunakan cookies/session
	}))

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))

	// ==========================================
	// 5. DAFTARKAN SEMUA ROUTE KE FIBER
	// ==========================================
	myHttp.SetupRoutes(app, handlers)

	// ==========================================
	// 6. START SERVER & GRACEFUL SHUTDOWN
	// ==========================================
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		slog.Info("Starting Inventory API Server", slog.String("port", port))
		if err := app.Listen(":" + port); err != nil {
			slog.Error("Server failed to start", slog.String("detail", err.Error()))
		}
	}()

	// Menunggu sinyal interupsi (Ctrl+C atau Docker stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	slog.Info("Menerima sinyal mati, mematikan server dengan sopan...")

	if err := app.Shutdown(); err != nil {
		slog.Error("Server dipaksa mati karena error", slog.String("detail", err.Error()))
	}

	slog.Info("Inventory API Server berhasil dimatikan dengan aman. Sampai jumpa!")
}

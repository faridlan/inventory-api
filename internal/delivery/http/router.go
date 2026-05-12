package http

import (
	"github.com/faridlan/inventory-api/internal/delivery/http/dto"
	"github.com/gofiber/fiber/v2"
)

// AppHandlers menampung semua handler untuk diinjeksi ke router
type AppHandlers struct {
	Product dto.ServerInterface
}

func SetupRoutes(app *fiber.App, h AppHandlers) {

	// ==========================================
	// PUBLIC ROUTES
	// ==========================================
	api := app.Group("/api")

	// Health check sederhana
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Inventory API Server is up and running. 🚀",
			"version": "1.0.0",
		})
	})

	// ==========================================
	// PRODUCT ROUTES
	// ==========================================
	products := api.Group("/products")

	products.Get("/", h.Product.GetProducts)
	products.Post("/", h.Product.PostProducts)
	products.Get("/:id", h.Product.GetProductsId)
	products.Put("/:id", h.Product.PutProductsId)
	products.Patch("/:id", h.Product.PatchProductsId)
	products.Delete("/:id", h.Product.DeleteProductsId)

	// ==========================================
	// ADMIN ROUTES (Endpoint Rahasia)
	// ==========================================
	admin := api.Group("/admin")

	// Hanya Anda yang tahu cara mengakses ini
	admin.Post("/products/reset", h.Product.ResetProducts)
}

package dto

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type ProductResponse struct {
	Id          string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	IsAvailable *bool      `json:"is_available,omitempty" example:"true"`
	Name        string     `json:"name" example:"Kemeja Taktikal Hitam"`
	Price       float64    `json:"price" example:"150000"`
	Sizes       *[]string  `json:"sizes,omitempty" example:"M,L,XL"`
	Stock       *int       `json:"stock,omitempty" example:"50"`
	CreatedAt   *time.Time `json:"created_at,omitempty" example:"2026-05-06T12:00:00Z"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" example:"2026-05-06T12:00:00Z"`
}

type ProductRequest struct {
	Name        string    `json:"name" validate:"required" example:"Kemeja Taktikal Hitam"`
	Price       float64   `json:"price" validate:"required,gte=0" example:"150000"`
	Stock       *int      `json:"stock,omitempty" validate:"omitempty,gte=0" example:"50"`
	Sizes       *[]string `json:"sizes,omitempty" example:"M,L,XL"`
	IsAvailable *bool     `json:"is_available,omitempty" example:"true"`
}

type PatchProductsIdJSONBody struct {
	// Diberi contoh false, karena biasanya endpoint PATCH digunakan
	// oleh FE saat tombol "Setel ke Habis" diklik.
	IsAvailable *bool `json:"is_available" validate:"required" example:"false"`
}

type ServerInterface interface {
	GetProducts(c *fiber.Ctx) error
	PostProducts(c *fiber.Ctx) error
	DeleteProductsId(c *fiber.Ctx) error
	GetProductsId(c *fiber.Ctx) error
	PatchProductsId(c *fiber.Ctx) error
	PutProductsId(c *fiber.Ctx) error
}

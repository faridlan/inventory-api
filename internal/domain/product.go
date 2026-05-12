package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Product merepresentasikan entitas database sekaligus balikan JSON ke Frontend.
type Product struct {
	ID          string
	Name        string
	Price       float64
	Stock       int
	Sizes       []string
	IsAvailable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type ProductInput struct {
	Name        string
	Price       float64
	Stock       int
	Sizes       []string
	IsAvailable *bool
}

// ProductRepository adalah kontrak untuk layer Data Access (GORM/Database).
type ProductRepository interface {
	Fetch() ([]Product, error)
	GetByID(ctx context.Context, id string) (Product, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	PatchAvailability(ctx context.Context, id string, isAvailable bool) error
	Delete(ctx context.Context, id string) error
	ResetDefault(ctx context.Context) error
}

// ProductUsecase adalah kontrak untuk layer Business Logic.
// Inputnya menerima ProductInput dari HTTP handler, memprosesnya, lalu mengirim ke Repository.
type ProductUsecase interface {
	Fetch() ([]Product, error)
	GetByID(ctx context.Context, id string) (Product, error)
	Create(ctx context.Context, input *ProductInput) (Product, error)
	Update(ctx context.Context, id string, input *ProductInput) (Product, error)
	PatchAvailability(ctx context.Context, id string, isAvailable bool) error
	Delete(ctx context.Context, id string) error
	ResetDefault(ctx context.Context) error
}

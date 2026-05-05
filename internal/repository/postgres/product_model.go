package postgres

import (
	"time"

	"github.com/faridlan/inventory-api/internal/domain"
	"gorm.io/gorm"
)

type ProductModel struct {
	// Menggunakan uuid.UUID sebagai tipe data ID
	ID          string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Price       float64        `gorm:"type:numeric(12,2);not null"`
	Stock       int            `gorm:"not null;default:0"`
	Sizes       []string       `gorm:"type:jsonb;serializer:json"`
	IsAvailable bool           `gorm:"not null;default:true"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (ProductModel) TableName() string {
	return "products"
}

func (m *ProductModel) ToDomain() *domain.Product {
	return &domain.Product{
		ID:          m.ID,
		Name:        m.Name,
		Price:       m.Price,
		Stock:       m.Stock,
		Sizes:       m.Sizes,
		IsAvailable: m.IsAvailable,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   gorm.DeletedAt{},
	}
}

func FromDomain(d *domain.Product) *ProductModel {
	return &ProductModel{
		ID:          d.ID,
		Name:        d.Name,
		Price:       d.Price,
		Stock:       d.Stock,
		Sizes:       d.Sizes,
		IsAvailable: d.IsAvailable,
	}
}

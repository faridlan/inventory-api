package postgres

import (
	"context"
	"encoding/json"

	"github.com/faridlan/inventory-api/internal/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Fetch() ([]domain.Product, error) {
	var models []ProductModel

	err := r.db.Find(&models).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	var products []domain.Product
	for _, m := range models {
		products = append(products, *m.ToDomain())
	}

	return products, nil
}

func (r *productRepository) GetByID(ctx context.Context, id string) (domain.Product, error) {
	var model ProductModel

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		return domain.Product{}, TranslateError(err)
	}

	return *model.ToDomain(), nil
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {

	model := FromDomain(product)
	model.ID = ""

	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	product.ID = model.ID
	product.CreatedAt = model.CreatedAt
	product.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	model := FromDomain(product)

	sizesJSON, _ := json.Marshal(model.Sizes)
	result := r.db.WithContext(ctx).Model(&ProductModel{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
		"name":         model.Name,
		"price":        model.Price,
		"stock":        model.Stock,
		"sizes":        gorm.Expr("?::jsonb", sizesJSON),
		"is_available": model.IsAvailable,
	})

	if result.Error != nil {
		return TranslateError(result.Error)
	}

	return nil
}

func (r *productRepository) PatchAvailability(ctx context.Context, id string, isAvailable bool) error {
	result := r.db.WithContext(ctx).Model(&ProductModel{}).Where("id = ?", id).Update("is_available", isAvailable)

	if result.Error != nil {
		return TranslateError(result.Error)
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&ProductModel{})

	if result.Error != nil {
		return TranslateError(result.Error)
	}

	return nil
}

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

func (r *productRepository) ResetDefault(ctx context.Context) error {
	// 1. Kosongkan tabel secara paksa dan bersih
	err := r.db.WithContext(ctx).Exec("TRUNCATE TABLE products").Error
	if err != nil {
		return TranslateError(err)
	}

	// 2. Siapkan 10 data default
	// Karena ID menggunakan UUID dari DB (default:uuid_generate_v4()), kita tidak perlu mengisi ID-nya
	defaultProducts := []ProductModel{
		{Name: "Kemeja Taktikal Hitam", Price: 150000, Stock: 50, Sizes: []string{"M", "L", "XL"}, IsAvailable: true},
		{Name: "Kemeja Taktikal Hijau", Price: 150000, Stock: 30, Sizes: []string{"L", "XL"}, IsAvailable: true},
		{Name: "Kemeja Taktikal Khaki", Price: 155000, Stock: 20, Sizes: []string{"M", "L"}, IsAvailable: true},
		{Name: "Daster Motif Bunga", Price: 85000, Stock: 100, Sizes: []string{"All Size"}, IsAvailable: true},
		{Name: "Daster Arab Renda", Price: 95000, Stock: 75, Sizes: []string{"All Size"}, IsAvailable: true},
		{Name: "Stelan Wanita Kasual", Price: 120000, Stock: 40, Sizes: []string{"M", "L"}, IsAvailable: true},
		{Name: "Stelan Olahraga Muslimah", Price: 180000, Stock: 25, Sizes: []string{"L", "XL", "XXL"}, IsAvailable: true},
		{Name: "Kaos Polos Cotton Combed", Price: 45000, Stock: 200, Sizes: []string{"S", "M", "L", "XL"}, IsAvailable: true},
		{Name: "Celana Sirwal", Price: 110000, Stock: 60, Sizes: []string{"M", "L", "XL"}, IsAvailable: true},
		{Name: "Jaket Windbreaker", Price: 250000, Stock: 15, Sizes: []string{"L", "XL"}, IsAvailable: true},
	}

	// 3. Masukkan 10 data tersebut ke database
	err = r.db.WithContext(ctx).Create(&defaultProducts).Error
	if err != nil {
		return TranslateError(err)
	}

	return nil
}

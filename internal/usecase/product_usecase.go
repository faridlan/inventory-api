package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/faridlan/inventory-api/internal/domain"
)

type productUsecase struct {
	productRepo domain.ProductRepository
}

// NewProductUsecase adalah constructor untuk inisiasi layer usecase
func NewProductUsecase(repo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: repo,
	}
}

func (u *productUsecase) Fetch() ([]domain.Product, error) {
	products, err := u.productRepo.Fetch()
	if err != nil {
		return nil, domain.NewError(domain.ErrInternalServerError, "Gagal mengambil daftar produk")
	}

	return products, nil
}

func (u *productUsecase) GetByID(ctx context.Context, id string) (domain.Product, error) {
	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Product{}, domain.NewError(domain.ErrNotFound, fmt.Sprintf("Produk dengan ID %s tidak ditemukan", id))
		}
		return domain.Product{}, domain.NewError(domain.ErrInternalServerError, "Gagal mengambil detail produk")
	}

	return product, nil
}

func (u *productUsecase) Create(ctx context.Context, input *domain.ProductInput) (domain.Product, error) {
	// Handling pointer boolean untuk IsAvailable
	isAvailable := true // Default sesuai rancangan DB
	if input.IsAvailable != nil {
		isAvailable = *input.IsAvailable
	}

	// Mapping input ke domain product
	product := domain.Product{
		Name:        input.Name,
		Price:       input.Price,
		Stock:       input.Stock,
		Sizes:       input.Sizes,
		IsAvailable: isAvailable,
	}

	err := u.productRepo.Create(ctx, &product)
	if err != nil {
		return domain.Product{}, domain.NewError(domain.ErrInternalServerError, "Gagal menyimpan produk baru")
	}

	return product, nil
}

func (u *productUsecase) Update(ctx context.Context, id string, input *domain.ProductInput) (domain.Product, error) {
	// 1. Cek ketersediaan data terlebih dahulu (Logic ditarik ke Usecase)
	existingProduct, err := u.GetByID(ctx, id)
	if err != nil {
		// Error sudah di-wrap oleh GetByID, langsung kembalikan
		return domain.Product{}, err
	}

	// Handling pointer boolean untuk IsAvailable saat method PUT (Replace utuh)
	isAvailable := false // Default asumsi false jika FE benar-benar tidak mengirim (walau sebaiknya di-validate required)
	if input.IsAvailable != nil {
		isAvailable = *input.IsAvailable
	}

	// 2. Timpa data yang ada dengan data input baru
	existingProduct.Name = input.Name
	existingProduct.Price = input.Price
	existingProduct.Stock = input.Stock
	existingProduct.Sizes = input.Sizes
	existingProduct.IsAvailable = isAvailable

	// 3. Simpan pembaruan ke repository
	err = u.productRepo.Update(ctx, &existingProduct)
	if err != nil {
		return domain.Product{}, err
	}

	return existingProduct, nil
}

func (u *productUsecase) PatchAvailability(ctx context.Context, id string, isAvailable bool) error {
	// 1. Cek ketersediaan data terlebih dahulu
	_, err := u.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Lakukan patch jika data ada
	err = u.productRepo.PatchAvailability(ctx, id, isAvailable)
	if err != nil {
		return domain.NewError(domain.ErrInternalServerError, "Gagal memperbarui status ketersediaan produk")
	}

	return nil
}

func (u *productUsecase) Delete(ctx context.Context, id string) error {
	// 1. Cek ketersediaan data terlebih dahulu
	_, err := u.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Lakukan soft-delete jika data ada
	err = u.productRepo.Delete(ctx, id)
	if err != nil {
		return domain.NewError(domain.ErrInternalServerError, "Gagal menghapus data produk")
	}

	return nil
}

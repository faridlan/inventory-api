package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/faridlan/inventory-api/internal/domain"
	"github.com/faridlan/inventory-api/internal/domain/mocks"
	"github.com/faridlan/inventory-api/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// helper untuk membuat pointer boolean
func boolPtr(b bool) *bool {
	return &b
}

func TestProductUsecase_Fetch(t *testing.T) {
	mockRepo := new(mocks.ProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)

	t.Run("Sukses mengambil data", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil // Reset mock
		expectedProducts := []domain.Product{
			{ID: "1", Name: "Kemeja"},
		}

		mockRepo.On("Fetch").Return(expectedProducts, nil).Once()

		res, err := uc.Fetch()

		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Gagal mengambil data - Error Internal", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("Fetch").Return(nil, errors.New("database down")).Once()

		res, err := uc.Fetch()

		assert.Error(t, err)
		assert.Nil(t, res)
		// Memastikan error yang dikembalikan adalah custom error ErrInternalServerError
		assert.ErrorIs(t, err, domain.ErrInternalServerError)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductUsecase_GetByID(t *testing.T) {
	mockRepo := new(mocks.ProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)
	ctx := context.Background()
	id := "prod-123"

	t.Run("Sukses mengambil detail", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		expectedProduct := domain.Product{ID: id, Name: "Kemeja"}

		mockRepo.On("GetByID", ctx, id).Return(expectedProduct, nil).Once()

		res, err := uc.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, res)
	})

	t.Run("Gagal - Data Tidak Ditemukan", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(domain.Product{}, domain.ErrNotFound).Once()

		res, err := uc.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Empty(t, res)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestProductUsecase_Create(t *testing.T) {
	mockRepo := new(mocks.ProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)
	ctx := context.Background()

	t.Run("Sukses membuat produk", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		input := &domain.ProductInput{
			Name:        "Daster",
			Price:       50000,
			Stock:       10,
			IsAvailable: boolPtr(false), // Test input boolean explicit
		}

		// Mock repository Create, gunakan mock.AnythingOfType karena ID/Timestamp digenerate di Repo
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Product")).Return(nil).Run(func(args mock.Arguments) {
			// Simulasikan kelakuan repository yang mengisi ID
			arg := args.Get(1).(*domain.Product)
			arg.ID = "new-id"
		}).Once()

		res, err := uc.Create(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, "new-id", res.ID)
		assert.Equal(t, input.Name, res.Name)
		assert.False(t, res.IsAvailable) // Memastikan pointer bool di-handle dengan benar
	})
}

func TestProductUsecase_Update(t *testing.T) {
	mockRepo := new(mocks.ProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)
	ctx := context.Background()
	id := "prod-123"

	t.Run("Sukses update produk", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		existingProduct := domain.Product{ID: id, Name: "Lama", IsAvailable: true}
		input := &domain.ProductInput{Name: "Baru", IsAvailable: boolPtr(false)}

		// 1. Usecase akan memanggil GetByID terlebih dahulu
		mockRepo.On("GetByID", ctx, id).Return(existingProduct, nil).Once()

		// 2. Kemudian memanggil Update
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Product")).Return(nil).Once()

		res, err := uc.Update(ctx, id, input)

		assert.NoError(t, err)
		assert.Equal(t, "Baru", res.Name)
		assert.False(t, res.IsAvailable)
	})

	t.Run("Gagal update - Data tidak ditemukan", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		input := &domain.ProductInput{Name: "Baru"}

		// GetByID mengembalikan error NotFound
		mockRepo.On("GetByID", ctx, id).Return(domain.Product{}, domain.ErrNotFound).Once()
		// Update TIDAK BOLEH dipanggil
		mockRepo.AssertNotCalled(t, "Update")

		res, err := uc.Update(ctx, id, input)

		assert.Error(t, err)
		assert.Empty(t, res)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestProductUsecase_PatchAvailability(t *testing.T) {
	mockRepo := new(mocks.ProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)
	ctx := context.Background()
	id := "prod-123"

	t.Run("Sukses patch", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(domain.Product{ID: id}, nil).Once()
		mockRepo.On("PatchAvailability", ctx, id, false).Return(nil).Once()

		err := uc.PatchAvailability(ctx, id, false)

		assert.NoError(t, err)
	})
}

func TestProductUsecase_Delete(t *testing.T) {
	mockRepo := new(mocks.ProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)
	ctx := context.Background()
	id := "prod-123"

	t.Run("Sukses delete", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(domain.Product{ID: id}, nil).Once()
		mockRepo.On("Delete", ctx, id).Return(nil).Once()

		err := uc.Delete(ctx, id)

		assert.NoError(t, err)
	})
}

func TestProductUsecase_ResetDefault(t *testing.T) {
	mockRepo := new(mocks.ProductRepository)
	uc := usecase.NewProductUsecase(mockRepo)
	ctx := context.Background()

	t.Run("Sukses reset data default", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil // Reset status mock

		// Ekspektasi: Repository mengembalikan nil (sukses)
		mockRepo.On("ResetDefault", ctx).Return(nil).Once()

		err := uc.ResetDefault(ctx)

		// Assert: Pastikan tidak ada error yang bocor ke atas
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Gagal reset data - Error dari Repository", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil

		// Ekspektasi: Repository mengalami error (misal koneksi DB terputus)
		mockRepo.On("ResetDefault", ctx).Return(errors.New("database timeout")).Once()

		err := uc.ResetDefault(ctx)

		// Assert: Pastikan mengeluarkan error
		assert.Error(t, err)

		// Assert: Pastikan error aslinya sudah dibungkus dengan baik menjadi ErrInternalServerError
		assert.ErrorIs(t, err, domain.ErrInternalServerError)
		mockRepo.AssertExpectations(t)
	})
}

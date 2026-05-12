package http

import (
	"os"

	"github.com/faridlan/inventory-api/internal/delivery/http/dto"
	"github.com/faridlan/inventory-api/internal/domain"
	"github.com/faridlan/inventory-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type productHandler struct {
	productUsecase domain.ProductUsecase
}

// NewProductHandler menginisiasi handler dan memastikan implementasi ServerInterface
func NewProductHandler(usecase domain.ProductUsecase) dto.ServerInterface {
	return &productHandler{
		productUsecase: usecase,
	}
}

// === HELPER MAPPERS ===

func toProductResponse(product *domain.Product) dto.ProductResponse {
	if product == nil {
		return dto.ProductResponse{}
	}

	// Mapping nilai dari value ke pointer untuk DTO
	isAvailable := product.IsAvailable
	stock := product.Stock
	sizes := product.Sizes

	return dto.ProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Stock:       &stock,
		Sizes:       &sizes,
		IsAvailable: &isAvailable,
		CreatedAt:   &product.CreatedAt,
		UpdatedAt:   &product.UpdatedAt,
	}
}

// === IMPLEMENTASI METHOD INTERFACE ===

// GetProducts mendapatkan semua daftar produk
// @Summary Mendapatkan semua produk
// @Description Mengambil daftar semua produk pakaian yang tersedia di inventori
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} fiber.Map{data=[]dto.ProductResponse,message=string}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/products [get]
func (h *productHandler) GetProducts(c *fiber.Ctx) error {
	products, err := h.productUsecase.Fetch()
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	var res []dto.ProductResponse
	for _, p := range products {
		res = append(res, toProductResponse(&p))
	}

	// Jika data kosong, pastikan mengembalikan array kosong [] bukan null
	if res == nil {
		res = []dto.ProductResponse{}
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil daftar produk", res)
}

// PostProducts menambahkan produk baru
// @Summary Menambahkan produk baru
// @Description Menyimpan data produk pakaian baru ke dalam database
// @Tags Products
// @Accept json
// @Produce json
// @Param request body dto.ProductRequest true "Payload Produk Baru"
// @Success 201 {object} fiber.Map{data=dto.ProductResponse,message=string}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/products [post]
func (h *productHandler) PostProducts(c *fiber.Ctx) error {
	var req dto.ProductRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON tidak valid")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Safely unwrap pointers untuk field opsional
	stock := 0
	if req.Stock != nil {
		stock = *req.Stock
	}

	var sizes []string
	if req.Sizes != nil {
		sizes = *req.Sizes
	}

	reqInput := domain.ProductInput{
		Name:        req.Name,
		Price:       req.Price,
		Stock:       stock,
		Sizes:       sizes,
		IsAvailable: req.IsAvailable,
	}

	product, err := h.productUsecase.Create(c.Context(), &reqInput)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := toProductResponse(&product)
	return utils.SendSuccess(c, fiber.StatusCreated, "Produk berhasil ditambahkan", res)
}

// GetProductsId mendapatkan detail satu produk
// @Summary Mendapatkan detail produk
// @Description Mengambil data spesifik satu produk berdasarkan ID (UUID)
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} fiber.Map{data=dto.ProductResponse,message=string}
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/products/{id} [get]
func (h *productHandler) GetProductsId(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	product, err := h.productUsecase.GetByID(c.Context(), id)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := toProductResponse(&product)
	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil detail produk", res)
}

// PutProductsId mengubah seluruh data produk
// @Summary Mengubah seluruh data produk (Replace)
// @Description Memperbarui seluruh data produk pakaian secara utuh berdasarkan ID (UUID)
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Param request body dto.ProductRequest true "Payload Update Produk"
// @Success 200 {object} fiber.Map{data=dto.ProductResponse,message=string}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/products/{id} [put]
func (h *productHandler) PutProductsId(c *fiber.Ctx) error {

	id := c.Params("id")

	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req dto.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON tidak valid")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	stock := 0
	if req.Stock != nil {
		stock = *req.Stock
	}

	var sizes []string
	if req.Sizes != nil {
		sizes = *req.Sizes
	}

	reqInput := domain.ProductInput{
		Name:        req.Name,
		Price:       req.Price,
		Stock:       stock,
		Sizes:       sizes,
		IsAvailable: req.IsAvailable,
	}

	product, err := h.productUsecase.Update(c.Context(), id, &reqInput)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := toProductResponse(&product)
	return utils.SendSuccess(c, fiber.StatusOK, "Data produk berhasil diperbarui", res)
}

// PatchProductsId mengubah status ketersediaan produk
// @Summary Mengubah status ketersediaan produk (Partial Update)
// @Description Memperbarui hanya status ketersediaan (is_available) dari sebuah produk. Sangat berguna untuk efisiensi payload dari Frontend.
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Param request body dto.PatchProductsIdJSONBody true "Payload Patch Status"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/products/{id} [patch]
func (h *productHandler) PatchProductsId(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req dto.PatchProductsIdJSONBody
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON tidak valid")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	err := h.productUsecase.PatchAvailability(c.Context(), id, *req.IsAvailable)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Status ketersediaan berhasil diubah", nil)
}

// DeleteProductsId menghapus data produk
// @Summary Menghapus produk
// @Description Menghapus data produk berdasarkan ID. Data tidak benar-benar dihapus dari database (Soft Delete).
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/products/{id} [delete]
func (h *productHandler) DeleteProductsId(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	err := h.productUsecase.Delete(c.Context(), id)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Produk berhasil dihapus", nil)
}

// ResetProducts mengembalikan data ke setelan awal (default)
// @Summary Reset Data Default (Secret Endpoint)
// @Description Endpoint rahasia untuk mengosongkan tabel dan mengisi kembali dengan 10 data awal produk.
// @Tags Admin
// @Accept json
// @Produce json
// @Param X-Admin-Token header string true "Token Rahasia Admin"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/admin/products/reset [post]
func (h *productHandler) ResetProducts(c *fiber.Ctx) error {
	// 1. Ambil token dari header request
	adminToken := c.Get("X-Admin-Token")

	// 2. Ambil token asli dari environment variable (VM Azure Anda)
	secretKey := os.Getenv("ADMIN_SECRET_TOKEN")

	// Jika di .env belum diset, kita beri fallback sementara
	if secretKey == "" {
		secretKey = "rahasia-mentor-123"
	}

	// 3. Validasi Token
	if adminToken != secretKey {
		return utils.SendError(c, fiber.StatusUnauthorized, "Akses ditolak: Token admin tidak valid")
	}

	// 4. Eksekusi Reset
	err := h.productUsecase.ResetDefault(c.Context())
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Data produk berhasil direset ke setelan pabrik secara manual", nil)
}

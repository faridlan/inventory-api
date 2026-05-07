package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/faridlan/inventory-api/internal/domain"
	"github.com/robfig/cron/v3"
)

// SetupCronJob akan merakit dan menjalankan jadwal otomatis
func SetupCronJob(productUsecase domain.ProductUsecase) *cron.Cron {
	// Pastikan zona waktu menggunakan WIB (Asia/Jakarta)
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		slog.Error("Gagal memuat zona waktu", "error", err)
		location = time.Local
	}

	// Inisialisasi cron dengan zona waktu spesifik
	c := cron.New(cron.WithLocation(location))

	// Format waktu Cron: "Menit Jam Tanggal Bulan Hari"
	// "0 0 * * *" artinya: Menit ke-0, Jam ke-0 (12 Malam), Setiap Tanggal, Setiap Bulan, Setiap Hari.
	_, err = c.AddFunc("0 0 * * *", func() {
		slog.Info("Mengeksekusi Cron Job: Reset Data Default...")

		ctx := context.Background()
		err := productUsecase.ResetDefault(ctx)
		if err != nil {
			slog.Error("Cron Job Gagal: Gagal mereset data", "error", err)
			return
		}

		slog.Info("Cron Job Sukses: Data berhasil direset ke setelan pabrik (10 produk)!")
	})

	if err != nil {
		slog.Error("Gagal mendaftarkan fungsi cron", "error", err)
	}

	// Mulai jalankan mesin waktu
	c.Start()

	slog.Info("Mesin Cron Job berhasil dinyalakan. Menunggu jadwal eksekusi...")

	return c
}

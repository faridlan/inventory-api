# ==========================================
# STAGE 1: BUILDER (Pabrik Perakitan)
# ==========================================
# Gunakan image Golang versi alpine agar ringan untuk build
FROM golang:1.25-alpine AS builder

# Pasang git jika ada private module/dependency tertentu
RUN apk add --no-cache git

# Set working directory di dalam container
WORKDIR /app

# Copy go.mod dan go.sum lebih dulu
# Trik ini agar Docker me-nyimpan cache download library, 
# sehingga build berikutnya jauh lebih cepat.
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasinya!
# CGO_ENABLED=0 membuat binary yang dihasilkan benar-benar mandiri (statically linked)
# -o inventory-api adalah nama output file binary-nya
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o inventory-api ./cmd/api/main.go

# ==========================================
# STAGE 2: RUNNER (Etalase / Mesin Eksekusi)
# ==========================================
# Gunakan Alpine Linux super kecil (sekitar 5MB)
FROM alpine:3.22

RUN adduser -D appuser
WORKDIR /app

# Install tzdata agar aplikasi Golang membaca zona waktu dengan benar
RUN apk --no-cache add tzdata

# Copy HANYA file binary dari STAGE 1
COPY --from=builder /app/inventory-api .

# Pastikan file binary dimiliki oleh appuser
RUN chown appuser:appuser inventory-api

# Jalankan sebagai user biasa, bukan root
USER appuser

# (Opsional) Copy file .env jika ada konfigurasi default 
# Biasanya saat produksi, variabel di-inject via docker-compose/server
# COPY .env .

# Beritahu Docker port berapa yang digunakan aplikasimu (Misal: 8080)
EXPOSE 8080

# Perintah untuk menjalankan aplikasi saat container menyala
CMD ["./inventory-api"]
# Database Migration

Folder ini menyimpan migration SQL untuk PostgreSQL. Migration dijalankan manual melalui CLI, bukan dari kode aplikasi, agar perubahan skema bersifat versioned, reproducible, mudah direview, dan aman untuk production.

## Struktur Folder

```text
backend/
\-- migrations/
    +-- 000001_create_images_table.up.sql
    +-- 000001_create_images_table.down.sql
    \-- README.md
```

## Konvensi Penamaan

- Gunakan format `NNNNNN_deskripsi.up.sql` untuk apply migration.
- Gunakan format `NNNNNN_deskripsi.down.sql` untuk rollback migration.
- Satu nomor versi harus selalu punya pasangan file `up` dan `down`.

## Tool yang Digunakan

Contoh di bawah menggunakan `golang-migrate`, yang merupakan tool migration umum di ekosistem Go.

Install CLI:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Pastikan binary `migrate` sudah masuk ke `PATH`.

## Format Database URL

Gunakan connection string PostgreSQL yang eksplisit, misalnya:

```text
postgres://postgres:postgres@localhost:5432/image_processing_service?sslmode=disable
```

Sesuaikan user, password, host, port, dan nama database dengan environment Anda.

## Menjalankan Migration

Jalankan dari folder `backend/`:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/image_processing_service?sslmode=disable" up
```

Command di atas akan menjalankan seluruh file `*.up.sql` yang belum pernah diaplikasikan.

## Rollback Migration

Rollback satu versi:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/image_processing_service?sslmode=disable" down 1
```

Rollback seluruh migration sampai bersih:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/image_processing_service?sslmode=disable" down
```

## Melihat Versi Migration

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/image_processing_service?sslmode=disable" version
```

## Force Version

Gunakan hanya jika status migration kotor (`dirty`) dan Anda sudah memastikan kondisi database secara manual:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/image_processing_service?sslmode=disable" force 1
```

Ganti `1` dengan versi yang ingin dijadikan acuan.

## Best Practice

- Jangan gunakan `AutoMigrate()` untuk skema production.
- Review setiap file SQL seperti review kode aplikasi.
- Hindari mengubah file migration yang sudah pernah dijalankan di environment bersama.
- Jika ada perubahan skema baru, tambahkan file migration baru dengan versi berikutnya.
- Simpan logika schema change di SQL murni agar perilaku database jelas dan portable untuk proses review.
- Jalankan migration terlebih dahulu sebelum repository layer mulai menggunakan tabel baru.
- Uji `up` dan `down` di database lokal atau staging sebelum dipakai di environment lain.

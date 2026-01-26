package infrastructure

import (
	"bytes" // <--- WAJIB: Buat convert []byte ke io.Reader
	"fmt"
	"os"

	"github.com/nedpals/supabase-go"
)

type SupabaseStorage struct {
	Client *supabase.Client
	Bucket string
}

func NewSupabaseStorage() *SupabaseStorage {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")
	bucketName := os.Getenv("SUPABASE_BUCKET")

	if url == "" || key == "" {
		panic("Supabase URL/Key belum diisi di .env!")
	}

	client := supabase.CreateClient(url, key)

	return &SupabaseStorage{
		Client: client,
		Bucket: bucketName,
	}
}

// Fitur: Initialize Bucket (Cek Koneksi)
func (s *SupabaseStorage) InitializeBucket() {
	fmt.Println("⚡ Menghubungkan ke Supabase Storage...")

	// PERBAIKAN LIST:
	// List() di library ini mengembalikan []FileObject (1 value), bukan error.
	// Kita tes koneksi dengan mencoba melist 1 file saja.
	results := s.Client.Storage.From(s.Bucket).List("", supabase.FileSearchOptions{
		Limit: 1,
	})

	// Kalau results tidak nil (meskipun kosong), berarti koneksi ke bucket "nyambung".
	// (Library ini gak return error eksplisit di List, jadi kita assume sukses kalo gak panic)
	fmt.Printf("✅ Supabase Storage Connected! Bucket: %s (Found %d items)\n", s.Bucket, len(results))
}

// Fitur: Upload File
func (s *SupabaseStorage) UploadFile(fileBytes []byte, filename string, contentType string) (string, error) {

	var upsert = true

	// 1. UPLOAD KE SUPABASE
	// Perbaikan: Tangkap hasilnya ke variabel 'resp' (bukan err)
	resp := s.Client.Storage.From(s.Bucket).Upload(filename, bytes.NewReader(fileBytes), &supabase.FileUploadOptions{
		ContentType: contentType,
		Upsert:      upsert,
	})

	// Cek manual: Kalau Key kosong, berarti gagal (karena library ini return struct, bukan error)
	if resp.Key == "" {
		return "", fmt.Errorf("gagal upload, tidak ada key yang dikembalikan")
	}

	// 2. GENERATE PUBLIC URL
	// Perbaikan: Hasilnya adalah struct, kita ambil string URL-nya
	urlData := s.Client.Storage.From(s.Bucket).GetPublicUrl(filename)

	// Di library nedpals, hasil URL-nya ada di field .SignedURL
	// (Nama fieldnya memang SignedURL, tapi isinya Public URL kalau bucketnya public)
	publicURL := urlData.SignedUrl

	return publicURL, nil
}

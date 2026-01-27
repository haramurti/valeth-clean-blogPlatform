package infrastructure

import (
	"bytes"
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

// Fitur: Initialize Bucket
func (s *SupabaseStorage) InitializeBucket() {
	// Kita skip cek List biar gak panic
	fmt.Printf("⚡ Supabase Config Loaded. Target Bucket: %s\n", s.Bucket)
}

// Fitur: Upload File
// Fitur: Upload File
func (s *SupabaseStorage) UploadFile(fileBytes []byte, filename string, contentType string) (string, error) {

	// 👇 1. DEBUGGING AWAL (Cek apa yang dikirim)
	fmt.Println("\n--- 🕵️‍♂️ DEBUG UPLOAD START ---")
	fmt.Printf("Target Bucket : '%s'\n", s.Bucket) // Penting! Cek ini di terminal nanti
	fmt.Printf("Nama File     : %s\n", filename)
	fmt.Printf("Ukuran File   : %d bytes\n", len(fileBytes))

	// 2. UPLOAD KE SUPABASE
	resp := s.Client.Storage.From(s.Bucket).Upload(filename, bytes.NewReader(fileBytes), &supabase.FileUploadOptions{
		ContentType: contentType,
		Upsert:      true,
	})

	// 👇 3. DEBUGGING HASIL (Cek apa balasannya)
	fmt.Printf("Status Response: %+v\n", resp) // Print isinya (Key, Id, dll)

	// Cek manual
	if resp.Key == "" {
		fmt.Println("❌ ERROR: Key Kosong! Kemungkinan bucket tidak ketemu atau ditolak Policy.")
		fmt.Println("--------------------------------")
		return "", fmt.Errorf("gagal upload: Key kosong. Pastikan bucket '%s' ada & Policy INSERT aktif", s.Bucket)
	}

	fmt.Println("✅ SUKSES: Key didapat ->", resp.Key)
	fmt.Println("--------------------------------")

	// 4. GENERATE PUBLIC URL
	urlData := s.Client.Storage.From(s.Bucket).GetPublicUrl(filename)
	publicURL := urlData.SignedUrl

	return publicURL, nil
}

package infrastructure

import (
	"bytes"
	"fmt"
	"os"

	"github.com/nedpals/supabase-go"
)

type SupabaseStorage struct {
	Client *supabase.Client
	// Bucket string <-- HAPUS INI, kita gak simpen bucket di struct lagi
}

func NewSupabaseStorage() *SupabaseStorage {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	if url == "" || key == "" {
		panic("Supabase URL/Key belum diisi di .env!")
	}

	client := supabase.CreateClient(url, key)

	return &SupabaseStorage{
		Client: client,
	}
}

func (s *SupabaseStorage) InitializeBucket() {
	fmt.Println("⚡ Supabase Client Ready. Buckets: 'avatars' & 'posts'")
}

// Perubahan: Nambah parameter 'bucketName'
func (s *SupabaseStorage) UploadFile(fileBytes []byte, filename string, contentType string, bucketName string) (string, error) {

	// Upload ke bucket yang diminta
	resp := s.Client.Storage.From(bucketName).Upload(filename, bytes.NewReader(fileBytes), &supabase.FileUploadOptions{
		ContentType: contentType,
		Upsert:      true,
	})

	if resp.Key == "" {
		return "", fmt.Errorf("gagal upload ke bucket '%s'. Cek nama bucket & policy", bucketName)
	}

	// Generate URL
	urlData := s.Client.Storage.From(bucketName).GetPublicUrl(filename)
	return urlData.SignedUrl, nil // Perhatikan huruf besar/kecil SignedURL tergantung versi library
}

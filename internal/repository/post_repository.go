package repository

import (
	"strings"
	"valeth-clean-blogPlatform/internal/domain"

	"gorm.io/gorm"
)

// mysqlPostRepository adalah struct yang megang koneksi DB
type postgresPostRepository struct {
	Conn *gorm.DB
}

// NewPostgresPostRepository buat inisialisasi repo
func NewPostgresPostRepository(Conn *gorm.DB) domain.PostRepository {
	return &postgresPostRepository{Conn}
}

// Fetch -> Ambil semua data (bisa filter search term)
func (m *postgresPostRepository) Fetch(search string) ([]domain.Post, error) {
	var posts []domain.Post

	query := m.Conn

	// Logic filter sesuai requirements: title, content, or category
	if search != "" {
		term := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(category) LIKE ?", term, term, term)
	}

	err := query.Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// GetByID -> Ambil satu
func (m *postgresPostRepository) GetByID(id int) (domain.Post, error) {
	var post domain.Post
	// First() bakal return error RecordNotFound kalo ga ada
	err := m.Conn.First(&post, id).Error
	return post, err
}

// Store -> Simpan baru
func (m *postgresPostRepository) Store(p *domain.Post) error {
	return m.Conn.Create(p).Error
}

// Update -> Update data
func (m *postgresPostRepository) Update(p *domain.Post) error {
	return m.Conn.Save(p).Error
}

// Delete -> Hapus data
func (m *postgresPostRepository) Delete(id int) error {
	return m.Conn.Delete(&domain.Post{}, id).Error
}

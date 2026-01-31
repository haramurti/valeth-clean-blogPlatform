package repository

import (
	"strings"
	"valeth-clean-blogPlatform/internal/domain"

	"gorm.io/gorm"
)

type postgresPostRepository struct {
	Conn *gorm.DB
}

func NewPostgresPostRepository(Conn *gorm.DB) domain.PostRepository {
	return &postgresPostRepository{Conn}
}

func (m *postgresPostRepository) Fetch(search string) ([]domain.Post, error) {
	var posts []domain.Post

	// Query dasar: Gabungin User
	query := m.Conn.Debug().Preload("User").Joins("User")

	if search != "" {
		// Ubah search term jadi lowercase biar gak case-sensitive
		term := "%" + strings.ToLower(search) + "%"

		// --- BAGIAN INI YANG KITA UPDATE ---
		// Kita tambahkan logika OR array_to_string(...)
		query = query.Where(
			`LOWER(posts.title) LIKE ? 
			OR LOWER(posts.content) LIKE ? 
			OR LOWER("User".name) LIKE ? 
			OR array_to_string(posts.tags, ',') ILIKE ?`, // 👈 INI BARU
			term, term, term, term, // 👈 Jangan lupa tambah variabel term ke-4
		)
	}

	// Urutkan dari terbaru
	err := query.Order("posts.created_at desc").Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *postgresPostRepository) GetByID(id int) (domain.Post, error) {
	var post domain.Post

	err := m.Conn.First(&post, id).Error
	return post, err
}

func (m *postgresPostRepository) Store(p *domain.Post) error {
	return m.Conn.Create(p).Error
}

func (m *postgresPostRepository) Update(p *domain.Post) error {
	return m.Conn.Save(p).Error
}

func (m *postgresPostRepository) Delete(id int) error {
	return m.Conn.Delete(&domain.Post{}, id).Error
}

func (m *postgresPostRepository) FetchByUserID(userID int) ([]domain.Post, error) {
	var posts []domain.Post
	// Ambil post where user_id = X, urutkan dari terbaru
	err := m.Conn.Preload("User").Where("user_id = ?", userID).Order("created_at desc").Find(&posts).Error
	return posts, err
}

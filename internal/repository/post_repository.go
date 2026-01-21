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

	query := m.Conn.Preload("User")

	if search != "" {
		term := "%" + strings.ToLower(search) + "%"
		// Logic search diperluas: Cari Judul ATAU Konten ATAU Nama Penulis
		query = query.Joins("User").Where(
			"LOWER(posts.title) LIKE ? OR LOWER(posts.content) LIKE ? OR LOWER(User.name) LIKE ?",
			term, term, term,
		)
	}

	err := query.Order("created_at desc").Find(&posts).Error
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

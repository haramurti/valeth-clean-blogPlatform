package repository

import (
	"valeth-clean-blogPlatform/internal/domain"

	"gorm.io/gorm"
)

type postgresCommentRepository struct {
	Conn *gorm.DB
}

func NewPostgresCommentRepository(Conn *gorm.DB) domain.CommentRepository {
	return &postgresCommentRepository{Conn}
}

func (r *postgresCommentRepository) Create(c *domain.Comment) error {
	return r.Conn.Create(c).Error
}

func (r *postgresCommentRepository) GetByPostID(postID int) ([]domain.Comment, error) {
	var comments []domain.Comment
	// Preload User biar foto & namanya kebawa
	err := r.Conn.Preload("User").Where("post_id = ?", postID).Order("created_at desc").Find(&comments).Error
	return comments, err
}

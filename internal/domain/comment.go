package domain

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"` // Biar tau siapa yg komen
	CreatedAt time.Time `json:"created_at"`
}

// Kontrak buat Gudang (Repository)
type CommentRepository interface {
	Create(comment *Comment) error
	GetByPostID(postID int) ([]Comment, error)
}

// Kontrak buat Manajer (UseCase)
type CommentUseCase interface {
	Create(comment *Comment) error
	GetByPostID(postID int) ([]Comment, error)
}

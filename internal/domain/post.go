package domain

import (
	"time"
)

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	Category  string    `json:"category" gorm:"not null"`
	Tags      []string  `json:"tags" gorm:"type:text[]"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id"` // Foreign Key
	User      User      `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type PostRepository interface {
	Fetch(search string) ([]Post, error)
	GetByID(id int) (Post, error)
	Store(p *Post) error
	Update(p *Post) error
	Delete(id int) error
	FetchByUserID(userID int) ([]Post, error)
}

type PostUseCase interface {
	Fetch(search string) ([]Post, error)
	GetByID(id int) (Post, error)
	Store(p *Post) error
	Update(id int, p *Post) error
	Delete(id int) error
	FetchByUserID(userID int) ([]Post, error)
}

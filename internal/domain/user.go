package domain

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	GoogleID  string    `json:"google_id" gorm:"uniqueIndex"` // Kunci utama login Google
	Email     string    `json:"email" gorm:"uniqueIndex"`     // Email harus unik
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"` // URL foto profil Google
	CreatedAt time.Time `json:"created_at"`

	// Relasi: Satu User bisa punya banyak Post
	// 'omitempty' biar kalau di-json-kan user gak bawa gerbong postingan kalau gak diminta
	Posts []Post `json:"posts,omitempty" gorm:"foreignKey:UserID"`
}

// Kontrak kerja buat ngurus User di Database nanti
type UserRepository interface {
	FindByGoogleID(googleID string) (User, error)
	Create(user *User) error
	GetByID(id int) (User, error)
	Update(user *User) error
}

// Kontrak kerja buat Logic User (Business Logic)
type UserUseCase interface {
	CheckGoogleLogin(string) (User, error)
	RegisterUser(user *User) error
	GetProfile(id int) (User, error)
	UpdateUser(user *User) error
}

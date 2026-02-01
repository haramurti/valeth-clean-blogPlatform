package domain

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"
)

// 1. Bikin Tipe Custom buat handle Array Postgres
type StringArray []string

// Scan: Mengubah data dari Database (string/bytes) menjadi []string di Go
func (a *StringArray) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New("incompatible type for StringArray")
	}

	// Format Postgres Array biasanya: "{tag1,tag2}"
	// Kita buang kurung kurawal {}
	source = strings.Trim(source, "{}")

	if source == "" {
		*a = []string{}
		return nil
	}

	// Pecah berdasarkan koma
	*a = strings.Split(source, ",")
	return nil
}

// Value: Mengubah []string di Go menjadi format Array Postgres "{tag1,tag2}"
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}

	// Gabungkan array jadi satu string dengan pemisah koma, lalu bungkus {}
	// Contoh hasil: "{nature,hiking}"
	return "{" + strings.Join(a, ",") + "}", nil
}

type Post struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Title    string `json:"title" gorm:"not null"`
	Content  string `json:"content" gorm:"not null"`
	Category string `json:"category" gorm:"not null"`

	// 2. Ganti tipe datanya dari []string jadi StringArray
	Tags StringArray `json:"tags" gorm:"type:text[]"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Image     string    `json:"image" form:"image"`

	IsBookmarked bool `json:"is_bookmarked" gorm:"-"`
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

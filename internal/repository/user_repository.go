package repository

import (
	"valeth-clean-blogPlatform/internal/domain"

	"gorm.io/gorm"
)

type postgresUserRepository struct {
	Conn *gorm.DB
}

func NewPostgresUserRepository(Conn *gorm.DB) domain.UserRepository {
	return &postgresUserRepository{Conn}
}

// Cari user berdasarkan GoogleID (Dipake pas Login)
func (r *postgresUserRepository) FindByGoogleID(googleID string) (domain.User, error) {
	var user domain.User
	// "Cari di tabel users dimana google_id = ?"
	err := r.Conn.Where("google_id = ?", googleID).First(&user).Error
	return user, err
}

// Simpan user baru (Dipake pas Register otomatis)
func (r *postgresUserRepository) Create(user *domain.User) error {
	return r.Conn.Create(user).Error
}

// Cari user berdasarkan ID biasa (Dipake buat Profile)
func (r *postgresUserRepository) GetByID(id int) (domain.User, error) {
	var user domain.User
	err := r.Conn.First(&user, id).Error
	return user, err
}

// Tambahin fungsi ini buat update data user
func (r *postgresUserRepository) Update(user *domain.User) error {
	return r.Conn.Save(user).Error
}

// baru kamis 22
// 3. [BARU] Cari by Email (Penting buat Login!)
func (r *postgresUserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	// SELECT * FROM users WHERE email = '...' LIMIT 1
	err := r.Conn.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 4. [BARU] Simpan User Baru (Register)
func (r *postgresUserRepository) Store(user *domain.User) error {
	return r.Conn.Create(user).Error
}

// 5. Update Data User (Sync Avatar)

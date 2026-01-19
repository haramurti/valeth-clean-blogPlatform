package usecase

import (
	"errors"
	"valeth-clean-blogPlatform/internal/domain"

	"gorm.io/gorm"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

// Ganti nama fungsinya biar lebih jelas, bukan LoginOrRegister lagi
func (u *userUseCase) CheckGoogleLogin(googleID string) (domain.User, error) {
	// 1. Cek di database
	user, err := u.userRepo.FindByGoogleID(googleID)

	// 2. Kalau ada error DB (selain Not Found), lapor error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, err
	}

	// 3. Kalau user ketemu, balikin datanya (Login Sukses)
	if user.ID != 0 {
		return user, nil
	}

	// 4. Kalau user GAK KETEMU (ID-nya 0 atau RecordNotFound)
	// Jangan Create user baru! Balikin error khusus.
	// Nanti Handler bakal nangkep error ini buat nge-redirect.
	return domain.User{}, errors.New("USER_NOT_FOUND")
}

func (u *userUseCase) GetProfile(id int) (domain.User, error) {
	return u.userRepo.GetByID(id)
}

// ... fungsi CheckGoogleLogin yang lama ...

// --- TAMBAHAN BARU ---
// Fungsi buat Simpan User Baru (Finalisasi Register)
func (u *userUseCase) RegisterUser(newUser *domain.User) error {
	// Validasi basic aja
	if newUser.Email == "" || newUser.Name == "" {
		return errors.New("nama dan email gaboleh kosong bro")
	}

	// Suruh Repo simpan
	return u.userRepo.Create(newUser)
}

package usecase

import (
	"errors"
	"valeth-clean-blogPlatform/internal/domain"

	"golang.org/x/crypto/bcrypt"

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

// kamis 22 baru
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

// Tambahin ini buat nyambungin Handler ke Repo
func (u *userUseCase) UpdateUser(user *domain.User) error {
	return u.userRepo.Update(user)
}

// ... kode yang lama ...

// 1. REGISTER MANUAL (HASHING)
func (u *userUseCase) RegisterManual(user *domain.User) error {
	// Cek dulu email udah ada belum?
	existingUser, _ := u.userRepo.GetByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email sudah terdaftar bro")
	}

	// Acak Password (Hashing)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Simpan versi hash-nya, BUKAN plain text
	user.Password = string(hashedPassword)

	// Set Avatar default kalau kosong

	return u.userRepo.Store(user)
}

// 2. LOGIN MANUAL (COMPARING)
func (u *userUseCase) LoginManual(email, password string) (*domain.User, error) {
	// Cari User by Email
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("email tidak ditemukan")
	}

	// Bandingkan Password Input vs Password Hash di DB
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("password salah bro")
	}

	return user, nil
}

// Implementasi
func (u *userUseCase) ToggleBookmark(userID int, postID int) error {
	return u.userRepo.ToggleBookmark(userID, postID)
}

func (u *userUseCase) GetBookmarks(userID int) ([]domain.Post, error) {
	return u.userRepo.GetBookmarks(userID)
}

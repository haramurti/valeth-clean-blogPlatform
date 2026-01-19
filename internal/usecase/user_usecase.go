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

// Logic Login Paling Sakti
func (u *userUseCase) LoginOrRegister(googleID, email, name, avatar string) (domain.User, error) {
	// 1. Cek dulu, user ini udah pernah login belum?
	user, err := u.userRepo.FindByGoogleID(googleID)

	// 2. Kalau errornya bukan "Record Not Found", berarti ada error DB beneran
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, err
	}

	// 3. Kalau user DITEMUKAN (err == nil), langsung balikin usernya
	if user.ID != 0 {
		return user, nil
	}

	// 4. Kalau user TIDAK DITEMUKAN, berarti dia User Baru -> REGISTER
	newUser := domain.User{
		GoogleID: googleID,
		Email:    email,
		Name:     name,
		Avatar:   avatar,
	}

	err = u.userRepo.Create(&newUser)
	if err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}

func (u *userUseCase) GetProfile(id int) (domain.User, error) {
	return u.userRepo.GetByID(id)
}

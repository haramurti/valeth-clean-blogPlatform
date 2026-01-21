package usecase

import (
	"errors"
	"valeth-clean-blogPlatform/internal/domain"
)

// 1. STRUCT (SI BOS)
// Dia nyimpen "Nomor HP" anak buahnya (Repository)
type postUseCase struct {
	postRepo domain.PostRepository
}

// 2. CONSTRUCTOR (REKRUTMEN)
// "Lu mau bikin Bos baru? Kasih tau dulu siapa anak buahnya (p)"
func NewPostUseCase(p domain.PostRepository) domain.PostUseCase {
	return &postUseCase{
		postRepo: p,
	}
}

// 1. Fetch (Udah dibahas)
func (u *postUseCase) Fetch(search string) ([]domain.Post, error) {
	return u.postRepo.Fetch(search)
}

// 2. GetByID (Udah dibahas)
func (u *postUseCase) GetByID(id int) (domain.Post, error) {
	return u.postRepo.GetByID(id)
}

// 3. Store (INI HARUS ADA)
func (u *postUseCase) Store(p *domain.Post) error {
	// Validasi bisnis, misal: Judul gak boleh kosong
	if p.Title == "" {
		return errors.New("judul tidak boleh kosong bro")
	}

	// Kalau lolos validasi, baru suruh repo simpan
	return u.postRepo.Store(p)
}

// 4. Update (INI HARUS ADA)
func (u *postUseCase) Update(id int, p *domain.Post) error {
	// Cek dulu datanya ada gak?
	_, err := u.postRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Kalau ada, baru update
	return u.postRepo.Update(p)
}

// 5. Delete (INI HARUS ADA)
func (u *postUseCase) Delete(id int) error {
	// Sama, cek dulu ada gak barangnya
	_, err := u.postRepo.GetByID(id)
	if err != nil {
		return err
	}

	return u.postRepo.Delete(id)
}

func (u *postUseCase) FetchByUserID(userID int) ([]domain.Post, error) {
	return u.postRepo.FetchByUserID(userID)
}

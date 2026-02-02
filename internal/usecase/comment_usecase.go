package usecase

import "valeth-clean-blogPlatform/internal/domain"

type commentUseCase struct {
	commentRepo domain.CommentRepository
}

func NewCommentUseCase(r domain.CommentRepository) domain.CommentUseCase {
	return &commentUseCase{commentRepo: r}
}

func (u *commentUseCase) Create(c *domain.Comment) error {
	return u.commentRepo.Create(c)
}

func (u *commentUseCase) GetByPostID(postID int) ([]domain.Comment, error) {
	return u.commentRepo.GetByPostID(postID)
}

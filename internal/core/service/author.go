package service

import (
	"context"
	"fmt"
	"github.com/olund/cool/internal/core/domain"
	"github.com/olund/cool/internal/core/ports"
)

var _ ports.Authors = &AuthorService{}

type AuthorService struct {
	authorStore ports.AuthorStore
}

func NewAuthorService(authorStore ports.AuthorStore) *AuthorService {
	return &AuthorService{authorStore: authorStore}
}

func (s *AuthorService) Create(ctx context.Context, req domain.CreateAuthorRequest) (domain.Author, error) {
	// todo validation.
	return s.authorStore.Insert(ctx, req)
}

func (s *AuthorService) GetById(ctx context.Context, id int64) (domain.Author, error) {
	return s.authorStore.GetById(ctx, id)
}

func (s *AuthorService) ListAll(ctx context.Context) ([]domain.Author, error) {
	authors, err := s.authorStore.ListAuthors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list authors: %w", err)
	}

	return authors, nil
}

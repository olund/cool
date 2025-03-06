package author

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/olund/cool/internal/core/domain"
	"github.com/olund/cool/internal/core/ports"
)

var _ ports.AuthorStore = &AuthorStore{}

type AuthorStore struct {
	queries *Queries
}

func NewAuthorStore(queries *Queries) *AuthorStore {
	return &AuthorStore{
		queries: queries,
	}
}

func (s *AuthorStore) Insert(ctx context.Context, author domain.CreateAuthorRequest) (domain.Author, error) {
	created, err := s.queries.CreateAuthor(ctx, CreateAuthorParams{
		Name: author.Name,
		Bio: pgtype.Text{
			String: author.Bio,
		},
	})
	if err != nil {
		return domain.Author{}, err
	}
	return ToAuthor(created), err
}

func (s *AuthorStore) GetById(ctx context.Context, id int64) (domain.Author, error) {
	author, err := s.queries.GetAuthor(ctx, id)
	if err != nil {
		return domain.Author{}, err
	}
	return ToAuthor(author), err
}

func (s *AuthorStore) ListAuthors(ctx context.Context) ([]domain.Author, error) {
	authors, err := s.queries.ListAuthors(ctx)
	if err != nil {
		return nil, err
	}

	var ret []domain.Author
	for _, author := range authors {
		ret = append(ret, ToAuthor(author))
	}
	return ret, nil
}

func ToAuthor(author Author) domain.Author {
	return domain.Author{
		Id:   author.ID,
		Name: author.Name,
		Bio:  author.Bio.String,
	}
}

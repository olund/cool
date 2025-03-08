package domain

import "context"

type Author struct {
	Id   int64
	Name string
	Bio  string
}

type CreateAuthorRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

func (c CreateAuthorRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if len(c.Name) == 0 {
		problems["Name"] = "Name cannot be empty"
	}
	if len(c.Bio) == 0 {
		problems["Bio"] = "Bio cannot be empty"
	}
	return problems
}

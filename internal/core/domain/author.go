package domain

type Author struct {
	Id   int64
	Name string
	Bio  string
}

type CreateAuthorRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

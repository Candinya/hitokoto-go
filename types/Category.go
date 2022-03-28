package types

type Category struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	Key         string `json:"key"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Path        string `json:"path"`
}

type Categories []Category

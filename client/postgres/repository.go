package postgres

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

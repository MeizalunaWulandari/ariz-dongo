package repository

type Repository struct {
	// Database akan ditambahkan di sini nanti.
	//
	// Contoh future:
	// DB *sql.DB
}

func New() *Repository {
	return &Repository{}
}

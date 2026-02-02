package postgres

import (
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"
	postgresInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

type AuthorRepository interface {
	WithinTX(tx *gorm.DB) AuthorRepository
	CreateAuthor(author *model.Author) error
	FindByID(ID uint) (*model.Author, error)
	DeleteAuthor(ID uint) error
}

type authorRepository struct {
	name string
	db   *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &authorRepository{name: "Author", db: db}
}

func (r *authorRepository) WithinTX(tx *gorm.DB) AuthorRepository {
	return &authorRepository{db: tx}
}

func (r *authorRepository) CreateAuthor(author *model.Author) error {
	if err := r.db.Create(author).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}

func (r *authorRepository) FindByID(ID uint) (*model.Author, error) {
	var author model.Author
	if err := r.db.Where("id = ?", ID).First(&author).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return &author, nil
}

func (r *authorRepository) DeleteAuthor(ID uint) error {
	if err := r.db.Delete(&model.Author{}, ID).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}

package repository

import (
	"errors"

	"github.com/Yarik7610/library-backend/catalog-service/internal/dto"
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"gorm.io/gorm"
)

type BookRepository interface {
	GetCategories() ([]string, error)
	CreateBook(book *model.Book) error
	CountBooks() (int64, error)
	FindByID(ID uint) (*model.Book, error)
	GetAuthorsBooks(authorName string) ([]dto.AuthorBooksRaw, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) GetCategories() ([]string, error) {
	categories := make([]string, 0)
	if err := r.db.Model(&model.Book{}).Distinct().Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *bookRepository) CreateBook(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) CountBooks() (int64, error) {
	var bookCount int64
	err := r.db.Model(&model.Book{}).Count(&bookCount).Error
	return bookCount, err
}

func (r *bookRepository) FindByID(ID uint) (*model.Book, error) {
	var book model.Book
	if err := r.db.Where("id = ?", ID).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) GetAuthorsBooks(authorName string) ([]dto.AuthorBooksRaw, error) {
	var rawAuthorsBooks []dto.AuthorBooksRaw

	const query = `
		SELECT 
			b.author_id, 
			a.fullname,
			json_agg(
				json_build_object(
					'id', b.id,
					'created_at', b.created_at,
					'title', b.title,
					'year', b.year,
					'category', b.category
				)
			) books
		FROM books b
		INNER JOIN authors a
		ON b.author_id = a.id
		WHERE a.fullname ILIKE ?
		GROUP BY b.author_id, a.fullname
		ORDER BY b.author_id
	`

	if err := r.db.Raw(query, "%"+authorName+"%").Scan(&rawAuthorsBooks).Error; err != nil {
		return nil, err
	}

	return rawAuthorsBooks, nil
}

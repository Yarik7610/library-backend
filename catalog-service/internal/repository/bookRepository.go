package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Yarik7610/library-backend/catalog-service/internal/dto"
	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"gorm.io/gorm"
)

type BookRepository interface {
	WithinTX(tx *gorm.DB) BookRepository
	GetCategories() ([]string, error)
	CountBooks() (int64, error)
	FindByID(ID uint) (*model.Book, error)
	FindByTitleAndAuthorID(title string, authorID uint) (*model.Book, error)
	GetBooksByAuthorID(authorID uint) ([]model.Book, error)
	ListBooksByAuthorName(authorName string, page, count uint, sort, order string) ([]dto.ListedBooksRaw, error)
	ListBooksByTitle(title string, page, count uint, sort, order string) ([]dto.ListedBooksRaw, error)
	ListBooksByAuthorNameAndTitle(authorName, title string, page, count uint, sort, order string) ([]dto.ListedBooksRaw, error)
	ListBooksByCategory(categoryName string, page, count uint, sort, order string) ([]dto.ListedBooksRaw, error)
	CreateBook(book *model.Book) error
	DeleteBook(ID uint) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) WithinTX(tx *gorm.DB) BookRepository {
	return &bookRepository{db: tx}
}

func (r *bookRepository) GetCategories() ([]string, error) {
	time.Sleep(200 * time.Millisecond)
	categories := make([]string, 0)
	if err := r.db.Model(&model.Book{}).Distinct().Order("category").Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *bookRepository) CreateBook(book *model.Book) error {
	book.Category = strings.ToLower(book.Category)
	return r.db.Create(book).Error
}

func (r *bookRepository) DeleteBook(ID uint) error {
	return r.db.Delete(&model.Book{}, ID).Error
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

func (r *bookRepository) FindByTitleAndAuthorID(title string, authorID uint) (*model.Book, error) {
	var book model.Book
	if err := r.db.Where("title = ?", title).Where("author_id = ?", authorID).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) GetBooksByAuthorID(authorID uint) ([]model.Book, error) {
	var books []model.Book
	if err := r.db.Where("author_id = ?", authorID).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) ListBooksByAuthorName(authorName string, page, count uint, sort, order string) ([]dto.ListedBooksRaw, error) {
	return r.listBooksBy(map[string]string{"author": authorName}, page, count, sort, order)
}

func (r *bookRepository) ListBooksByTitle(title string, page, count uint, sort, order string) ([]dto.ListedBooksRaw, error) {
	return r.listBooksBy(map[string]string{"title": title}, page, count, sort, order)
}

func (r *bookRepository) ListBooksByAuthorNameAndTitle(authorName, title string, page, count uint, sort, order string) ([]dto.ListedBooksRaw, error) {
	return r.listBooksBy(map[string]string{"author": authorName, "title": title}, page, count, sort, order)
}

func (r *bookRepository) ListBooksByCategory(category string, page, count uint, sort, order string) ([]dto.ListedBooksRaw, error) {
	return r.listBooksBy(map[string]string{"category": category}, page, count, sort, order)
}

func (r *bookRepository) listBooksBy(filters map[string]string, page, count uint, sort, order string) ([]dto.ListedBooksRaw, error) {
	var rawBooks []dto.ListedBooksRaw
	offset := (page - 1) * count

	sort, order = validateOrderParams(sort, order)

	whereClauses := []string{}
	args := []any{}

	if v, ok := filters["author"]; ok && v != "" {
		whereClauses = append(whereClauses, "a.fullname ILIKE ?")
		args = append(args, "%"+v+"%")
	}
	if v, ok := filters["title"]; ok && v != "" {
		whereClauses = append(whereClauses, "b.title ILIKE ?")
		args = append(args, "%"+v+"%")
	}
	if v, ok := filters["category"]; ok && v != "" {
		whereClauses = append(whereClauses, "b.category ILIKE ?")
		args = append(args, "%"+v+"%")
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT 
			b.author_id, 
			a.fullname,
			json_agg(
				json_build_object(
					'id', b.id,
					'title', b.title,
					'year', b.year,
					'category', b.category,
					'created_at', b.created_at
				) ORDER BY %s %s
			) AS books
		FROM books b
		INNER JOIN authors a ON b.author_id = a.id
		%s
		GROUP BY b.author_id, a.fullname
		ORDER BY b.author_id
		LIMIT ? OFFSET ?
	`, sort, order, whereSQL)

	args = append(args, count, offset)

	if err := r.db.Raw(query, args...).Scan(&rawBooks).Error; err != nil {
		return nil, err
	}

	return rawBooks, nil
}

func (r *bookRepository) AddBook(bookDTO *dto.AddBook) (*model.Book, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var author model.Author
		if err := tx.Where("author_id = ?", bookDTO.AuthorID).First(&author).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
			}
		}

		return nil
	})
	return nil, err
}

func validateOrderParams(sort, order string) (string, string) {
	sort = strings.ToLower(sort)
	if sort != "title" && sort != "year" && sort != "category" {
		sort = "title"
	}
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}
	return sort, order
}

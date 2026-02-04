package postgres

import (
	"fmt"
	"strings"

	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"

	postgresInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

type BookRepository interface {
	WithinTX(tx *gorm.DB) BookRepository
	GetCategories() ([]string, error)
	GetNewBooks() ([]model.Book, error)
	GetBooksByIDs(bookIDs []string) ([]model.Book, error)
	GetBooksByAuthorID(authorID uint) ([]model.Book, error)
	FindByID(ID uint) (*model.Book, error)
	CountBooks() (int64, error)
	CreateBook(book *model.Book) error
	DeleteBook(ID uint) error
	ListBooksByAuthorName(authorName string, page, count uint, sort, order string) ([]model.BookWithAuthor, error)
	ListBooksByTitle(title string, page, count uint, sort, order string) ([]model.BookWithAuthor, error)
	ListBooksByAuthorNameAndTitle(authorName, title string, page, count uint, sort, order string) ([]model.BookWithAuthor, error)
	ListBooksByCategory(categoryName string, page, count uint, sort, order string) ([]model.BookWithAuthor, error)
}

type bookRepository struct {
	name string
	db   *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{name: "Book(s)", db: db}
}

func (r *bookRepository) WithinTX(tx *gorm.DB) BookRepository {
	return &bookRepository{name: "Books(s)", db: tx}
}

func (r *bookRepository) GetCategories() ([]string, error) {
	var categories []string
	if err := r.db.Model(&model.Book{}).Distinct().Order("category").Pluck("category", &categories).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return categories, nil
}

func (r *bookRepository) GetNewBooks() ([]model.Book, error) {
	const NEW_BOOKS_COUNT = 10

	var newBooks []model.Book
	if err := r.db.Order("created_at DESC").Limit(NEW_BOOKS_COUNT).Find(&newBooks).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return newBooks, nil
}

func (r *bookRepository) GetBooksByIDs(bookIDs []string) ([]model.Book, error) {
	var books []model.Book
	if err := r.db.Where("id IN ?", bookIDs).Find(&books).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return books, nil
}

func (r *bookRepository) GetBooksByAuthorID(authorID uint) ([]model.Book, error) {
	var books []model.Book
	if err := r.db.Where("author_id = ?", authorID).Find(&books).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return books, nil
}

func (r *bookRepository) FindByID(ID uint) (*model.Book, error) {
	var book model.Book
	if err := r.db.Where("id = ?", ID).First(&book).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return &book, nil
}

func (r *bookRepository) CountBooks() (int64, error) {
	var bookCount int64
	if err := r.db.Model(&model.Book{}).Count(&bookCount).Error; err != nil {
		return 0, postgresInfrastructure.NewError(err, r.name)
	}
	return bookCount, nil
}

func (r *bookRepository) CreateBook(book *model.Book) error {
	book.Category = strings.ToLower(book.Category)
	if err := r.db.Create(book).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}

func (r *bookRepository) DeleteBook(ID uint) error {
	if err := r.db.Delete(&model.Book{}, ID).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}

func (r *bookRepository) ListBooksByAuthorName(authorName string, page, count uint, sort, order string) ([]model.BookWithAuthor, error) {
	return r.listBooksBy(map[string]string{"author": authorName}, page, count, sort, order)
}

func (r *bookRepository) ListBooksByTitle(title string, page, count uint, sort, order string) ([]model.BookWithAuthor, error) {
	return r.listBooksBy(map[string]string{"title": title}, page, count, sort, order)
}

func (r *bookRepository) ListBooksByAuthorNameAndTitle(authorName, title string, page, count uint, sort, order string) ([]model.BookWithAuthor, error) {
	return r.listBooksBy(map[string]string{"author": authorName, "title": title}, page, count, sort, order)
}

func (r *bookRepository) ListBooksByCategory(category string, page, count uint, sort, order string) ([]model.BookWithAuthor, error) {
	return r.listBooksBy(map[string]string{"category": category}, page, count, sort, order)
}

func (r *bookRepository) listBooksBy(filters map[string]string, page, count uint, sort, order string) ([]model.BookWithAuthor, error) {
	var bookModels []model.BookWithAuthor
	offset := (page - 1) * count

	sort, order = sanitizeListBooksParams(sort, order)

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
			b.id,
			b.author_id,
			a.fullname author_fullname,
			b.title, 
			b.year,
			b.category
		FROM books b
		INNER JOIN authors a 
		ON b.author_id = a.id
		%s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereSQL, sort, order)

	args = append(args, count, offset)
	if err := r.db.Raw(query, args...).Scan(&bookModels).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return bookModels, nil
}

func sanitizeListBooksParams(sort, order string) (string, string) {
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

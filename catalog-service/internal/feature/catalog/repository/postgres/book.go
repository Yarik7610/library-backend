package postgres

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"

	postgresInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

type BookRepository interface {
	WithinTX(tx *gorm.DB) BookRepository
	GetCategories() ([]string, error)
	GetNewBooks() ([]model.BookWithAuthor, error)
	GetBooksByIDs(bookIDs []string) ([]model.BookWithAuthor, error)
	GetBooksByAuthorID(authorID uint) ([]model.BookWithAuthor, error)
	FindByID(ID uint) (*model.BookWithAuthor, error)
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

func (r *bookRepository) GetNewBooks() ([]model.BookWithAuthor, error) {
	const NEW_BOOKS_COUNT = 10

	var booksWithAuthor []model.BookWithAuthor

	err := r.buildBaseBookWithAuthorQuery().
		Order("books.created_at DESC").
		Limit(NEW_BOOKS_COUNT).
		Find(&booksWithAuthor).Error
	if err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}

	return booksWithAuthor, nil
}

func (r *bookRepository) GetBooksByIDs(bookIDs []string) ([]model.BookWithAuthor, error) {
	var booksWithAuthor []model.BookWithAuthor

	err := r.buildBaseBookWithAuthorQuery().
		Where("books.id IN ?", bookIDs).
		Find(&booksWithAuthor).Error
	if err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}

	return booksWithAuthor, nil
}

func (r *bookRepository) GetBooksByAuthorID(authorID uint) ([]model.BookWithAuthor, error) {
	var booksWithAuthor []model.BookWithAuthor

	err := r.buildBaseBookWithAuthorQuery().
		Where("books.author_id = ?", authorID).
		Find(&booksWithAuthor).Error
	if err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}

	return booksWithAuthor, nil
}

func (r *bookRepository) FindByID(ID uint) (*model.BookWithAuthor, error) {
	var bookWithAuthor model.BookWithAuthor

	err := r.buildBaseBookWithAuthorQuery().
		Where("books.id = ?", ID).
		First(&bookWithAuthor).Error
	if err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}

	return &bookWithAuthor, nil
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
			a.fullname AS author_fullname,
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

func (r *bookRepository) buildBaseBookWithAuthorQuery() *gorm.DB {
	return r.db.
		Model(&model.Book{}).
		Select("books.id, books.author_id, authors.fullname AS author_fullname, books.title, books.year, books.category").
		Joins("LEFT JOIN authors ON books.author_id = authors.id")
}

func sanitizeListBooksParams(sort, order string) (string, string) {
	allowedSortColumnValues := []string{"title", "year", "category"}

	sort = strings.ToLower(sort)
	if !slices.Contains(allowedSortColumnValues, sort) {
		sort = "title"
	}

	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}
	return sort, order
}

package postgres

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"

	postgresInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/postgres"
	"gorm.io/gorm"
)

type BookRepository interface {
	WithinTX(tx *gorm.DB) BookRepository
	GetCategories(ctx context.Context) ([]string, error)
	GetNew(ctx context.Context) ([]model.BookWithAuthor, error)
	GetBooksByIDs(ctx context.Context, bookIDs []string) ([]model.BookWithAuthor, error)
	GetBooksByAuthorID(ctx context.Context, authorID uint) ([]model.BookWithAuthor, error)
	FindByID(ctx context.Context, bookID uint) (*model.BookWithAuthor, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, book *model.Book) error
	Delete(ctx context.Context, bookID uint) error
	ListByAuthorName(ctx context.Context, authorName string, page, count uint, sort, order string) ([]model.BookWithAuthor, error)
	ListByTitle(ctx context.Context, title string, page, count uint, sort, order string) ([]model.BookWithAuthor, error)
	ListByAuthorNameAndTitle(ctx context.Context, authorName, title string, page, count uint, sort, order string) ([]model.BookWithAuthor, error)
	ListByCategory(ctx context.Context, categoryName string, page, count uint, sort, order string) ([]model.BookWithAuthor, error)
}

type bookRepository struct {
	name    string
	timeout time.Duration
	db      *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{name: "Book(s)", timeout: 500 * time.Millisecond, db: db}
}

func (r *bookRepository) WithinTX(tx *gorm.DB) BookRepository {
	return &bookRepository{name: "Books(s)", timeout: 500 * time.Millisecond, db: tx}
}

func (r *bookRepository) GetCategories(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var categories []string
	if err := r.db.WithContext(ctx).Model(&model.Book{}).Distinct().Order("category").Pluck("category", &categories).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return categories, nil
}

func (r *bookRepository) GetNew(ctx context.Context) ([]model.BookWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	const NEW_BOOKS_COUNT = 10

	var booksWithAuthor []model.BookWithAuthor

	err := r.buildBaseBookWithAuthorQuery(ctx).
		Order("books.created_at DESC").
		Limit(NEW_BOOKS_COUNT).
		Find(&booksWithAuthor).Error
	if err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}

	return booksWithAuthor, nil
}

func (r *bookRepository) GetBooksByIDs(ctx context.Context, bookIDs []string) ([]model.BookWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var booksWithAuthor []model.BookWithAuthor

	err := r.buildBaseBookWithAuthorQuery(ctx).
		Where("books.id IN ?", bookIDs).
		Find(&booksWithAuthor).Error
	if err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}

	return booksWithAuthor, nil
}

func (r *bookRepository) GetBooksByAuthorID(ctx context.Context, authorID uint) ([]model.BookWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var booksWithAuthor []model.BookWithAuthor

	err := r.buildBaseBookWithAuthorQuery(ctx).
		Where("books.author_id = ?", authorID).
		Find(&booksWithAuthor).Error
	if err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}

	return booksWithAuthor, nil
}

func (r *bookRepository) FindByID(ctx context.Context, bookID uint) (*model.BookWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var bookWithAuthor model.BookWithAuthor

	err := r.buildBaseBookWithAuthorQuery(ctx).
		Where("books.id = ?", bookID).
		First(&bookWithAuthor).Error
	if err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}

	return &bookWithAuthor, nil
}

func (r *bookRepository) Count(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var booksCount int64
	if err := r.db.WithContext(ctx).Model(&model.Book{}).Count(&booksCount).Error; err != nil {
		return 0, postgresInfrastructure.NewError(err, r.name)
	}
	return booksCount, nil
}

func (r *bookRepository) Create(ctx context.Context, book *model.Book) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	book.Category = strings.ToLower(book.Category)
	if err := r.db.WithContext(ctx).Create(book).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}

func (r *bookRepository) Delete(ctx context.Context, bookID uint) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if err := r.db.WithContext(ctx).Delete(&model.Book{}, bookID).Error; err != nil {
		return postgresInfrastructure.NewError(err, r.name)
	}
	return nil
}

func (r *bookRepository) ListByAuthorName(ctx context.Context, authorName string, page, count uint, sort, order string) ([]model.BookWithAuthor, error) {
	return r.listBooksBy(ctx, map[string]string{"author": authorName}, page, count, sort, order)
}

func (r *bookRepository) ListByTitle(ctx context.Context, title string, page, count uint, sort, order string) ([]model.BookWithAuthor, error) {
	return r.listBooksBy(ctx, map[string]string{"title": title}, page, count, sort, order)
}

func (r *bookRepository) ListByAuthorNameAndTitle(ctx context.Context, authorName, title string, page, count uint, sort, order string) ([]model.BookWithAuthor, error) {
	return r.listBooksBy(ctx, map[string]string{"author": authorName, "title": title}, page, count, sort, order)
}

func (r *bookRepository) ListByCategory(ctx context.Context, category string, page, count uint, sort, order string) ([]model.BookWithAuthor, error) {
	return r.listBooksBy(ctx, map[string]string{"category": category}, page, count, sort, order)
}

func (r *bookRepository) listBooksBy(ctx context.Context, filters map[string]string, page, count uint, sort, order string) ([]model.BookWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

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
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&bookModels).Error; err != nil {
		return nil, postgresInfrastructure.NewError(err, r.name)
	}
	return bookModels, nil
}

func (r *bookRepository) buildBaseBookWithAuthorQuery(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).
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

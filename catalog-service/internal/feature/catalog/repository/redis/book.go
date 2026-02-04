package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/redis/model"
	redisInfrastructure "github.com/Yarik7610/library-backend/catalog-service/internal/infrastructure/storage/redis"
	"github.com/redis/go-redis/v9"
)

const (
	CATEGORIES_KEY            = "categories"
	CATEGORIES_KEY_EXPIRATION = time.Minute * 15

	NEW_BOOKS_KEY            = "books:new"
	NEW_BOOKS_KEY_EXPIRATION = time.Minute * 15

	POPULAR_BOOKS_KEY   = "books:popular"
	POPULAR_BOOKS_COUNT = 10
)

type BookRepository interface {
	SetCategories(categories []string) error
	GetCategories() ([]string, error)
	SetNewBooks(newBooks []model.Book) error
	GetNewBooks() ([]model.Book, error)
	UpdateBookViewsCount(bookID, userID uint) error
	GetBookViewsCount(bookID uint) (int64, error)
	GetPopularBooksIDs() ([]string, error)
}

type bookRepository struct {
	rdb *redis.Client
}

func NewBookRepository(rdb *redis.Client) BookRepository {
	return &bookRepository{rdb: rdb}
}

func (r *bookRepository) SetCategories(categories []string) error {
	ctx := context.Background()

	if err := r.rdb.Del(ctx, CATEGORIES_KEY).Err(); err != nil {
		return redisInfrastructure.NewError(err)
	}

	if len(categories) > 0 {
		if err := r.rdb.RPush(ctx, CATEGORIES_KEY, stringSliceToAnySlice(categories)...).Err(); err != nil {
			return redisInfrastructure.NewError(err)
		}
		if err := r.rdb.Expire(ctx, CATEGORIES_KEY, CATEGORIES_KEY_EXPIRATION).Err(); err != nil {
			return redisInfrastructure.NewError(err)
		}
	}
	return nil
}

func (r *bookRepository) GetCategories() ([]string, error) {
	ctx := context.Background()

	categories, err := r.rdb.LRange(ctx, CATEGORIES_KEY, 0, -1).Result()
	if err != nil {
		if redisInfrastructure.IsNil(err) {
			return nil, nil
		}
		return nil, redisInfrastructure.NewError(err)
	}
	return categories, nil
}

func (r *bookRepository) SetNewBooks(newBooks []model.Book) error {
	ctx := context.Background()

	newBooksBytes, err := json.Marshal(newBooks)
	if err != nil {
		return err
	}

	if err := r.rdb.Set(ctx, NEW_BOOKS_KEY, newBooksBytes, NEW_BOOKS_KEY_EXPIRATION).Err(); err != nil {
		return redisInfrastructure.NewError(err)
	}
	return nil
}

func (r *bookRepository) GetNewBooks() ([]model.Book, error) {
	ctx := context.Background()

	newBooksString, err := r.rdb.Get(ctx, NEW_BOOKS_KEY).Result()
	if err != nil {
		if redisInfrastructure.IsNil(err) {
			return nil, nil
		}
		return nil, redisInfrastructure.NewError(err)
	}

	var newBooks []model.Book
	if err := json.Unmarshal([]byte(newBooksString), &newBooks); err != nil {
		return nil, err
	}
	return newBooks, nil
}

func (r *bookRepository) UpdateBookViewsCount(bookID, userID uint) error {
	ctx := context.Background()

	bookViewsCountKey := fmt.Sprintf("books:%d:views", bookID)
	addedCount, err := r.rdb.PFAdd(ctx, bookViewsCountKey, userID).Result()
	if err != nil {
		return redisInfrastructure.NewError(err)
	}

	if addedCount > 0 {
		if err := r.rdb.ZIncrBy(ctx, POPULAR_BOOKS_KEY, 1, strconv.Itoa(int(bookID))).Err(); err != nil {
			return redisInfrastructure.NewError(err)
		}
	}
	return nil
}

func (r *bookRepository) GetBookViewsCount(bookID uint) (int64, error) {
	ctx := context.Background()

	bookViewsCountKey := fmt.Sprintf("books:%d:views", bookID)
	bookViewsCount, err := r.rdb.PFCount(ctx, bookViewsCountKey).Result()
	if err != nil {
		if redisInfrastructure.IsNil(err) {
			return 0, nil
		}
		return 0, redisInfrastructure.NewError(err)
	}
	return bookViewsCount, nil
}

func (r *bookRepository) GetPopularBooksIDs() ([]string, error) {
	ctx := context.Background()

	popularBookIDs, err := r.rdb.ZRevRange(ctx, POPULAR_BOOKS_KEY, 0, POPULAR_BOOKS_COUNT-1).Result()
	if err != nil {
		if redisInfrastructure.IsNil(err) {
			return nil, nil
		}
		return nil, redisInfrastructure.NewError(err)
	}
	return popularBookIDs, nil
}

func stringSliceToAnySlice(slice []string) []any {
	res := make([]any, len(slice))
	for i, v := range slice {
		res[i] = v
	}
	return res
}

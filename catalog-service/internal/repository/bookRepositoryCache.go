package repository

import (
	"context"
	"time"

	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"github.com/redis/go-redis/v9"
)

const (
	CATEGORIES_KEY            = "categories"
	CATEGORIES_KEY_EXPIRATION = time.Minute * 15

	NEW_BOOKS_KEY            = "books:new"
	NEW_BOOKS_KEY_EXPIRATION = time.Minute * 15

	POPULAR_BOOKS_KEY            = "books:popular"
	POPULAR_BOOKS_KEY_EXPIRATION = time.Minute * 15
)

type BookRepositoryCache interface {
	SetCategories(categories []string)
	GetCategories() ([]string, error)
	SetNewBooks(newBooks []model.Book)
	GetNewBooks() ([]model.Book, error)
	SetPopularBooks(popularBooks []model.Book)
	GetPopularBooks() ([]model.Book, error)
}

type bookRepositoryCache struct {
	rdb *redis.Client
}

func NewBookRepositoryCache(rdb *redis.Client) BookRepositoryCache {
	return &bookRepositoryCache{rdb: rdb}
}

func (r *bookRepositoryCache) SetCategories(categories []string) {
	ctx := context.Background()
	r.rdb.Del(ctx, CATEGORIES_KEY)
	if len(categories) > 0 {
		r.rdb.RPush(ctx, CATEGORIES_KEY, stringSliceToAnySlice(categories)...)
		r.rdb.Expire(ctx, CATEGORIES_KEY, CATEGORIES_KEY_EXPIRATION)
	}
}

func (r *bookRepositoryCache) GetCategories() ([]string, error) {
	ctx := context.Background()
	return r.rdb.LRange(ctx, CATEGORIES_KEY, 0, -1).Result()
}

func (r *bookRepositoryCache) SetNewBooks(newBooks []model.Book) {

}

func (r *bookRepositoryCache) GetNewBooks() ([]model.Book, error) {
	return nil, nil
}

func (r *bookRepositoryCache) SetPopularBooks(newBooks []model.Book) {

}

func (r *bookRepositoryCache) GetPopularBooks() ([]model.Book, error) {
	return nil, nil
}

func stringSliceToAnySlice(slice []string) []any {
	res := make([]any, len(slice))
	for i, v := range slice {
		res[i] = v
	}
	return res
}

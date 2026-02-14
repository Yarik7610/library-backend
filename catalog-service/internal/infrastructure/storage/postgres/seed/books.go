package seed

import (
	"context"
	"fmt"
	"sync"

	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres"
	"github.com/Yarik7610/library-backend/catalog-service/internal/feature/catalog/repository/postgres/model"
)

func Books(bookRepository postgres.BookRepository, pageRepository postgres.PageRepository, authorRepository postgres.AuthorRepository) error {
	ctx := context.Background()

	bookCount, err := bookRepository.Count(ctx)
	if err != nil {
		return err
	}

	if bookCount != 0 {
		return nil
	}

	return seedBooks(ctx, bookRepository, pageRepository, authorRepository)
}

func seedBooks(ctx context.Context, bookRepository postgres.BookRepository, pageRepository postgres.PageRepository, authorRepository postgres.AuthorRepository) error {
	const booksCount = 100
	const bookPagesCount = 5
	const workersCount = 10
	categories := []string{"Fantasy", "Mystery", "Romance", "Sci-Fi", "Thriller", "Horror", "Adventure", "Historical", "Biography", "Non-Fiction"}

	var wg sync.WaitGroup
	errors := make(chan error, booksCount)
	workers := make(chan struct{}, workersCount)

	for i := range booksCount {
		workers <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() {
				wg.Done()
				<-workers
			}()

			author := model.Author{
				Fullname: fmt.Sprintf("Author %d", i+1),
			}
			if err := authorRepository.Create(ctx, &author); err != nil {
				errors <- err
				return
			}

			book := model.Book{
				AuthorID: author.ID,
				Title:    fmt.Sprintf("Book %d", i+1),
				Year:     1900 + (i % 125),
				Category: categories[i%len(categories)],
			}
			if err := bookRepository.Create(ctx, &book); err != nil {
				errors <- err
				return
			}

			for p := 1; p <= bookPagesCount; p++ {
				page := model.Page{
					BookID:  book.ID,
					Number:  uint(p),
					Content: fmt.Sprintf("page%d", p),
				}
				if err := pageRepository.Create(ctx, &page); err != nil {
					errors <- err
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		return <-errors
	}

	return nil
}

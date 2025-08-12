package seed

import (
	"fmt"
	"sync"

	"github.com/Yarik7610/library-backend/catalog-service/internal/model"
	"github.com/Yarik7610/library-backend/catalog-service/internal/repository"
	"go.uber.org/zap"
)

func Books(bookRepository repository.BookRepository, pageRepository repository.PageRepository, authorRepository repository.AuthorRepository) {
	bookCount, err := bookRepository.CountBooks()
	if err != nil {
		zap.S().Fatalf("Failed to count books for seed need: %v", err)
	}

	if bookCount != 0 {
		zap.S().Info("Seeded books already exist, skip seeding")
		return
	}

	zap.S().Info("No books found, start seeding...")
	if err := seedBooks(bookRepository, pageRepository, authorRepository); err != nil {
		zap.S().Fatalf("Books seed error: %v", err)
	}
	zap.S().Info("Successfully seeded books")
}

func seedBooks(bookRepository repository.BookRepository, pageRepository repository.PageRepository, authorRepository repository.AuthorRepository) error {
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
			if err := authorRepository.CreateAuthor(&author); err != nil {
				errors <- err
				return
			}

			book := model.Book{
				AuthorID: author.ID,
				Title:    fmt.Sprintf("Book %d", i+1),
				Year:     1900 + (i % 125),
				Category: categories[i%len(categories)],
			}
			if err := bookRepository.CreateBook(&book); err != nil {
				errors <- err
				return
			}

			for p := 1; p <= bookPagesCount; p++ {
				page := model.Page{
					BookID:  book.ID,
					Number:  p,
					Content: fmt.Sprintf("page%d", p),
				}
				if err := pageRepository.CreatePage(&page); err != nil {
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

package model

import (
	"time"

	"github.com/Yarik7610/library-backend/catalog-service/internal/domain"
)

type Page struct {
	ID        uint `gorm:"primarykey"`
	BookID    uint `gorm:"uniqueIndex:book_id_number_index"`
	Number    uint `gorm:"uniqueIndex:book_id_number_index"`
	Content   string
	CreatedAt time.Time
}

func (p *Page) ToDomain() *domain.Page {
	return &domain.Page{
		ID:      p.ID,
		Number:  p.Number,
		Content: p.Content,
	}
}

package mocks

import (
	"time"

	"github.com/ryanbyrne30/htmx/web_app/internal/models"
)

var mockSnippet = &models.Snippet{
	ID: 1,
	Title: "An Old Silent Pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct {}

func (m *SnippetModel) Insert(title string, content string, expires time.Time) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil 
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{ mockSnippet }, nil
}
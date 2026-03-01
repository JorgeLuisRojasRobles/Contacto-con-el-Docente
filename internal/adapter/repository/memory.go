// Autonomo 2 Jorge Luis Rojas Robles - 2026
package repository

import (
	"sync"

	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/domain"
	"github.com/samber/mo"
)

type InMemoryBookRepo struct {
	mu    sync.RWMutex
	books map[domain.BookID]domain.Book
}

func NewInMemoryBookRepo() *InMemoryBookRepo {
	return &InMemoryBookRepo{
		books: make(map[domain.BookID]domain.Book),
	}
}

// FindByID retorna Option. Si no existe, retorna None, no nil ni error.
func (r *InMemoryBookRepo) FindByID(id domain.BookID) mo.Option[domain.Book] {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if book, exists := r.books[id]; exists {
		return mo.Some(book) // Envuelve el valor existente
	}
	return mo.None[domain.Book]() // Retorna explícitamente "Nada"
}

// Save guarda una copia del libro de manera segura.
func (r *InMemoryBookRepo) Save(book domain.Book) mo.Result[bool] {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.books[book.ID] = book
	return mo.Ok(true)
}

// ListAll (Extra para poder ver los libros guardados)
func (r *InMemoryBookRepo) ListAll() []domain.Book {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := make([]domain.Book, 0, len(r.books))
	for _, b := range r.books {
		books = append(books, b)
	}
	return books
}

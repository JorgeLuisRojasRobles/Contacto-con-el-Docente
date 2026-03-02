// Contacto con el Docente Jorge Luis Rojas Robles - 2026
package repository

import (
	"sync"

	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/domain"
	"github.com/samber/mo"
)

// Implemento la estructura del repositorio en memoria utilizando un mapa.
// Añado un RWMutex para garantizar que el acceso concurrente a los datos sea seguro, aplicando lo aprendido en la Unidad 4.
type InMemoryBookRepo struct {
	mu    sync.RWMutex
	books map[domain.BookID]domain.Book
}

// Constructor que inicializa el repositorio asegurándose de instanciar el mapa en memoria para evitar errores de tipo 'panic'.
func NewInMemoryBookRepo() *InMemoryBookRepo {
	return &InMemoryBookRepo{
		books: make(map[domain.BookID]domain.Book),
	}
}

// Desarrollé este método utilizando el patrón Option (mo.Option) en lugar de retornar punteros nulos.
// Esto me permite manejar la ausencia de un libro de forma segura a nivel de compilación.
func (r *InMemoryBookRepo) FindByID(id domain.BookID) mo.Option[domain.Book] {
	// Aplico un bloqueo de lectura para evitar lecturas sucias sin detener otras consultas simultáneas.
	r.mu.RLock()
	defer r.mu.RUnlock()

	if book, exists := r.books[id]; exists {
		// Envuelvo el valor existente en un Option válido.
		return mo.Some(book)
	}
	// Retorno explícitamente la ausencia de valor si el ID no coincide.
	return mo.None[domain.Book]()
}

// Este método persiste o actualiza el libro en el mapa de memoria.
func (r *InMemoryBookRepo) Save(book domain.Book) mo.Result[bool] {
	// Aplico un bloqueo de escritura completo porque estoy mutando el estado interno del repositorio, garantizando la integridad de los datos.
	r.mu.Lock()
	defer r.mu.Unlock()

	r.books[book.ID] = book
	return mo.Ok(true)
}

// Agregué este método auxiliar para extraer todo el catálogo actual del mapa y transformarlo en un Slice,
// facilitando su serialización posterior en los servicios web JSON.
func (r *InMemoryBookRepo) ListAll() []domain.Book {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := make([]domain.Book, 0, len(r.books))
	for _, b := range r.books {
		books = append(books, b)
	}
	return books
}

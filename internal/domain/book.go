// Autonomo 2 Jorge Luis Rojas Robles - 2026
package domain

import (
	"fmt"
	"time"

	"github.com/samber/mo"
)

type EPUBMetadata struct {
	Title  string
	Author string
}

type DomainError struct {
	Message string
	Code    int
}

func (e DomainError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewDomainError(msg string, code int) error {
	return DomainError{Message: msg, Code: code}
}

type BookRepository interface {
	FindByID(id BookID) mo.Option[Book]
	Save(book Book) mo.Result[bool]
	ListAll() []Book
}

type BookID string

type Book struct {
	ID          BookID
	Title       string
	Author      string
	Description mo.Option[string]
	PublishDate time.Time
	PageCount   int
	FilePath    string
	Format      string
	Status      string // 🟢 ¡NUEVO CAMPO PARA SABER SI ESTÁ PRESTADO!
}

func GenerateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Constructor Seguro
func NewBook(id string, title string, author string, pages int, path string) mo.Result[Book] {
	if id == "" {
		return mo.Err[Book](NewDomainError("ID no puede estar vacío", 400))
	}
	if title == "" {
		return mo.Err[Book](NewDomainError("El título es obligatorio", 400))
	}
	if pages <= 0 {
		return mo.Err[Book](NewDomainError("El conteo de páginas debe ser positivo", 400))
	}

	return mo.Ok(Book{
		ID:          BookID(id),
		Title:       title,
		Author:      author,
		Description: mo.None[string](),
		PublishDate: time.Now(),
		PageCount:   pages,
		FilePath:    path,
		Format:      "EPUB",
		Status:      "Disponible", // 🟢 Por defecto, todo libro nuevo está disponible
	})
}

// Setters Funcionales (Unidad 3)
func (b Book) WithTitle(newTitle string) Book {
	b.Title = newTitle
	return b
}

func (b Book) WithDescription(desc string) Book {
	b.Description = mo.Some(desc)
	return b
}

// Para prestar o devolver el libro de forma inmutable
func (b Book) WithStatus(newStatus string) Book {
	b.Status = newStatus
	return b
}

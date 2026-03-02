// Contacto con el Docente Jorge Luis Rojas Robles - 2026
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
	Status      string
	BorrowedBy  string // Agregué este campo para llevar el control exacto de qué usuario alquiló el libro.
}

func GenerateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Implemento un constructor seguro para validar las reglas de negocio antes de instanciar un libro en memoria.
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
		Status:      "Disponible",
		BorrowedBy:  "", // Lo inicializo vacío ya que un libro recién ingresado al sistema no tiene un arrendatario asignado.
	})
}

// Aplico el concepto de Setters Funcionales (Unidad 3) para proteger la encapsulación y garantizar la inmutabilidad de la estructura.
func (b Book) WithTitle(newTitle string) Book {
	b.Title = newTitle
	return b
}

func (b Book) WithDescription(desc string) Book {
	b.Description = mo.Some(desc)
	return b
}

// Este método me permite gestionar las transacciones (cambiar de Disponible a Prestado o Vendido) devolviendo una copia segura del objeto.
func (b Book) WithStatus(newStatus string) Book {
	b.Status = newStatus
	return b
}

// Desarrollo este método para asociar la identidad del cliente con el libro prestado en el momento de la transacción.
func (b Book) WithBorrowedBy(user string) Book {
	b.BorrowedBy = user
	return b
}

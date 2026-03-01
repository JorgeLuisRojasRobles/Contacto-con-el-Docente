// Autonomo 2 Jorge Luis Rojas Robles - 2026

package service

import (
	"fmt"

	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/adapter/epub"
	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/domain"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

type LibraryService struct {
	repo domain.BookRepository
}

func NewLibraryService(repo domain.BookRepository) *LibraryService {
	return &LibraryService{repo: repo}
}

func (s *LibraryService) ImportBooks(paths []string) []domain.Book {

	// 1. Pipeline de Transformación
	results := lo.Map(paths, func(path string, _ int) mo.Result[domain.Book] {

		// A. Ejecutamos el Parser
		resMeta := epub.ParseMetadata(path)

		// B. Si falló, devolvemos el error inmediatamente
		if resMeta.IsError() {
			return mo.Err[domain.Book](resMeta.Error())
		}

		// C. Si tuvo éxito, extraemos la data y creamos el Libro
		meta := resMeta.MustGet()

		return domain.NewBook(
			domain.GenerateID(),
			meta.Title,
			meta.Author,
			100,
			path,
		)
	})

	// 2. Filtrado de Errores
	validBooks := lo.FilterMap(results, func(res mo.Result[domain.Book], _ int) (domain.Book, bool) {

		val, err := res.Get()
		return val, err == nil
	})

	// 3. Persistencia (Usando la Interfaz)
	lo.ForEach(validBooks, func(book domain.Book, _ int) {
		s.repo.Save(book)
	})

	return validBooks
}

func (s *LibraryService) GetAllBooks() []domain.Book {
	return s.repo.ListAll()
}

// Agregar un libro manualmente desde la interfaz web
func (s *LibraryService) AddManualBook(title, author string) domain.Book {

	resBook := domain.NewBook(
		domain.GenerateID(),
		title,
		author,
		100,
		"Registro Manual",
	)

	newBook := resBook.MustGet()

	s.repo.Save(newBook)

	return newBook
}

func (s *LibraryService) BorrowBook(id string) error {
	// 1. Buscamos el libro en la base de datos
	optBook := s.repo.FindByID(domain.BookID(id))

	// 2. Si la caja está vacía, devolvemos un error
	if optBook.IsAbsent() {
		return fmt.Errorf("error: libro no encontrado")
	}

	// 3. Extraemos el libro de la caja
	book := optBook.MustGet()

	// 4. No se puede prestar algo que ya está prestado
	if book.Status == "Prestado" {
		return fmt.Errorf("error: el libro ya se encuentra prestado")
	}

	// 5. Usamos tu Setter Funcional para cambiar el estado de forma segura e inmutable
	updatedBook := book.WithStatus("Prestado")

	// 6. Guardamos el libro actualizado (esto sobreescribe el anterior)
	s.repo.Save(updatedBook)

	return nil
}

// Busca el libro prestado y lo vuelve a marcar como disponible
func (s *LibraryService) ReturnBook(id string) error {
	optBook := s.repo.FindByID(domain.BookID(id))

	if optBook.IsAbsent() {
		return fmt.Errorf("error: libro no encontrado")
	}

	book := optBook.MustGet()

	if book.Status == "Disponible" {
		return fmt.Errorf("error: el libro ya está en la biblioteca")
	}

	updatedBook := book.WithStatus("Disponible")
	s.repo.Save(updatedBook)

	return nil
}

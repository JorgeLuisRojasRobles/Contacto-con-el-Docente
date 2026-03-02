// Contacto con el Docente Jorge Luis Rojas Robles - 2026

package service

import (
	"fmt"

	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/adapter/epub"
	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/domain"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

// Defino el servicio principal de la biblioteca. Utilizo inyección de dependencias
// recibiendo la interfaz del repositorio para mantener el dominio de negocio desacoplado de la persistencia.
type LibraryService struct {
	repo domain.BookRepository
}

func NewLibraryService(repo domain.BookRepository) *LibraryService {
	return &LibraryService{repo: repo}
}

// Implemento la importación masiva de libros utilizando los principios de programación funcional.
func (s *LibraryService) ImportBooks(paths []string) []domain.Book {
	// Fase 1: Mapeo cada ruta de archivo a un Result de un libro, delegando la extracción de metadatos al adaptador.
	results := lo.Map(paths, func(path string, _ int) mo.Result[domain.Book] {
		resMeta := epub.ParseMetadata(path)
		if resMeta.IsError() {
			return mo.Err[domain.Book](resMeta.Error())
		}
		meta := resMeta.MustGet()

		// Instancio el libro a través del constructor seguro del dominio.
		return domain.NewBook(
			domain.GenerateID(),
			meta.Title,
			meta.Author,
			100,
			path,
		)
	})

	// Fase 2: Filtro los resultados, descartando silenciosamente los archivos que generaron error durante el parseo de metadatos.
	validBooks := lo.FilterMap(results, func(res mo.Result[domain.Book], _ int) (domain.Book, bool) {
		val, err := res.Get()
		return val, err == nil
	})

	// Fase 3: Persisto los libros válidos en el repositorio.
	lo.ForEach(validBooks, func(book domain.Book, _ int) {
		s.repo.Save(book)
	})

	return validBooks
}

func (s *LibraryService) GetAllBooks() []domain.Book {
	return s.repo.ListAll()
}

// Desarrollo este método para permitir el registro manual de libros desde la interfaz web.
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

// Gestiono la transacción de alquiler. Requiere el ID del libro y el nombre del usuario para fines de control y auditoría.
func (s *LibraryService) BorrowBook(id string, username string) error {
	optBook := s.repo.FindByID(domain.BookID(id))

	if optBook.IsAbsent() {
		return fmt.Errorf("error: libro no encontrado")
	}

	book := optBook.MustGet()

	// Valido las reglas de negocio fundamentales: el libro debe estar disponible para poder prestarse.
	if book.Status == "Prestado" {
		return fmt.Errorf("error: el libro ya se encuentra prestado")
	}
	if book.Status == "Vendido" {
		return fmt.Errorf("error: este libro ya fue vendido")
	}

	// Aplico los setters funcionales para mutar el estado de forma inmutable y registro al usuario responsable.
	updatedBook := book.WithStatus("Prestado").WithBorrowedBy(username)
	s.repo.Save(updatedBook)

	return nil
}

// Proceso la devolución del libro, integrando una capa de control de acceso basada en el usuario actual.
func (s *LibraryService) ReturnBook(id string, username string) error {
	optBook := s.repo.FindByID(domain.BookID(id))

	if optBook.IsAbsent() {
		return fmt.Errorf("error: libro no encontrado")
	}

	book := optBook.MustGet()

	if book.Status == "Disponible" {
		return fmt.Errorf("error: el libro ya está en la biblioteca")
	}

	// Regla de seguridad y privacidad: Valido que únicamente el usuario que originó el préstamo,
	// o el administrador central del sistema, tengan autorización para registrar la devolución.
	if book.BorrowedBy != username && username != "jorge" {
		return fmt.Errorf("solo el usuario que lo alquiló (%s) puede devolverlo", book.BorrowedBy)
	}

	// Libero el libro y limpio el registro del usuario.
	updatedBook := book.WithStatus("Disponible").WithBorrowedBy("")
	s.repo.Save(updatedBook)

	return nil
}

// Implemento la transacción de compra, la cual opera como un estado final dentro del ciclo de vida del libro.
func (s *LibraryService) BuyBook(id string) error {
	optBook := s.repo.FindByID(domain.BookID(id))

	if optBook.IsAbsent() {
		return fmt.Errorf("error: libro no encontrado")
	}

	book := optBook.MustGet()

	// Valido restricciones de negocio cruzadas para evitar corromper el inventario.
	if book.Status == "Vendido" {
		return fmt.Errorf("error: este libro ya fue vendido")
	}
	if book.Status == "Prestado" {
		return fmt.Errorf("error: no puedes comprar un libro prestado")
	}

	// Actualizo el estado a vendido de forma inmutable y lo persisto en memoria.
	updatedBook := book.WithStatus("Vendido")
	s.repo.Save(updatedBook)

	return nil
}

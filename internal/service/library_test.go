// Autonomo 2 Jorge Luis Rojas Robles - 2026
package service

import (
	"testing"

	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/domain"
	"github.com/samber/mo"
)

// =============================================================================
// 1. MOCK (Simulacro) del Repositorio - TEMA UNIDAD 3: INTERFACES
// =============================================================================
// Este struct "finge" ser una base de datos. Sirve para engañar al servicio
// durante el test y verificar si llama a guardar o no.
type MockRepository struct {
	SaveCalls int // Contador de cuántas veces se llamó a Save
}

// Implementación de la Interfaz domain.BookRepository
func (m *MockRepository) FindByID(id domain.BookID) mo.Option[domain.Book] {
	return mo.None[domain.Book]()
}

func (m *MockRepository) Save(book domain.Book) mo.Result[bool] {
	m.SaveCalls++ // ¡Espía! Contamos la llamada
	return mo.Ok(true)
}

func (m *MockRepository) ListAll() []domain.Book {
	return []domain.Book{}
}

// =============================================================================
// 2. PRUEBA UNITARIA - TEMA SECCIÓN 9 PDF
// =============================================================================

func TestImportBooks_FiltradoDeErrores(t *testing.T) {
	// A. PREPARACIÓN (Arrange)
	// Creamos el mock
	mockRepo := &MockRepository{}

	// Inyectamos el mock en lugar del repositorio real (Polimorfismo)
	svc := NewLibraryService(mockRepo)

	// Definimos entradas que SABEMOS que van a fallar (porque no existen los archivos)
	rutasInvalidas := []string{
		"archivo_inexistente_1.epub",
		"archivo_inexistente_2.epub",
	}

	// B. EJECUCIÓN (Act)
	// Ejecutamos la lógica funcional
	librosProcesados := svc.ImportBooks(rutasInvalidas)

	// C. ASERCIÓN (Assert)

	// 1. Verificamos Robustez: El sistema NO debe explotar (panic) y debe devolver una lista vacía.
	// Esto confirma que lo.FilterMap descartó los errores de la Vía Roja.
	if len(librosProcesados) != 0 {
		t.Errorf("Se esperaba 0 libros válidos (todos son errores), pero llegaron %d", len(librosProcesados))
	}

	// 2. Verificamos Comportamiento: No se debió intentar guardar nada en el repo.
	if mockRepo.SaveCalls != 0 {
		t.Errorf("El repositorio no debió ser llamado, pero SaveCalls = %d", mockRepo.SaveCalls)
	}
}

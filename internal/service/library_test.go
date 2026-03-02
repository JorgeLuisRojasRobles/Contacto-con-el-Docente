// Contacto con el Docente Jorge Luis Rojas Robles - 2026
package service

import (
	"testing"

	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/domain"
	"github.com/samber/mo"
)

// Implemento un Mock (simulación) del repositorio para aislar las pruebas unitarias.
// Esto demuestra la utilidad práctica de las Interfaces que estudiamos en la Unidad 3,
// permitiéndome desacoplar la lógica de negocio de la base de datos real.
type MockRepository struct {
	SaveCalls int // Utilizo este contador para auditar cuántas veces el servicio intenta guardar datos.
}

// Cumplo con el contrato de la interfaz domain.BookRepository retornando valores vacíos controlados.
func (m *MockRepository) FindByID(id domain.BookID) mo.Option[domain.Book] {
	return mo.None[domain.Book]()
}

// Sobrescribo el método Save para que funcione como un espía de comportamiento durante mis pruebas.
func (m *MockRepository) Save(book domain.Book) mo.Result[bool] {
	m.SaveCalls++
	return mo.Ok(true)
}

func (m *MockRepository) ListAll() []domain.Book {
	return []domain.Book{}
}

// Desarrollo esta prueba unitaria para validar la robustez del sistema frente a errores,
// aplicando los conceptos de Testing correspondientes a la última unidad.
func TestImportBooks_FiltradoDeErrores(t *testing.T) {

	// Fase de Preparación (Arrange):
	// Instancio el mock y lo inyecto en el servicio de biblioteca, aprovechando el polimorfismo.
	mockRepo := &MockRepository{}
	svc := NewLibraryService(mockRepo)

	// Defino un escenario de falla controlada proporcionando rutas de archivos que no existen.
	rutasInvalidas := []string{
		"archivo_inexistente_1.epub",
		"archivo_inexistente_2.epub",
	}

	// Fase de Ejecución (Act):
	// Pongo a prueba la lógica funcional de importación.
	librosProcesados := svc.ImportBooks(rutasInvalidas)

	// Fase de Aserción (Assert):

	// 1. Verifico la robustez del código. El sistema no debe colapsar y debe retornar una lista vacía.
	// Esto me confirma que la función lo.FilterMap filtró exitosamente los errores de lectura.
	if len(librosProcesados) != 0 {
		t.Errorf("Se esperaba 0 libros válidos debido a errores de lectura, pero se procesaron %d", len(librosProcesados))
	}

	// 2. Verifico el comportamiento de la persistencia. Al no haber libros válidos procesados,
	// el servicio no debió interactuar con la base de datos simulada en ningún momento.
	if mockRepo.SaveCalls != 0 {
		t.Errorf("El repositorio no debió registrar llamadas a Save, pero se detectaron %d llamadas", mockRepo.SaveCalls)
	}
}

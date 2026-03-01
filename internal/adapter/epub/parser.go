// Autonomo 2 Jorge Luis Rojas Robles - 2026
package epub

import (
	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/domain"
	"github.com/kapmahc/epub"
	"github.com/samber/mo"
)

// ParseMetadata envuelve la librería externa en una función pura que retorna Result.
func ParseMetadata(path string) mo.Result[domain.EPUBMetadata] {
	// mo.Try atrapa errores y los convierte en Result (Rieles de tren)
	return mo.Try(func() (domain.EPUBMetadata, error) {

		// 1. Abrir el archivo usando la librería kapmahc/epub
		book, err := epub.Open(path)
		if err != nil {
			return domain.EPUBMetadata{}, err
		}
		// Cerramos el archivo al terminar
		defer book.Close()

		// 2. Extraer Título (La librería devuelve un array de títulos, tomamos el primero)
		title := "Desconocido"
		if len(book.Opf.Metadata.Title) > 0 {
			title = book.Opf.Metadata.Title[0]
		}

		// 3. Extraer Autor (La librería devuelve un array de creadores)
		author := "Anónimo"
		if len(book.Opf.Metadata.Creator) > 0 {
			// El campo 'Data' contiene el nombre del autor
			author = book.Opf.Metadata.Creator[0].Data
		}

		// 4. Retornar nuestro objeto de dominio limpio
		return domain.EPUBMetadata{
			Title:  title,
			Author: author,
		}, nil
	})
}

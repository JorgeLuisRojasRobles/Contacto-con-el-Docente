// Contacto con el Docente Jorge Luis Rojas Robles - 2026
package epub

import (
	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/domain"
	"github.com/kapmahc/epub"
	"github.com/samber/mo"
)

// Encapsulo la lógica de la librería externa dentro de una función pura que retorna un Result, cumpliendo con los principios de abstracción.
func ParseMetadata(path string) mo.Result[domain.EPUBMetadata] {
	// Utilizo mo.Try para capturar cualquier error de lectura y convertirlo en un Result, manejando el flujo de errores de forma segura.
	return mo.Try(func() (domain.EPUBMetadata, error) {

		// Intento abrir el archivo físico utilizando la librería de terceros.
		book, err := epub.Open(path)
		if err != nil {
			return domain.EPUBMetadata{}, err
		}
		// Me aseguro de liberar el recurso en memoria al finalizar la lectura mediante la directiva defer.
		defer book.Close()

		// Extraigo el título. Como la metadata de EPUB puede contener múltiples títulos, valido su longitud y tomo el primero por defecto.
		title := "Desconocido"
		if len(book.Opf.Metadata.Title) > 0 {
			title = book.Opf.Metadata.Title[0]
		}

		// Realizo el mismo proceso de validación para el autor, accediendo al campo de datos del primer creador registrado.
		author := "Anónimo"
		if len(book.Opf.Metadata.Creator) > 0 {
			author = book.Opf.Metadata.Creator[0].Data
		}

		// Finalmente, devuelvo la estructura de dominio limpia, desacoplando la lógica central de mi proyecto de la librería externa.
		return domain.EPUBMetadata{
			Title:  title,
			Author: author,
		}, nil
	})
}

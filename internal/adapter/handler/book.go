// Contacto con el Docente Jorge Luis Rojas Robles - 2026
package handler

import (
	"net/http"

	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/domain"
	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/samber/mo"
)

// Defino el manejador (Handler) para conectar las peticiones HTTP con la lógica de negocio de los libros.
type BookHandler struct {
	service *service.LibraryService
}

// Constructor para inyectar el servicio de biblioteca en el manejador.
func NewBookHandler(s *service.LibraryService) *BookHandler {
	return &BookHandler{service: s}
}

// Implemento el método Import para procesar las peticiones POST de importación masiva de archivos.
func (h *BookHandler) Import(c echo.Context) error {

	// Creo un DTO (Data Transfer Object) interno para mapear y validar estrictamente la estructura del JSON entrante.
	type Request struct {
		Paths []string `json:"paths"`
	}

	// Fase de Binding: Utilizo la librería 'mo' para capturar de forma segura la entrada impura (HTTP)
	// y transformarla en una estructura tipada en Go, evitando que el sistema colapse si el JSON viene malformado.
	inputResult := mo.Try(func() (Request, error) {
		var req Request
		if err := c.Bind(&req); err != nil {
			return Request{}, err
		}
		return req, nil
	})

	// Declaro una variable de tipo Result para encapsular la respuesta de la lógica de negocio pura.
	var processingResult mo.Result[[]domain.Book]

	if inputResult.IsError() {
		// Si la validación de los datos de entrada falla, propago el error hacia la capa de respuesta.
		processingResult = mo.Err[[]domain.Book](inputResult.Error())
	} else {
		// Si la entrada es válida, extraigo las rutas y delego el procesamiento intensivo al servicio de dominio.
		req := inputResult.MustGet()
		books := h.service.ImportBooks(req.Paths)
		processingResult = mo.Ok(books)
	}

	// Fase de Respuesta: Gestiono la salida HTTP evaluando el estado del Result final.
	// Opto por una verificación explícita de errores para garantizar la estabilidad del compilador.

	if processingResult.IsError() {
		// Flujo de error (Bad Path): Retorno un código HTTP 400 estructurado en formato JSON para el cliente.
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":  "Solicitud inválida",
			"detail": processingResult.Error().Error(),
		})
	}

	// Flujo de éxito (Happy Path): Retorno un código HTTP 200, serializando la lista de libros procesados.
	books := processingResult.MustGet()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  books,
		"count": len(books),
	})
}

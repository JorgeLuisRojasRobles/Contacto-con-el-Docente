// Autonomo 2 Jorge Luis Rojas Robles - 2026
package handler

import (
	"net/http"

	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/domain"
	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/samber/mo"
)

type BookHandler struct {
	service *service.LibraryService
}

func NewBookHandler(s *service.LibraryService) *BookHandler {
	return &BookHandler{service: s}
}

// Import maneja POST /books/import
func (h *BookHandler) Import(c echo.Context) error {
	// Definición de DTO (Data Transfer Object) para la petición
	type Request struct {
		Paths []string `json:"paths"`
	}

	// 1. Binding (Entrada Impura -> Estructura Segura)
	inputResult := mo.Try(func() (Request, error) {
		var req Request
		if err := c.Bind(&req); err != nil {
			return Request{}, err
		}
		return req, nil
	})

	// 2. Procesamiento (Lógica Pura)
	var processingResult mo.Result[[]domain.Book]

	if inputResult.IsError() {
		// Si falló la entrada, propagamos el error
		processingResult = mo.Err[[]domain.Book](inputResult.Error())
	} else {
		// Si todo bien, ejecutamos el servicio
		req := inputResult.MustGet()
		books := h.service.ImportBooks(req.Paths)
		processingResult = mo.Ok(books)
	}

	// 3. Respuesta (Salida -> HTTP) - CORREGIDO
	// Usamos verificación explícita para evitar errores de compilación con 'Match'
	if processingResult.IsError() {
		// Caso Error (Red Track)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":  "Solicitud inválida",
			"detail": processingResult.Error().Error(),
		})
	}

	// Caso Éxito (Green Track)
	books := processingResult.MustGet()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  books,
		"count": len(books),
	})
}

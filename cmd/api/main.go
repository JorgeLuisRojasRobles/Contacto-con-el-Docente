// Autonomo 2 Jorge Luis Rojas Robles - 2026
package main

import (
	"fmt"

	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/adapter/handler"
	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/adapter/repository"
	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/service"
	"github.com/labstack/echo/v4"
)

func main() {
	repo := repository.NewInMemoryBookRepo()
	svc := service.NewLibraryService(repo)
	bookHandler := handler.NewBookHandler(svc)

	e := echo.New()

	// --- 1. RUTAS DE LA INTERFAZ WEB ---
	e.GET("/", func(c echo.Context) error {
		return c.File("views/login.html")
	})

	e.POST("/login", func(c echo.Context) error {
		user := c.FormValue("username")
		pass := c.FormValue("password")
		if user == "jorge" && pass == "1234" {
			return c.Redirect(302, "/dashboard")
		}
		return c.String(401, "Acceso Denegado: Credenciales inválidas")
	})

	e.GET("/dashboard", func(c echo.Context) error {
		return c.File("views/dashboard.html")
	})

	// 2 Recibe los datos del botón "Guardar Libro"
	e.POST("/books/manual", func(c echo.Context) error {
		var input struct {
			Title  string `json:"title"`
			Author string `json:"author"`
		}

		if err := c.Bind(&input); err != nil {
			fmt.Println("Error leyendo formulario:", err)
			return c.String(400, "Error en los datos")
		}

		fmt.Printf("\n🟢 Libro guardado-> Libro: '%s' | Autor: '%s'\n", input.Title, input.Author)

		svc.AddManualBook(input.Title, input.Author)

		return c.String(200, "OK")
	})

	e.GET("/books", func(c echo.Context) error {
		books := svc.GetAllBooks()
		return c.JSON(200, books)
	})

	e.POST("/books/import", bookHandler.Import)

	// --- 3. ARRANCAR SERVIDOR ---
	fmt.Println("Servidor de Autonomo 2 corriendo en el puerto 8080")
	e.Logger.Fatal(e.Start(":8080"))
}

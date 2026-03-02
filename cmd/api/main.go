// Contacto con el Docente Jorge Luis Rojas Robles - 2026
package main

import (
	"fmt"
	"net/http"

	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/adapter/handler"
	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/adapter/repository"
	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/domain"
	"github.com/JorgeLuisRojasRobles/Autonomo-2/internal/service"
	"github.com/labstack/echo/v4"
)

// Implemento un mapa en memoria para simular una base de datos de usuarios de forma ágil para este proyecto.
var usersDB = map[string]string{
	"cliente": "1234",
}

func main() {
	// Inicializo las dependencias aplicando el patrón de inyección de dependencias y separando responsabilidades.
	repo := repository.NewInMemoryBookRepo()
	svc := service.NewLibraryService(repo)
	bookHandler := handler.NewBookHandler(svc)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.File("views/login.html")
	})

	e.POST("/login", func(c echo.Context) error {
		user := c.FormValue("username")
		pass := c.FormValue("password")

		// Valido credenciales para el administrador principal del sistema.
		if user == "jorge" && pass == "1234" {
			c.SetCookie(&http.Cookie{Name: "role", Value: "jorge", Path: "/"})
			c.SetCookie(&http.Cookie{Name: "username", Value: user, Path: "/"}) // Guardo el nombre en la cookie para auditoría de acciones.
			return c.Redirect(302, "/dashboard")
		}

		// Valido credenciales contra la base de datos en memoria para los clientes regulares.
		if storedPass, exists := usersDB[user]; exists && storedPass == pass {
			c.SetCookie(&http.Cookie{Name: "role", Value: "cliente", Path: "/"})
			c.SetCookie(&http.Cookie{Name: "username", Value: user, Path: "/"})
			return c.Redirect(302, "/dashboard")
		}

		return c.String(401, "Acceso Denegado: Credenciales inválidas o usuario no existe")
	})

	e.GET("/dashboard", func(c echo.Context) error {
		return c.File("views/dashboard.html")
	})

	// Configuración de los Servicios Web para cumplir con los requerimientos de la Unidad 4.

	e.GET("/api/books", func(c echo.Context) error {
		books := svc.GetAllBooks()
		return c.JSON(200, books)
	})

	e.POST("/api/books/manual", func(c echo.Context) error {
		var input struct {
			Title  string `json:"title"`
			Author string `json:"author"`
		}
		if err := c.Bind(&input); err != nil {
			return c.JSON(400, map[string]string{"error": "Datos inválidos"})
		}
		svc.AddManualBook(input.Title, input.Author)
		return c.JSON(200, map[string]string{"message": "Libro creado exitosamente"})
	})

	e.POST("/api/books/import", bookHandler.Import)

	e.GET("/api/books/:id", func(c echo.Context) error {
		id := c.Param("id")
		optBook := repo.FindByID(domain.BookID(id))
		if optBook.IsAbsent() {
			return c.JSON(404, map[string]string{"error": "Libro no encontrado"})
		}
		return c.JSON(200, optBook.MustGet())
	})

	// Servicio para alquilar un libro. Extraigo la cookie del usuario para vincular el préstamo a su identidad.
	e.POST("/api/transactions/borrow/:id", func(c echo.Context) error {
		id := c.Param("id")
		cookie, _ := c.Cookie("username")
		if err := svc.BorrowBook(id, cookie.Value); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		return c.JSON(200, map[string]string{"message": "Libro prestado exitosamente"})
	})

	// Servicio para devolver un libro. Envío el usuario logueado al servicio para que valide si tiene los permisos adecuados.
	e.POST("/api/transactions/return/:id", func(c echo.Context) error {
		id := c.Param("id")
		cookie, _ := c.Cookie("username")
		if err := svc.ReturnBook(id, cookie.Value); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		return c.JSON(200, map[string]string{"message": "Libro devuelto exitosamente"})
	})

	e.POST("/api/transactions/buy/:id", func(c echo.Context) error {
		id := c.Param("id")
		if err := svc.BuyBook(id); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		return c.JSON(200, map[string]string{"message": "Libro comprado exitosamente"})
	})

	e.GET("/api/system/status", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "online", "version": "1.0.0", "author": "Jorge Rojas"})
	})

	// Implemento un endpoint adicional para el registro de nuevos usuarios en tiempo real.
	e.POST("/api/users/register", func(c echo.Context) error {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.Bind(&input); err != nil {
			return c.JSON(400, map[string]string{"error": "Datos inválidos"})
		}

		// Valido la disponibilidad del nombre de usuario antes de procesar el registro.
		if _, exists := usersDB[input.Username]; exists || input.Username == "jorge" {
			return c.JSON(400, map[string]string{"error": "El nombre ya está en uso"})
		}

		usersDB[input.Username] = input.Password
		fmt.Printf("Registro de sistema: Nuevo usuario incorporado (%s)\n", input.Username)
		return c.JSON(200, map[string]string{"message": "Usuario registrado exitosamente"})
	})

	fmt.Println("Iniciando el servidor del Proyecto Integrador en el puerto 8080...")
	e.Logger.Fatal(e.Start(":8080"))
}

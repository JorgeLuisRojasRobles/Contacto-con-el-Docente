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

var usersDB = map[string]string{
	"cliente": "1234",
}

func main() {
	repo := repository.NewInMemoryBookRepo()
	svc := service.NewLibraryService(repo)
	bookHandler := handler.NewBookHandler(svc)

	// ==========================================
	// PRECARGA DE DATOS (SEEDING) - LIBROS MUNDIALES
	// ==========================================
	svc.AddManualBook("Cien años de soledad", "Gabriel García Márquez")
	svc.AddManualBook("Don Quijote de la Mancha", "Miguel de Cervantes")
	svc.AddManualBook("El Principito", "Antoine de Saint-Exupéry")
	svc.AddManualBook("1984", "George Orwell")
	svc.AddManualBook("Harry Potter y la piedra filosofal", "J.K. Rowling")
	// ==========================================

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.File("views/login.html")
	})

	e.POST("/login", func(c echo.Context) error {
		user := c.FormValue("username")
		pass := c.FormValue("password")

		// 1. Rol: Administrador (Acceso total)
		if user == "jorge" && pass == "1234" {
			c.SetCookie(&http.Cookie{Name: "role", Value: "admin", Path: "/"})
			c.SetCookie(&http.Cookie{Name: "username", Value: user, Path: "/"})
			return c.Redirect(302, "/dashboard")
		}

		// 2. Rol: Visitante (Solo lectura)
		if user == "visitante" && pass == "1234" {
			c.SetCookie(&http.Cookie{Name: "role", Value: "visitante", Path: "/"})
			c.SetCookie(&http.Cookie{Name: "username", Value: user, Path: "/"})
			return c.Redirect(302, "/dashboard")
		}

		// 3. Rol: Cliente (Puede alquilar y comprar)
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

	e.POST("/api/transactions/borrow/:id", func(c echo.Context) error {
		id := c.Param("id")
		cookie, _ := c.Cookie("username")
		if err := svc.BorrowBook(id, cookie.Value); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		return c.JSON(200, map[string]string{"message": "Libro prestado exitosamente"})
	})

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

	e.POST("/api/users/register", func(c echo.Context) error {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.Bind(&input); err != nil {
			return c.JSON(400, map[string]string{"error": "Datos inválidos"})
		}

		if _, exists := usersDB[input.Username]; exists || input.Username == "jorge" || input.Username == "visitante" {
			return c.JSON(400, map[string]string{"error": "El nombre ya está en uso"})
		}

		usersDB[input.Username] = input.Password
		fmt.Printf("Registro de E-Commerce: Nuevo cliente incorporado (%s)\n", input.Username)
		return c.JSON(200, map[string]string{"message": "Cliente registrado exitosamente"})
	})

	fmt.Println("🛒 Iniciando el servidor del E-Commerce de Libros en el puerto 8080...")
	e.Logger.Fatal(e.Start(":8080"))
}

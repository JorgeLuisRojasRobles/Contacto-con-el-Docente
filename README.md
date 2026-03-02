# 🛒 E-Commerce de Libros (Sistema de Gestión y Transacciones RESTful)

**Autor:** Jorge Luis Rojas Robles  
**Institución:** Universidad Internacional del Ecuador (UIDE)  
**Asignatura:** Programación Orientada a Objetos 2  
**Fecha:** 01 de Marzo de 2026  

---

## 🎯 Objetivo del Programa
Desarrollar una plataforma web robusta para la gestión de un catálogo de libros electrónicos y físicos. El sistema permite administrar inventarios, registrar nuevos usuarios y procesar transacciones comerciales (alquiler y compra) aplicando un enfoque de programación funcional y arquitectura limpia en Go (Golang).

## ⚙️ Funcionalidades Principales
El software integra los conocimientos de las 4 unidades de la materia, destacando:
1. **Control de Acceso Basado en Roles (RBAC):** Gestión de sesiones mediante cookies para Administradores, Clientes y Visitantes.
2. **Persistencia Segura en Memoria:** Uso de Mutex (`sync.RWMutex`) para garantizar la integridad de los datos en entornos concurrentes.
3. **Programación Funcional y Mónadas:** Implementación de flujos seguros y manejo de errores determinista mediante la librería `samber/mo` (Result/Option) y procesamiento de colecciones con `samber/lo`.
4. **API RESTful (Servicios Web JSON):** Exposición de más de 8 endpoints para la comunicación entre el cliente (Frontend HTML/JS) y el servidor (Backend Go).

### 🌐 Servicios Web Implementados (Endpoints)
La serialización de datos se realiza estrictamente en formato **JSON**:
* `GET /api/books` - Lista todo el catálogo.
* `GET /api/books/:id` - Obtiene detalles de un libro específico.
* `POST /api/books/manual` - Crea un libro de forma manual.
* `POST /api/books/import` - Importación masiva extrayendo metadatos de archivos EPUB.
* `POST /api/transactions/borrow/:id` - Procesa el alquiler digital de un libro.
* `POST /api/transactions/return/:id` - Procesa la devolución (validando identidad).
* `POST /api/transactions/buy/:id` - Procesa la compra física (Agotar stock).
* `POST /api/users/register` - Registro asíncrono de nuevos clientes.
* `GET /api/system/status` - Healthcheck del sistema.

## 🚀 Instrucciones de Ejecución y Pruebas

Para levantar el servidor localmente, ejecute:
```bash
go run cmd/api/main.go

El sistema estará disponible en http://localhost:8080

# Makefile - Automatización del flujo de trabajo

# Variables
BINARY_NAME=server
MAIN_PATH=cmd/api/main.go

# Ejecuta air para desarrollo con recarga en caliente 
run:
	air

# Compila el binario optimizado para producción 
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Ejecuta pruebas unitarias con cobertura 
test:
	go test -cover ./...

# Limpia dependencias y formatea el código 
tidy:
	go mod tidy
	go fmt ./...
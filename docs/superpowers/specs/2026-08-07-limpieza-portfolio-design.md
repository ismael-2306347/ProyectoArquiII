# Diseño — Limpieza y mejoras para portfolio

**Fecha:** 2026-08-07
**Objetivo:** Dejar el repo presentable como pieza de portfolio/CV. README claro, arquitectura correcta, setup reproducible, config higienizada y estructura prolija. Idioma del README: español.

## Contexto

Sistema de reservas de hotel con arquitectura de microservicios:
- **Backend (Go + Gin):** users-api (8080), rooms-api (8081), reservations-api (8082), search-api (8083).
- **Datos:** MySQL (users, rooms), MongoDB (reservations), Apache Solr (search).
- **Infra transversal:** Memcached (caché), RabbitMQ (eventos: rooms publica → search consume).
- **Frontend:** React 18 + TypeScript + Vite + TailwindCSS + axios + react-router. No está en docker-compose; se corre con npm.

Fase 1 (ya completada): purga de secretos/binarios del historial, `.gitignore` raíz, force-push, borrado de ramas viejas.

## Alcance de esta fase

### A. README raíz (reescritura en español)
Reemplazar `README.md` (actualmente describe una estructura "buscada" irreal con `.js`/`bookings-api`/puertos incorrectos) por uno que refleje el código real:
- Título y descripción del sistema.
- Stack tecnológico.
- Diagrama de arquitectura corregido (4 microservicios + flujo de eventos rooms→search).
- Tabla de servicios y puertos.
- Endpoints reales por servicio (verificados en los `main.go`).
- Instrucciones de ejecución: `docker-compose up` (backend) + `cd frontend && npm install && npm run dev` (frontend).
- Referencia a variables de entorno (`.env.example`).
- Estructura real del proyecto.

### B. Config: credenciales a `.env`
- Crear `.env.example` en la raíz con todas las variables y valores de ejemplo.
- Modificar `docker-compose.yml` para usar `${VAR:-default}`, conservando defaults para que `docker-compose up` siga funcionando sin `.env`.
- `.env` real ya está ignorado por `.gitignore`.

### C. Limpieza liviana
- **Borrar docs redundantes (9):** `FRONTEND_QUICKSTART.md`, `SEARCH_API_INTEGRATION.md`, `frontend/BACKEND_INTEGRATION.md`, `frontend/CHANGELOG.md`, `frontend/INSTALL_TROUBLESHOOTING.md`, `frontend/README.md`, `services/search-api/FRONTEND_INTEGRATION.md`, `services/search-api/QUICKSTART.md`.
  - Conservar: `services/rooms-api/README.md`, `services/search-api/README.md`, `services/search-api/scripts/README.md`.
- **Borrar cruft de raíz:** `go.mod` (módulo fantasma "users-api"), `go.sum` (vacío), `package-lock.json` (stub vacío). Cada servicio tiene su propio `go.mod`.
- **Normalizar Dockerfiles:** renombrar los 4 `dockerfile` → `Dockerfile` y ajustar `docker-compose.yml` (corrige el build en Linux/CI, hoy roto por casing).
- **Corregir typo:** `services/rooms-api/config/rebbitmq.go` → `rabbitmq.go`.

### Fuera de alcance
Lógica de negocio, tests, refactors grandes. `POST /health` de reservations-api queda como está (opcional a futuro).

## Endpoints verificados (para el README)

- **users-api (8080):** `GET /users`, `POST /users`, `GET /users/:id`, `POST /login`.
- **rooms-api (8081):** `GET /health`; `GET /api/v1/rooms`, `GET /api/v1/rooms/available`, `GET /api/v1/rooms/number/:number`, `GET /api/v1/rooms/:id`; admin (auth+admin): `POST/PUT/PATCH/DELETE /api/v1/admin/rooms...`.
- **reservations-api (8082):** `POST /health`; `GET/POST /api/reservations`, `GET/DELETE /api/reservations/:id`, `GET /api/reservations/users/:user_id/myreservations`.
- **search-api (8083):** `GET /health`, `GET /api/search/rooms`.

## Criterio de éxito

- `docker-compose up` levanta todo el backend en Linux y Windows.
- README permite a un tercero entender y correr el proyecto sin leer el código.
- Sin archivos basura/duplicados en el árbol.

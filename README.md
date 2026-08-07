# 🏨 Sistema de Reservas de Hotel — Arquitectura de Microservicios

Aplicación web para la gestión y reserva de habitaciones de hotel, construida con una
**arquitectura de microservicios** en Go y un frontend en React. Los servicios se
comunican de forma **event-driven** mediante RabbitMQ y usan cachés (Memcached + caché
local) y un motor de búsqueda (Apache Solr) para escalar las consultas.

> Proyecto desarrollado para la materia **Arquitectura de Software II**.

---

## 🧰 Stack tecnológico

| Capa | Tecnologías |
|------|-------------|
| **Backend** | Go 1.25 · [Gin](https://gin-gonic.com/) |
| **Frontend** | React 18 · TypeScript · Vite · TailwindCSS · React Router · Axios |
| **Bases de datos** | MySQL 8 (users, rooms) · MongoDB 7 (reservations) |
| **Búsqueda** | Apache Solr 9 |
| **Caché** | Memcached · caché local en memoria |
| **Mensajería** | RabbitMQ (eventos) |
| **Infraestructura** | Docker · Docker Compose |
| **Auth** | JWT (roles: usuario / admin) |

---

## 🏗️ Arquitectura

```
                            ┌──────────────┐
                            │   Frontend   │  React + Vite (:3000)
                            └──────┬───────┘
                                   │  HTTP (proxy)
        ┌──────────────┬──────────┴──────────┬──────────────────┐
        │              │                      │                  │
 ┌──────▼──────┐ ┌─────▼──────┐      ┌────────▼────────┐ ┌───────▼──────┐
 │  users-api  │ │ rooms-api  │      │ reservations-api│ │  search-api  │
 │   (:8080)   │ │  (:8081)   │      │     (:8082)     │ │   (:8083)    │
 └──────┬──────┘ └─────┬──────┘      └────────┬────────┘ └───────┬──────┘
        │              │                      │                  │
 ┌──────▼──────┐ ┌─────▼──────┐      ┌────────▼────────┐ ┌───────▼──────┐
 │ MySQL       │ │ MySQL      │      │    MongoDB      │ │  Solr        │
 │ (usersdb)   │ │ (roomsdb)  │      │                 │ │              │
 └─────────────┘ └─────┬──────┘      └─────────────────┘ └───────▲──────┘
                       │                                         │
                       │        ┌───────────────┐                │
                       └───────►│   RabbitMQ    │────────────────┘
                    publica     │   (eventos)   │   consume + indexa
                   room.*       └───────────────┘

 Cachés: Memcached (distribuida) + caché local  ·  usada por users, rooms y search
```

**Flujo de eventos:** cuando `rooms-api` crea/actualiza/elimina una habitación, publica un
evento en RabbitMQ. `search-api` lo consume y actualiza su índice en Solr, manteniendo la
búsqueda sincronizada sin acoplar los servicios.

---

## 🚀 Microservicios

| Servicio | Puerto | Base de datos | Responsabilidad |
|----------|:------:|---------------|-----------------|
| **users-api** | `8080` | MySQL | Registro, login y gestión de usuarios (JWT) |
| **rooms-api** | `8081` | MySQL | CRUD de habitaciones; publica eventos |
| **reservations-api** | `8082` | MongoDB | Alta y consulta de reservas |
| **search-api** | `8083` | Solr | Búsqueda de habitaciones; consume eventos |

### Endpoints principales

**users-api** (`:8080`)
```
POST   /login              Iniciar sesión (devuelve JWT)
POST   /users              Registrar usuario
GET    /users              Listar usuarios
GET    /users/:id          Obtener usuario por ID
```

**rooms-api** (`:8081`)
```
GET    /health
GET    /api/v1/rooms                    Listar habitaciones
GET    /api/v1/rooms/available          Disponibles (vía search-api)
GET    /api/v1/rooms/number/:number     Por número
GET    /api/v1/rooms/:id                Por ID
# Rutas admin (requieren JWT con rol admin)
POST   /api/v1/admin/rooms
PUT    /api/v1/admin/rooms/:id
PATCH  /api/v1/admin/rooms/:id/status
DELETE /api/v1/admin/rooms/:id
```

**reservations-api** (`:8082`)
```
GET    /api/reservations                                  Listar todas
POST   /api/reservations                                  Crear reserva
GET    /api/reservations/:id                              Por ID
DELETE /api/reservations/:id                              Eliminar
GET    /api/reservations/users/:user_id/myreservations   Reservas de un usuario
```

**search-api** (`:8083`)
```
GET    /health
GET    /api/search/rooms   Búsqueda de habitaciones (con filtros)
```

---

## ⚙️ Puesta en marcha

### Requisitos
- [Docker](https://www.docker.com/) y Docker Compose
- [Node.js](https://nodejs.org/) 18+ (para el frontend)

### 1. Variables de entorno (opcional)
El proyecto funciona con valores por defecto. Para personalizarlos:
```bash
cp .env.example .env   # y editá los valores
```

### 2. Levantar el backend
```bash
docker-compose up --build
```
Esto inicia los 4 microservicios y toda su infraestructura (MySQL x2, MongoDB, Solr,
Memcached, RabbitMQ). El panel de RabbitMQ queda en http://localhost:15672.

### 3. Levantar el frontend
```bash
cd frontend
npm install
npm run dev
```
La aplicación queda disponible en **http://localhost:3000** (Vite hace de proxy hacia los
servicios del backend).

---

## 📁 Estructura del proyecto

```
.
├── docker-compose.yml          # Orquestación de todos los servicios
├── .env.example                # Plantilla de variables de entorno
├── init-db.sql                 # Inicialización de bases MySQL
├── frontend/                   # SPA en React + TypeScript + Vite
│   └── src/
│       ├── components/         # UI, layout y rutas protegidas
│       ├── pages/              # Vistas (Home, Login, Rooms, Reservas, Admin…)
│       ├── services/           # Clientes HTTP por dominio
│       ├── context/            # AuthContext (JWT)
│       └── lib/                # Configuración de axios
└── services/
    ├── users-api/              # Go · MySQL
    ├── rooms-api/              # Go · MySQL · RabbitMQ (publisher)
    ├── reservations-api/       # Go · MongoDB · RabbitMQ
    ├── search-api/             # Go · Solr · RabbitMQ (consumer) · Memcached
    └── solr/                   # Configset del core de Solr
```

Cada microservicio sigue una arquitectura por capas: `controllers` → `services` →
`repositories` → `domain`, con `config` para infraestructura y `utils` para helpers
(JWT, hashing, errores).

---

## 🔐 Notas de seguridad

Las credenciales de `docker-compose.yml` y `.env.example` son **valores de ejemplo para
desarrollo local**. En un entorno real, definí un `.env` propio con secretos fuertes y
nunca lo subas al repositorio (ya está en `.gitignore`).

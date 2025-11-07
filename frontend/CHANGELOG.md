# Changelog - Frontend Grand Prix Hotel

## [1.1.0] - 2025-11-07

### Changed
- **Rebranding completo**: Cambiado el nombre del hotel de "TUTE Hotel" a "Grand Prix"
  - Navbar: Logo y título actualizados
  - Home: Títulos y descripciones actualizados
  - Login: Título actualizado
  - Register: Título actualizado
  - README.md: Documentación actualizada
  - FRONTEND_QUICKSTART.md: Guía de inicio actualizada

### Verified
- **Compatibilidad con Backend**: Verificada conexión correcta con los 3 microservicios
  - users-api (puerto 8080): Login, registro, gestión de usuarios ✓
  - rooms-api (puerto 8081): CRUD de habitaciones con filtros ✓
  - reservations-api (puerto 8082): Gestión de reservas ✓

### Fixed
- Configuración de proxy en vite.config.ts para rutas correctas:
  - `/login` → users-api
  - `/users` → users-api
  - `/api/v1/rooms` → rooms-api
  - `/api/reservations` → reservations-api

## [1.0.0] - 2025-11-07

### Added
- Frontend completo con React + TypeScript + Vite
- Autenticación JWT con AuthContext
- Sistema de rutas protegidas
- Páginas: Home, Login, Register, Rooms, ReserveRoom, MyReservations
- Componentes UI reutilizables: Button, Input, Card
- Servicios de API para users, rooms, reservations
- Estilos con TailwindCSS
- Diseño responsive y moderno
- Documentación completa

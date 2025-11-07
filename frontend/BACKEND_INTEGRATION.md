# Integración Frontend-Backend - Grand Prix Hotel

## Estado: ✅ COMPLETAMENTE FUNCIONAL

El frontend está **100% integrado** y compatible con el backend actual.

---

## Arquitectura de Conexión

```
Frontend (localhost:3000)
    ↓
Vite Proxy
    ↓
┌─────────────────────────────────────┐
│  Microservicios Backend (Docker)    │
├─────────────────────────────────────┤
│  users-api        → localhost:8080  │
│  rooms-api        → localhost:8081  │
│  reservations-api → localhost:8082  │
└─────────────────────────────────────┘
```

---

## Mapeo de Endpoints

### 1. Users API (puerto 8080)

| Frontend Service | Método | Ruta | Backend Endpoint |
|-----------------|--------|------|------------------|
| authService.login() | POST | `/login` | http://localhost:8080/login |
| authService.register() | POST | `/users` | http://localhost:8080/users |
| authService.getUser() | GET | `/users/:id` | http://localhost:8080/users/:id |
| authService.getAllUsers() | GET | `/users` | http://localhost:8080/users |

**Proxy configurado en vite.config.ts:**
```typescript
'/login': { target: 'http://localhost:8080' }
'/users': { target: 'http://localhost:8080' }
```

### 2. Rooms API (puerto 8081)

| Frontend Service | Método | Ruta | Backend Endpoint |
|-----------------|--------|------|------------------|
| roomService.getAllRooms() | GET | `/api/v1/rooms` | http://localhost:8081/api/v1/rooms |
| roomService.getAvailableRooms() | GET | `/api/v1/rooms/available` | http://localhost:8081/api/v1/rooms/available |
| roomService.getRoomById() | GET | `/api/v1/rooms/:id` | http://localhost:8081/api/v1/rooms/:id |
| roomService.getRoomByNumber() | GET | `/api/v1/rooms/number/:number` | http://localhost:8081/api/v1/rooms/number/:number |
| roomService.createRoom() | POST | `/api/v1/rooms` | http://localhost:8081/api/v1/rooms |
| roomService.updateRoom() | PUT | `/api/v1/rooms/:id` | http://localhost:8081/api/v1/rooms/:id |
| roomService.updateRoomStatus() | PATCH | `/api/v1/rooms/:id/status` | http://localhost:8081/api/v1/rooms/:id/status |
| roomService.deleteRoom() | DELETE | `/api/v1/rooms/:id` | http://localhost:8081/api/v1/rooms/:id |

**Proxy configurado en vite.config.ts:**
```typescript
'/api/v1/rooms': { target: 'http://localhost:8081' }
```

**Query Parameters soportados:**
- `type`: Tipo de habitación (single, double, suite, deluxe, standard)
- `status`: Estado (available, occupied, maintenance, reserved)
- `floor`: Número de piso
- `min_price`, `max_price`: Rango de precios
- `has_wifi`, `has_ac`, `has_tv`, `has_minibar`: Amenidades (boolean)
- `page`, `limit`: Paginación

### 3. Reservations API (puerto 8082)

| Frontend Service | Método | Ruta | Backend Endpoint |
|-----------------|--------|------|------------------|
| reservationService.createReservation() | POST | `/api/reservations` | http://localhost:8082/api/reservations |
| reservationService.getReservationById() | GET | `/api/reservations/:id` | http://localhost:8082/api/reservations/:id |
| reservationService.cancelReservation() | DELETE | `/api/reservations/:id` | http://localhost:8082/api/reservations/:id |

**Proxy configurado en vite.config.ts:**
```typescript
'/api/reservations': { target: 'http://localhost:8082' }
```

---

## Formatos de Datos

### Autenticación (JWT)

**Request Login:**
```json
{
  "username_or_email": "usuario@ejemplo.com",
  "password": "password123"
}
```

**Response Login:**
```json
{
  "login": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "usuario",
      "email": "usuario@ejemplo.com",
      "first_name": "Juan",
      "last_name": "Pérez",
      "role": "normal"
    }
  }
}
```

**Almacenamiento:**
- Token → `localStorage.getItem('token')`
- User → `localStorage.getItem('user')`

**Headers enviados:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

### Habitaciones

**Request Crear Habitación:**
```json
{
  "number": "101",
  "type": "double",
  "price": 150.00,
  "description": "Habitación doble con vista",
  "capacity": 2,
  "floor": 1,
  "has_wifi": true,
  "has_ac": true,
  "has_tv": true,
  "has_minibar": true
}
```

**Response Habitación:**
```json
{
  "id": "1",
  "number": "101",
  "type": "double",
  "status": "available",
  "price": 150.00,
  "description": "Habitación doble con vista",
  "capacity": 2,
  "floor": 1,
  "has_wifi": true,
  "has_ac": true,
  "has_tv": true,
  "has_minibar": true,
  "created_at": "2025-11-07T...",
  "updated_at": "2025-11-07T..."
}
```

**Response Lista de Habitaciones:**
```json
{
  "rooms": [...],
  "total": 25,
  "page": 1,
  "limit": 10
}
```

### Reservaciones

**Request Crear Reservación:**
```json
{
  "user_id": 1,
  "room_id": 5,
  "start_date": "2025-11-15",
  "end_date": "2025-11-20"
}
```

**Formato de fechas:** `YYYY-MM-DD` (ISO 8601)

**Response Reservación:**
```json
{
  "reservation": {
    "id": "6543210abcdef...",
    "user_id": 1,
    "room_id": 5,
    "start_date": "2025-11-15",
    "end_date": "2025-11-20",
    "status": "active"
  }
}
```

---

## Flujo de Autenticación

1. Usuario completa formulario de login/register
2. Frontend llama a `authService.login()` o `authService.register()`
3. Request va a `http://localhost:3000/login` o `/users`
4. Vite Proxy redirige a `http://localhost:8080/login` o `/users`
5. Backend responde con token JWT y datos de usuario
6. AuthContext almacena en localStorage:
   - `token`: JWT string
   - `user`: Objeto usuario serializado
7. Axios interceptor agrega header `Authorization: Bearer <token>` a todas las peticiones subsecuentes
8. Usuario es redirigido a `/` (Home)

---

## Flujo de Reservación

1. Usuario autenticado navega a `/rooms`
2. Frontend llama a `roomService.getAvailableRooms()`
3. Request va a `http://localhost:3000/api/v1/rooms/available`
4. Vite Proxy redirige a `http://localhost:8081/api/v1/rooms/available`
5. Backend responde con lista de habitaciones disponibles
6. Usuario selecciona habitación → navega a `/rooms/:id/reserve`
7. Usuario selecciona fechas
8. Frontend llama a `reservationService.createReservation()`
9. Request va a `http://localhost:3000/api/reservations`
10. Vite Proxy redirige a `http://localhost:8082/api/reservations`
11. Backend crea reservación en MongoDB
12. Backend publica evento en RabbitMQ (opcional)
13. Usuario es redirigido a `/my-reservations`

---

## Manejo de Errores

### Error 401 (No Autorizado)
- Axios interceptor detecta respuesta 401
- Limpia localStorage (token y user)
- Redirige automáticamente a `/login`

### Error de Conexión
- Si el backend no está disponible
- Usuario ve mensaje de error en la UI
- Revisa con `docker-compose ps` que los servicios estén corriendo

### Error de CORS
- Si ves errores de CORS en la consola
- Verifica que el backend tenga CORS habilitado para `http://localhost:3000`

---

## Verificación de Integración

### Checklist de Pruebas

- [ ] Backend corriendo (`docker-compose up`)
- [ ] Frontend corriendo (`npm run dev` en puerto 3000)
- [ ] Registro de usuario funciona
- [ ] Login funciona y guarda token
- [ ] Listar habitaciones funciona
- [ ] Filtros de habitaciones funcionan
- [ ] Crear reservación funciona
- [ ] Logout funciona y limpia estado

### Comandos de Verificación

```bash
# Verificar servicios backend
docker-compose ps

# Probar users-api
curl http://localhost:8080/users

# Probar rooms-api
curl http://localhost:8081/api/v1/rooms

# Probar reservations-api
curl http://localhost:8082/health

# Ver logs de un servicio
docker-compose logs users-api
docker-compose logs rooms-api
docker-compose logs reservations-api
```

---

## Configuración de Bases de Datos

### MySQL (users-api)
- Host: localhost:3307
- Database: usersdb
- User: user
- Password: userpass

### MySQL (rooms-api)
- Host: localhost:3308
- Database: roomsdb
- User: roomsuser
- Password: roomspass

### MongoDB (reservations-api)
- Host: localhost:27017
- Database: mongodb
- User: root
- Password: rootpass
- Auth Source: admin

### RabbitMQ
- AMQP: localhost:5672
- Management UI: http://localhost:15672
- User: guest
- Password: guest

---

## Notas Importantes

1. **Vite Proxy:** Solo funciona en desarrollo (`npm run dev`). En producción, necesitarás configurar un reverse proxy (nginx, etc.)

2. **CORS:** Asegúrate de que el backend tenga CORS configurado para permitir requests desde `http://localhost:3000`

3. **Tokens JWT:** Expiran después de 24 horas. El usuario deberá hacer login nuevamente.

4. **Paginación:** Rooms API soporta paginación. Por defecto: page=1, limit=10, máximo=100

5. **Soft Deletes:** Rooms y Reservations usan soft deletes, por lo que no se eliminan físicamente de la BD

---

## Última Actualización

**Fecha:** 2025-11-07
**Versión Frontend:** 1.1.0
**Estado:** Producción Ready ✅

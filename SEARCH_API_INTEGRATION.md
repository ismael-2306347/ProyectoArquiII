# Search API - Integración Completa

## Resumen de la Implementación

Se ha creado un nuevo microservicio **search-api** completamente funcional e integrado con tu sistema de reservas de hotel. El servicio implementa búsqueda avanzada con Apache Solr, doble caché (CCache + Memcached), y sincronización en tiempo real vía RabbitMQ.

## Archivos Creados/Modificados

### Nuevos Archivos en `services/search-api/`

```
services/search-api/
├── cmd/server/main.go                      ✅ Servidor principal
├── controllers/search-controller.go        ✅ Handlers HTTP
├── services/search-service.go              ✅ Lógica de negocio con doble caché
├── repositories/
│   ├── solr-repository.go                  ✅ Repositorio Solr (indexar/buscar)
│   ├── cache-local-repo.go                 ✅ Repositorio CCache
│   └── cache-distributed-repo.go           ✅ Repositorio Memcached
├── consumers/rooms-consumer.go             ✅ Consumer RabbitMQ
├── domain/
│   ├── room_search.go                      ✅ Modelos y DTOs de Solr
│   └── dto.go                              ✅ Request/Response DTOs
├── config/
│   ├── solr.go                             ✅ Configuración Solr
│   ├── rabbitmq.go                         ✅ Configuración RabbitMQ
│   ├── cache.go                            ✅ Configuración cachés
│   └── http_clients.go                     ✅ Cliente HTTP para rooms-api
├── utils/
│   ├── jwt.go                              ✅ Validación JWT
│   └── errors.go                           ✅ Errores personalizados
├── scripts/
│   ├── schema.json                         ✅ Schema de Solr
│   └── README.md                           ✅ Documentación de Solr
├── Dockerfile                               ✅ Dockerfile multi-stage
├── .dockerignore                            ✅ Exclusiones Docker
├── go.mod                                   ✅ Dependencias Go
├── go.sum                                   ✅ Checksums
├── README.md                                ✅ Documentación completa
└── QUICKSTART.md                            ✅ Guía paso a paso
```

### Archivos Modificados

1. **`services/rooms-api/events/publisher.go`** ✅
   - Cambio de exchange de `room_events` a `rooms`
   - Simplificación de eventos para formato compatible con search-api
   - Adición de campo `timestamp`
   - Routing keys: `room.created`, `room.updated`, `room.deleted`

2. **`docker-compose.yml`** ✅
   - Agregado servicio `solr` con core `rooms-core`
   - Agregado servicio `search-api` con todas las dependencias
   - Configuración de variables de entorno

## Características Implementadas

### 1. Motor de Búsqueda con Solr

- **Core**: `rooms-core`
- **Campos indexados**:
  - id, number, type, status
  - price, capacity, floor
  - has_wifi, has_ac, has_tv, has_minibar
  - description (texto completo)
  - created_at, updated_at

### 2. Doble Caché (L1 + L2)

**Estrategia**:
```
Request → CCache (L1, 60s TTL)
          ├─ HIT → Return
          └─ MISS → Memcached (L2, 300s TTL)
                    ├─ HIT → Save L1 → Return
                    └─ MISS → Solr → Save L2 → Save L1 → Return
```

**Ventajas**:
- L1 (CCache): Ultra-rápido, in-memory local
- L2 (Memcached): Compartido entre instancias, persistente

### 3. Sincronización en Tiempo Real

**Flujo**:
1. rooms-api publica evento a RabbitMQ (exchange: `rooms`)
2. search-api consume evento
3. Consumer obtiene datos completos desde rooms-api
4. Indexa/actualiza/elimina en Solr
5. Invalida caché

**Eventos soportados**:
- `room.created` → Indexar en Solr
- `room.updated` → Actualizar en Solr
- `room.deleted` → Eliminar de Solr

### 4. API RESTful

**Endpoint principal**: `GET /api/search/rooms`

**Filtros soportados**:
- `q`: Búsqueda de texto libre
- `type`: Tipo de habitación
- `status`: Estado
- `floor`: Piso
- `min_price` / `max_price`: Rango de precio
- `has_wifi`, `has_ac`, `has_tv`, `has_minibar`: Amenities
- `sort`: Ordenamiento (ej: `price`, `-price`)
- `page` / `limit`: Paginación

**Ejemplo**:
```bash
GET /api/search/rooms?status=available&has_wifi=true&min_price=100&max_price=200&sort=-price&page=1&limit=10
```

### 5. Autenticación JWT (Opcional)

- Middleware `AuthMiddleware()` para proteger endpoints
- Middleware `OptionalAuthMiddleware()` para auth opcional
- Mismo secret que users-api para compatibilidad

## Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                     SEARCH API                               │
│                                                               │
│  ┌──────────────┐      ┌───────────────┐                    │
│  │ Controller   │ ───▶ │   Service     │                    │
│  │ (HTTP Gin)   │      │ (Business     │                    │
│  └──────────────┘      │  Logic)       │                    │
│                        └───────┬───────┘                     │
│                                │                              │
│         ┌──────────────────────┼──────────────────┐         │
│         │                      │                  │         │
│    ┌────▼────┐          ┌─────▼──────┐     ┌────▼─────┐   │
│    │ CCache  │          │ Memcached  │     │   Solr   │   │
│    │ (Local) │          │(Distributed)│     │ (Search) │   │
│    │  60s    │          │   300s     │     │  Engine  │   │
│    └─────────┘          └────────────┘     └──────────┘   │
│                                                              │
│  ┌───────────────────────────────────────────────────────┐ │
│  │         RabbitMQ Consumer                              │ │
│  │  (Escucha: room.created, room.updated, room.deleted)  │ │
│  └───────────────────────────┬───────────────────────────┘ │
└──────────────────────────────┼─────────────────────────────┘
                               │
                    ┌──────────▼──────────┐
                    │     RabbitMQ        │
                    │   Exchange: rooms   │
                    └──────────▲──────────┘
                               │
                    ┌──────────┴──────────┐
                    │    Rooms API        │
                    │  (Event Publisher)  │
                    └─────────────────────┘
```

## Dependencias Go

```go
// go.mod
require (
    github.com/bradfitz/gomemcache        // Memcached client
    github.com/gin-gonic/gin              // HTTP framework
    github.com/golang-jwt/jwt/v5          // JWT validation
    github.com/karlseguin/ccache/v3       // Local cache (CCache)
    github.com/rabbitmq/amqp091-go        // RabbitMQ client
)
```

## Variables de Entorno Requeridas

```bash
# Solr
SOLR_URL=http://solr:8983/solr
SOLR_CORE=rooms-core

# Cachés
MEMCACHED_HOST=memcached
MEMCACHED_PORT=11211
LOCAL_CACHE_TTL_SECONDS=60
DISTRIBUTED_CACHE_TTL_SECONDS=300

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

# Integración con rooms-api
ROOMS_API_BASE_URL=http://rooms-api:8080

# Autenticación (mismo secret que users-api)
JWT_SECRET=supersecreto_cámbialo

# Server
PORT=8083
GIN_MODE=release
```

## Pasos de Inicialización

### Opción 1: Con Docker Compose (Recomendado)

```bash
# 1. Levantar Solr y crear core
docker-compose up -d solr

# 2. Esperar 15 segundos y aplicar schema
curl -X POST -H 'Content-type:application/json' \
  http://localhost:8983/solr/rooms-core/schema \
  --data-binary @services/search-api/scripts/schema.json

# 3. Levantar resto de servicios
docker-compose up -d memcached rabbitmq mysql-rooms rooms-api search-api

# 4. Verificar
curl http://localhost:8083/health
```

### Opción 2: Desarrollo Local

```bash
# 1. Levantar dependencias
docker-compose up -d solr memcached rabbitmq mysql-rooms rooms-api

# 2. Aplicar schema de Solr
curl -X POST -H 'Content-type:application/json' \
  http://localhost:8983/solr/rooms-core/schema \
  --data-binary @services/search-api/scripts/schema.json

# 3. Configurar variables de entorno locales
export SOLR_URL=http://localhost:8983/solr
export SOLR_CORE=rooms-core
export MEMCACHED_HOST=localhost
export MEMCACHED_PORT=11211
export RABBITMQ_URL=amqp://guest:guest@localhost:5672/
export ROOMS_API_BASE_URL=http://localhost:8081
export JWT_SECRET=supersecreto_cámbialo
export PORT=8083

# 4. Ejecutar search-api
cd search-api
go mod download
go run cmd/server/main.go
```

## Testing Rápido

### 1. Crear habitación en rooms-api

```bash
curl -X POST http://localhost:8081/api/v1/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "number": "101",
    "type": "double",
    "price": 120.50,
    "description": "Habitación doble con vista al mar",
    "capacity": 2,
    "floor": 1,
    "has_wifi": true,
    "has_ac": true,
    "has_tv": true,
    "has_minibar": false
  }'
```

### 2. Verificar indexación en Solr

```bash
# Directamente en Solr
curl "http://localhost:8983/solr/rooms-core/select?q=*:*&wt=json"

# Ver logs del consumer
docker logs search-api | grep "indexed successfully"
```

### 3. Buscar con search-api

```bash
# Todas las habitaciones
curl "http://localhost:8083/api/search/rooms"

# Con filtros
curl "http://localhost:8083/api/search/rooms?status=available&has_wifi=true"

# Rango de precio
curl "http://localhost:8083/api/search/rooms?min_price=100&max_price=200&sort=price"

# Texto libre
curl "http://localhost:8083/api/search/rooms?q=mar"
```

### 4. Probar sincronización

```bash
# Actualizar habitación
curl -X PUT http://localhost:8081/api/v1/rooms/1 \
  -H "Content-Type: application/json" \
  -d '{"price": 135.00}'

# Ver logs de search-api
docker logs search-api | tail -20

# Buscar de nuevo para ver cambios
curl "http://localhost:8083/api/search/rooms?q=101"
```

## Monitoreo y Debugging

### Logs

```bash
# Todos los servicios
docker-compose logs -f

# Solo search-api
docker logs -f search-api

# Filtrar eventos RabbitMQ
docker logs search-api | grep "Received message"

# Filtrar indexaciones
docker logs search-api | grep "indexed successfully"

# Ver cache hits/misses
docker logs search-api | grep "Cache"
```

### RabbitMQ Management UI

- URL: http://localhost:15672
- Usuario: `guest`
- Contraseña: `guest`

Verificar:
- Exchange `rooms` existe
- Queue `search-api-rooms-queue` existe
- Bindings a routing keys: `room.created`, `room.updated`, `room.deleted`

### Solr Admin UI

- URL: http://localhost:8983/solr/#/rooms-core

Verificar:
- Core `rooms-core` está activo
- Schema tiene todos los campos definidos
- Query puede ver documentos indexados

### Health Checks

```bash
# Search API
curl http://localhost:8083/health

# Solr
curl http://localhost:8983/solr/rooms-core/admin/ping

# RabbitMQ
curl http://localhost:15672/api/healthchecks/node
```

## Troubleshooting Común

### Consumer no recibe eventos

**Problema**: Habitaciones creadas en rooms-api no se indexan en Solr

**Solución**:
```bash
# 1. Verificar que rooms-api publique eventos
docker logs rooms-api | grep "Evento publicado"

# 2. Verificar exchange en RabbitMQ
curl -u guest:guest http://localhost:15672/api/exchanges/%2F/rooms

# 3. Verificar queue existe
curl -u guest:guest http://localhost:15672/api/queues/%2F/search-api-rooms-queue

# 4. Reiniciar search-api
docker-compose restart search-api
```

### Solr retorna errores

**Problema**: Búsquedas fallan con error 400 o 500

**Solución**:
```bash
# Ver logs de Solr
docker logs solr

# Re-aplicar schema
curl -X POST -H 'Content-type:application/json' \
  http://localhost:8983/solr/rooms-core/schema \
  --data-binary @services/search-api/scripts/schema.json

# Si persiste, re-crear core
curl "http://localhost:8983/solr/admin/cores?action=UNLOAD&core=rooms-core&deleteIndex=true"
docker-compose restart solr
# Esperar 15 segundos
curl -X POST -H 'Content-type:application/json' \
  http://localhost:8983/solr/rooms-core/schema \
  --data-binary @services/search-api/scripts/schema.json
```

### Cache no funciona

**Problema**: Todas las búsquedas muestran "Cache MISS"

**Solución**:
```bash
# Verificar Memcached
docker ps | grep memcached
docker logs memcached

# Test manual
telnet localhost 11211
> stats
> quit

# Reiniciar Memcached
docker-compose restart memcached
```

## Próximos Pasos

1. **Integración Frontend**: Agregar componente de búsqueda en React
2. **Filtros Avanzados**: Búsqueda por disponibilidad de fechas
3. **Autocompletado**: Sugerencias en tiempo real
4. **Facets/Aggregations**: Contadores por tipo, piso, etc.
5. **Métricas**: Prometheus + Grafana
6. **Rate Limiting**: Protección contra abuso
7. **Re-indexación**: Script para re-indexar todas las habitaciones

## Documentación Adicional

- [README.md](services/search-api/README.md): Documentación completa de search-api
- [QUICKSTART.md](services/search-api/QUICKSTART.md): Guía paso a paso de inicio
- [scripts/README.md](services/search-api/scripts/README.md): Documentación del schema de Solr

## Contacto y Soporte

Para reportar bugs o solicitar features, crea un issue en el repositorio.

---

**Creado con**: Go 1.21, Gin, Apache Solr 9.5, RabbitMQ, Memcached, CCache
**Autor**: Claude Code Assistant
**Fecha**: Noviembre 2024

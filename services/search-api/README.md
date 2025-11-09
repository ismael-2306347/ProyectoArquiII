# Search API - Microservicio de Búsqueda

API de búsqueda avanzada para el sistema de reservas de habitaciones utilizando Apache Solr.

## 🚀 Características

- **Búsqueda full-text** con Apache Solr
- **Caché inteligente** con Memcached
- **Sincronización automática** vía RabbitMQ
- **Filtros avanzados** (precio, capacidad, amenidades)
- **Autocompletado** en tiempo real
- **Facetas dinámicas** para filtros

## 📋 Requisitos

- Go 1.21+
- Apache Solr 9.5
- Memcached
- RabbitMQ
- Rooms API (para reindexación)

## 🛠️ Instalación

### 1. Preparar estructura de directorios

```bash
mkdir -p services/search-api/{controllers,services,repositories,domain,events,config}
mkdir -p solr/configsets/rooms_core/conf
```

### 2. Copiar archivos de configuración de Solr

Coloca el archivo `managed-schema.xml` en:
```
solr/configsets/rooms_core/conf/managed-schema.xml
```

### 3. Inicializar módulo Go

```bash
cd services/search-api
go mod init search-api
go mod tidy
```

### 4. Compilar

```bash
go build -o search-api
```

## 🐳 Docker Compose

El servicio está configurado en `docker-compose.yml`:

```yaml
search-api:
  build:
    context: ./services/search-api
    dockerfile: dockerfile
  container_name: search-api
  restart: always
  depends_on:
    solr:
      condition: service_healthy
    memcached:
      condition: service_started
    rabbitmq:
      condition: service_healthy
  environment:
    - SOLR_URL=http://solr:8983/solr/rooms_core
    - MEMCACHED_HOST=memcached
    - MEMCACHED_PORT=11211
    - MEMCACHED_TTL=300
    - RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
    - ROOMS_API_URL=http://rooms-api:8080
  ports:
    - "8083:8080"
```

## 🔌 Endpoints

### Búsqueda

#### `GET /api/v1/search/rooms`

Busca habitaciones con filtros avanzados.

**Query Parameters:**
- `q`: Búsqueda por texto libre
- `min_price`: Precio mínimo por noche
- `max_price`: Precio máximo por noche
- `min_capacity`: Capacidad mínima
- `room_type`: Tipo de habitación
- `is_available`: Solo disponibles (true/false)
- `floor`: Número de piso
- `has_wifi`: Tiene WiFi (true/false)
- `has_ac`: Tiene aire acondicionado (true/false)
- `has_tv`: Tiene TV (true/false)
- `has_minibar`: Tiene minibar (true/false)
- `page`: Número de página (default: 1)
- `limit`: Resultados por página (default: 10, max: 100)
- `sort`: Ordenamiento (price_asc, price_desc, capacity_asc, capacity_desc)

**Ejemplo:**
```bash
curl "http://localhost:8083/api/v1/search/rooms?q=suite&min_price=100&max_price=300&is_available=true&page=1&limit=10"
```

**Respuesta:**
```json
{
  "total_results": 25,
  "results": [
    {
      "id": "1",
      "room_number": "101",
      "room_type": "suite",
      "capacity": 2,
      "price_per_night": 250.00,
      "status": "available",
      "description": "Suite de lujo con vista al mar",
      "amenities": ["WiFi", "Aire Acondicionado", "TV", "Minibar"],
      "floor": 1,
      "is_available": true
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_pages": 3
}
```

#### `GET /api/v1/search/rooms/suggestions`

Obtiene sugerencias de autocompletado.

**Query Parameters:**
- `q`: Prefijo de búsqueda (requerido)
- `limit`: Máximo de sugerencias (default: 10, max: 20)

**Ejemplo:**
```bash
curl "http://localhost:8083/api/v1/search/rooms/suggestions?q=sui&limit=5"
```

**Respuesta:**
```json
{
  "suggestions": ["suite", "suite deluxe", "suite presidential"]
}
```

#### `GET /api/v1/search/rooms/facets`

Obtiene facetas para filtros dinámicos.

**Ejemplo:**
```bash
curl "http://localhost:8083/api/v1/search/rooms/facets"
```

**Respuesta:**
```json
{
  "room_types": {
    "single": 20,
    "double": 35,
    "suite": 15
  },
  "status_counts": {
    "available": 50,
    "occupied": 15,
    "maintenance": 5
  },
  "floor_counts": {
    "1": 20,
    "2": 25,
    "3": 25
  }
}
```

### Administración

#### `POST /api/v1/admin/index/rooms/full`

Reindexar todas las habitaciones desde rooms-api.

**Ejemplo:**
```bash
curl -X POST "http://localhost:8083/api/v1/admin/index/rooms/full"
```

#### `POST /api/v1/admin/index/room/:id`

Indexar una habitación manualmente.

**Body:**
```json
{
  "room_number": "101",
  "room_type": "suite",
  "capacity": 2,
  "price_per_night": 250.00,
  "status": "available",
  "description": "Suite de lujo",
  "amenities": ["WiFi", "TV"],
  "floor": 1,
  "is_available": true
}
```

#### `DELETE /api/v1/admin/index/room/:id`

Eliminar una habitación del índice.

#### `GET /api/v1/admin/index/stats`

Obtiene estadísticas del índice.

## 🔄 Eventos RabbitMQ

El servicio escucha los siguientes eventos:

- `room.created`: Nueva habitación creada
- `room.updated`: Habitación actualizada
- `room.deleted`: Habitación eliminada
- `room.status.changed`: Estado de habitación cambiado
- `reservation.created`: Nueva reserva (marca habitación como no disponible)
- `reservation.cancelled`: Reserva cancelada (marca habitación como disponible)

## 🧪 Pruebas

### Iniciar servicios
```bash
docker-compose up -d
```

### Verificar salud
```bash
curl http://localhost:8083/health
```

### Reindexar habitaciones
```bash
curl -X POST http://localhost:8083/api/v1/admin/index/rooms/full
```

### Buscar habitaciones disponibles
```bash
curl "http://localhost:8083/api/v1/search/rooms?is_available=true&limit=5"
```

### Buscar por precio
```bash
curl "http://localhost:8083/api/v1/search/rooms?min_price=100&max_price=200"
```

## 📊 Monitoreo

- **Solr Admin UI**: http://localhost:8983/solr
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)
- **API Health**: http://localhost:8083/health

## 🔧 Variables de Entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `SOLR_URL` | URL de Solr | `http://localhost:8983/solr/rooms_core` |
| `MEMCACHED_HOST` | Host de Memcached | `localhost` |
| `MEMCACHED_PORT` | Puerto de Memcached | `11211` |
| `MEMCACHED_TTL` | TTL del cache (segundos) | `300` |
| `RABBITMQ_URL` | URL de RabbitMQ | `amqp://guest:guest@localhost:5672/` |
| `ROOMS_API_URL` | URL de Rooms API | `http://localhost:8081` |
| `PORT` | Puerto del servidor | `8080` |

## 🏗️ Arquitectura

```
Frontend
    ↓
Search API (8083)
    ↓
├── Solr (8983) ← Motor de búsqueda
├── Memcached (11211) ← Caché
└── RabbitMQ (5672) ← Eventos
         ↑
    Rooms API (8081)
    Reservations API (8082)
```

## 📝 Notas

1. **Primera vez**: Ejecutar reindexación completa después de iniciar
2. **Caché**: Se invalida automáticamente con eventos
3. **Reintentos**: RabbitMQ reintenta conexión hasta 15 veces
4. **Límites**: Máximo 100 resultados por página

## 🐛 Troubleshooting

### Solr no responde
```bash
docker-compose logs solr
docker-compose restart solr
```

### No se indexan habitaciones
```bash
docker-compose logs search-api
# Verificar que RabbitMQ esté funcionando
docker-compose logs rabbitmq
```

### Cache no funciona
```bash
docker-compose logs memcached
# Reiniciar memcached
docker-compose restart memcached
```
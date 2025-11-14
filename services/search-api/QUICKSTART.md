# Search API - Guía de Inicio Rápido

Esta guía te llevará paso a paso para levantar y probar la search-api integrada con tu sistema de reservas.

## Prerrequisitos

- Docker y Docker Compose instalados
- Puerto 8083 disponible para search-api
- Puerto 8983 disponible para Solr

## Paso 1: Iniciar Solr

Primero, iniciemos Solr para crear el core automáticamente:

```bash
# Desde el directorio raíz del proyecto
docker-compose up -d solr
```

Espera unos 10-15 segundos y verifica que Solr esté corriendo:

```bash
curl http://localhost:8983/solr/rooms-core/admin/ping
```

Deberías ver una respuesta JSON con `"status":"OK"`.

## Paso 2: Configurar el Schema de Solr

Aplica el schema para definir los campos indexados:

```bash
curl -X POST -H 'Content-type:application/json' \
  http://localhost:8983/solr/rooms-core/schema \
  --data-binary @search-api/scripts/schema.json
```

Verifica que se aplicó correctamente:

```bash
curl http://localhost:8983/solr/rooms-core/schema/fields | grep "has_wifi"
```

## Paso 3: Iniciar Memcached y RabbitMQ

```bash
docker-compose up -d memcached rabbitmq
```

Espera a que RabbitMQ esté listo (unos 20 segundos):

```bash
# Ver logs hasta que veas "Server startup complete"
docker logs rabbitmq
```

Puedes acceder a RabbitMQ Management en http://localhost:15672 (usuario: guest, contraseña: guest)

## Paso 4: Iniciar Rooms API

```bash
docker-compose up -d mysql-rooms rooms-api
```

Espera a que rooms-api esté listo:

```bash
# Ver logs
docker logs -f rooms-api
```

Cuando veas "Server started on port 8080", presiona Ctrl+C.

## Paso 5: Iniciar Search API

```bash
docker-compose up -d search-api
```

Ver los logs para confirmar que todo se conectó correctamente:

```bash
docker logs -f search-api
```

Deberías ver:
```
Successfully connected to RabbitMQ
Queue search-api-rooms-queue bound to exchange rooms with routing key room.created
Queue search-api-rooms-queue bound to exchange rooms with routing key room.updated
Queue search-api-rooms-queue bound to exchange rooms with routing key room.deleted
RabbitMQ consumer started, waiting for messages...
Starting HTTP server on port 8083
```

## Paso 6: Verificar Health Checks

```bash
# Health check de search-api
curl http://localhost:8083/health

# Debería responder:
# {"status":"healthy","service":"search-api"}
```

## Paso 7: Crear Habitaciones de Prueba

Vamos a crear algunas habitaciones en rooms-api. El consumer de search-api automáticamente las indexará en Solr.

```bash
# Habitación 1: Double con WiFi y AC
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

# Habitación 2: Suite con todas las amenities
curl -X POST http://localhost:8081/api/v1/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "number": "201",
    "type": "suite",
    "price": 250.00,
    "description": "Suite de lujo con jacuzzi",
    "capacity": 4,
    "floor": 2,
    "has_wifi": true,
    "has_ac": true,
    "has_tv": true,
    "has_minibar": true
  }'

# Habitación 3: Single económica
curl -X POST http://localhost:8081/api/v1/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "number": "305",
    "type": "single",
    "price": 80.00,
    "description": "Habitación individual básica",
    "capacity": 1,
    "floor": 3,
    "has_wifi": true,
    "has_ac": false,
    "has_tv": false,
    "has_minibar": false
  }'
```

## Paso 8: Verificar Indexación en Solr

Verifica que las habitaciones se indexaron:

```bash
curl "http://localhost:8983/solr/rooms-core/select?q=*:*&wt=json&rows=10" | jq
```

Deberías ver 3 documentos en Solr.

También puedes ver los logs de search-api:

```bash
docker logs search-api | grep "indexed successfully"
```

## Paso 9: Probar Búsquedas

### Búsqueda simple - Todas las habitaciones

```bash
curl "http://localhost:8083/api/search/rooms" | jq
```

### Búsqueda con filtros - Solo disponibles con WiFi

```bash
curl "http://localhost:8083/api/search/rooms?status=available&has_wifi=true" | jq
```

### Búsqueda por rango de precio

```bash
curl "http://localhost:8083/api/search/rooms?min_price=100&max_price=200" | jq
```

### Búsqueda con ordenamiento por precio descendente

```bash
curl "http://localhost:8083/api/search/rooms?sort=-price" | jq
```

### Búsqueda de texto libre

```bash
curl "http://localhost:8083/api/search/rooms?q=suite" | jq
```

### Búsqueda compleja - Habitaciones en piso 1 o 2, con AC, precio < 150

```bash
curl "http://localhost:8083/api/search/rooms?max_price=150&has_ac=true" | jq
```

### Paginación

```bash
# Primera página (2 items)
curl "http://localhost:8083/api/search/rooms?page=1&limit=2" | jq

# Segunda página
curl "http://localhost:8083/api/search/rooms?page=2&limit=2" | jq
```

## Paso 10: Probar Cache

Ejecuta la misma búsqueda dos veces y observa los logs:

```bash
# Primera vez (Cache MISS)
curl "http://localhost:8083/api/search/rooms?status=available" | jq

# Segunda vez (Cache HIT - local)
curl "http://localhost:8083/api/search/rooms?status=available" | jq
```

Ver logs:

```bash
docker logs search-api | tail -20
```

Deberías ver:
```
Cache MISS: search:xxxxx - Querying Solr
Cache HIT (local): search:xxxxx
```

## Paso 11: Probar Sincronización en Tiempo Real

### Actualizar una habitación

```bash
curl -X PUT http://localhost:8081/api/v1/rooms/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 135.00,
    "description": "Habitación doble renovada con vista al mar"
  }'
```

Ver logs de search-api:

```bash
docker logs -f search-api
```

Deberías ver:
```
Received message: {"event_type":"updated","room_id":1,"timestamp":"..."}
Room 1 indexed successfully in Solr
Event processed successfully: updated for room 1
```

### Eliminar una habitación

```bash
curl -X DELETE http://localhost:8081/api/v1/rooms/3
```

Ver logs:

```bash
docker logs search-api | grep "deleted"
```

Deberías ver:
```
Room 3 deleted successfully from Solr
```

Verifica en Solr:

```bash
curl "http://localhost:8983/solr/rooms-core/select?q=id:3&wt=json" | jq
```

No debería haber resultados.

## Paso 12: Monitoreo con RabbitMQ Management

1. Abre http://localhost:15672
2. Usuario: `guest`, Contraseña: `guest`
3. Ve a la pestaña "Exchanges"
4. Busca el exchange `rooms`
5. Ve a la pestaña "Queues"
6. Busca la queue `search-api-rooms-queue`
7. Deberías ver los bindings a las routing keys: `room.created`, `room.updated`, `room.deleted`

## Troubleshooting

### "Connection refused" al hacer búsquedas

```bash
# Verificar que search-api esté corriendo
docker ps | grep search-api

# Ver logs de errores
docker logs search-api
```

### Las habitaciones no se indexan en Solr

```bash
# 1. Verificar que rooms-api esté publicando eventos
docker logs rooms-api | grep "Evento publicado"

# 2. Verificar que RabbitMQ tenga mensajes
# Ir a http://localhost:15672 y ver la queue

# 3. Verificar logs del consumer
docker logs search-api | grep "Received message"

# 4. Verificar que el exchange y queue existan
curl -u guest:guest http://localhost:15672/api/exchanges/%2F/rooms
curl -u guest:guest http://localhost:15672/api/queues/%2F/search-api-rooms-queue
```

### Solr retorna errores

```bash
# Ver logs de Solr
docker logs solr

# Re-crear el core si es necesario
curl "http://localhost:8983/solr/admin/cores?action=UNLOAD&core=rooms-core&deleteIndex=true"
curl "http://localhost:8983/solr/admin/cores?action=CREATE&name=rooms-core&configSet=_default"

# Re-aplicar schema
curl -X POST -H 'Content-type:application/json' \
  http://localhost:8983/solr/rooms-core/schema \
  --data-binary @search-api/scripts/schema.json
```

### Cache no funciona

```bash
# Verificar Memcached
docker ps | grep memcached
docker logs memcached

# Test manual de Memcached
telnet localhost 11211
> stats
> quit
```

## Detener Servicios

```bash
# Detener solo search-api
docker-compose stop search-api

# Detener todos los servicios
docker-compose down

# Detener y limpiar volúmenes (¡CUIDADO! Borra todos los datos)
docker-compose down -v
```

## Próximos Pasos

1. Integrar con el frontend para mostrar búsquedas en la UI
2. Agregar más filtros (por disponibilidad de fechas, etc.)
3. Implementar autocompletado
4. Agregar métricas y monitoreo
5. Configurar rate limiting en producción

## Comandos Útiles de Referencia

```bash
# Ver todos los logs en tiempo real
docker-compose logs -f

# Ver logs de un servicio específico
docker logs -f search-api

# Reiniciar un servicio
docker-compose restart search-api

# Ver estado de todos los contenedores
docker-compose ps

# Ejecutar comando dentro del contenedor
docker exec -it search-api sh

# Ver uso de recursos
docker stats

# Limpiar índice de Solr
curl http://localhost:8983/solr/rooms-core/update?commit=true \
  -H "Content-Type: application/json" \
  -d '{"delete":{"query":"*:*"}}'
```

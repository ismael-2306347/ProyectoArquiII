# Rooms API

A microservice for managing hotel rooms using MongoDB and Go with Gin framework.

## Features

- **CRUD Operations**: Create, read, update, and delete rooms
- **Room Management**: Track room status, type, pricing, and amenities
- **Filtering & Search**: Filter rooms by type, status, floor, price range, and amenities
- **Pagination**: Support for paginated room listings
- **MongoDB Integration**: Uses MongoDB for data persistence
- **RESTful API**: Clean REST API design with proper HTTP status codes

## Room Model

### Room Types
- `single` - Single occupancy room
- `double` - Double occupancy room  
- `suite` - Suite room
- `deluxe` - Deluxe room
- `standard` - Standard room

### Room Status
- `available` - Room is available for booking
- `occupied` - Room is currently occupied
- `maintenance` - Room is under maintenance
- `reserved` - Room is reserved but not yet occupied

### Room Properties
- **Basic Info**: Number, type, status, price, description, capacity, floor
- **Amenities**: WiFi, AC, TV, minibar availability
- **Timestamps**: Created and updated timestamps

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Room Management
- `POST /api/v1/rooms` - Create a new room
- `GET /api/v1/rooms` - Get all rooms (with filtering and pagination)
- `GET /api/v1/rooms/available` - Get available rooms only
- `GET /api/v1/rooms/:id` - Get room by ID
- `GET /api/v1/rooms/number/:number` - Get room by number
- `PUT /api/v1/rooms/:id` - Update room
- `PATCH /api/v1/rooms/:id/status` - Update room status
- `DELETE /api/v1/rooms/:id` - Delete room

### Query Parameters for GET /api/v1/rooms
- `type` - Filter by room type
- `status` - Filter by room status
- `floor` - Filter by floor number
- `min_price` - Minimum price filter
- `max_price` - Maximum price filter
- `has_wifi` - Filter by WiFi availability
- `has_ac` - Filter by AC availability
- `has_tv` - Filter by TV availability
- `has_minibar` - Filter by minibar availability
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)

## Environment Variables

- `MONGODB_URI` - MongoDB connection string (default: mongodb://localhost:27017)
- `MONGODB_DB` - Database name (default: roomsdb)
- `PORT` - Server port (default: 8080)

## Example Usage

### Create a Room
```bash
curl -X POST http://localhost:8081/api/v1/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "number": "101",
    "type": "double",
    "price": 150.00,
    "description": "Comfortable double room with city view",
    "capacity": 2,
    "floor": 1,
    "has_wifi": true,
    "has_ac": true,
    "has_tv": true,
    "has_minibar": false
  }'
```

### Get Available Rooms
```bash
curl "http://localhost:8081/api/v1/rooms/available?type=double&floor=1&min_price=100&max_price=200"
```

### Update Room Status
```bash
curl -X PATCH http://localhost:8081/api/v1/rooms/ROOM_ID/status \
  -H "Content-Type: application/json" \
  -d '{"status": "occupied"}'
```

## Running the Service

### Using Docker Compose
```bash
docker-compose up rooms-api
```

### Running Locally
```bash
# Install dependencies
go mod tidy

# Run the service
go run ./cmd/server
```

## Database Schema

The service uses MongoDB with the following collection structure:

```json
{
  "_id": "ObjectId",
  "number": "string",
  "type": "string",
  "status": "string", 
  "price": "number",
  "description": "string",
  "capacity": "number",
  "floor": "number",
  "has_wifi": "boolean",
  "has_ac": "boolean", 
  "has_tv": "boolean",
  "has_minibar": "boolean",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- `400 Bad Request` - Invalid request data
- `404 Not Found` - Room not found
- `409 Conflict` - Room number already exists
- `500 Internal Server Error` - Database or server error

## Development

### Project Structure
```
rooms-api/
├── cmd/server/          # Main application entry point
├── config/              # Database configuration
├── controllers/         # HTTP handlers
├── domain/              # Data models and DTOs
├── repositories/        # Data access layer
├── services/            # Business logic layer
├── utils/               # Utility functions
├── dockerfile           # Docker configuration
└── go.mod              # Go module dependencies
```

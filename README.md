
## 🏗️ Arquitectura

```
                        ┌─────────────┐
                        │   Frontend  │ (React - Puerto 3000)
                        │   (React)   │
                        └──────┬──────┘
                               │
       ┌───────────────────────┼───────────────────────┬─────────────────┐
       │                       │                       │                 │
┌──────▼──────┐   ┌────────▼─────────┐   ┌──────▼───────┐   ┌─────▼────────┐
│  Users API  │   │   Rooms API      │   │ Bookings API │   │  Search API  │
│   (Go)      │   │   (Go)           │   │   (Go)       │   │   (Go)       │
│  Port 8080  │   │  Port 8081       │   │  Port 8083   │   │  Port 8082   │
└──────┬──────┘   └────────┬─────────┘   └──────┬───────┘   └─────┬────────┘
       │                   │                     │                 │
┌──────▼──────┐   ┌────────▼─────────┐   ┌──────▼───────┐   ┌─────▼────────┐
│   MySQL     │   │    MongoDB       │   │   MySQL      │   │  Solr + Mem  │
│ (users_db)  │   │   (rooms_db)     │   │(bookings_db) │   │  Port 8983   │
│  Port 3306  │   │   Port 27017     │   │  Port 3307   │   │Port 11211    │
└─────────────┘   └──────────────────┘   └──────────────┘   └──────────────┘
                               │
                        ┌──────▼──────┐
                        │  RabbitMQ   │
                        │ Port 5672   │
                        └─────────────┘
```


Estructura buscada:

hotel-reservation-system/
│
├── docker-compose.yml
├── README.md
├── .gitignore
│
├── users-api/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   │
│   ├── controllers/
│   │   └── user_controller.go
│   │
│   ├── services/
│   │   ├── user_service.go
│   │   └── user_service_test.go
│   │
│   ├── repositories/
│   │   └── user_repository.go
│   │
│   ├── domain/
│   │   ├── user.go
│   │   └── dto.go
│   │
│   ├── middlewares/
│   │   └── auth_middleware.go
│   │
│   ├── utils/
│   │   ├── jwt.go
│   │   ├── hash.go
│   │   └── errors.go
│   │
│   └── config/
│       └── database.go
│
├── rooms-api/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   │
│   ├── controllers/
│   │   └── room_controller.go
│   │
│   ├── services/
│   │   ├── room_service.go
│   │   └── room_service_test.go
│   │
│   ├── repositories/
│   │   └── room_repository.go
│   │
│   ├── domain/
│   │   ├── room.go
│   │   └── dto.go
│   │
│   ├── clients/
│   │   └── user_client.go
│   │
│   ├── queue/
│   │   └── rabbitmq_publisher.go
│   │
│   ├── utils/
│   │   ├── concurrent.go
│   │   └── errors.go
│   │
│   └── config/
│       ├── database.go
│       └── rabbitmq.go
│
├── bookings-api/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   │
│   ├── controllers/
│   │   └── booking_controller.go
│   │
│   ├── services/
│   │   ├── booking_service.go
│   │   └── booking_service_test.go
│   │
│   ├── repositories/
│   │   └── booking_repository.go
│   │
│   ├── domain/
│   │   ├── booking.go
│   │   └── dto.go
│   │
│   ├── clients/
│   │   ├── user_client.go
│   │   └── room_client.go
│   │
│   ├── middlewares/
│   │   └── auth_middleware.go
│   │
│   ├── utils/
│   │   └── errors.go
│   │
│   └── config/
│       └── database.go
│
├── search-api/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   │
│   ├── controllers/
│   │   └── search_controller.go
│   │
│   ├── services/
│   │   ├── search_service.go
│   │   └── search_service_test.go
│   │
│   ├── repositories/
│   │   └── solr_repository.go
│   │
│   ├── domain/
│   │   ├── room_index.go
│   │   └── dto.go
│   │
│   ├── clients/
│   │   └── room_client.go
│   │
│   ├── queue/
│   │   └── rabbitmq_consumer.go
│   │
│   ├── cache/
│   │   ├── local_cache.go
│   │   └── memcached_cache.go
│   │
│   ├── utils/
│   │   └── errors.go
│   │
│   └── config/
│       ├── solr.go
│       ├── cache.go
│       └── rabbitmq.go
│
└── frontend/
    ├── Dockerfile
    ├── package.json
    ├── package-lock.json
    ├── .env
    │
    ├── public/
    │   ├── index.html
    │   └── favicon.ico
    │
    └── src/
        ├── App.js
        ├── index.js
        │
        ├── components/
        │   ├── Navbar.jsx
        │   ├── RoomCard.jsx
        │   ├── SearchBar.jsx
        │   ├── PrivateRoute.jsx
        │   └── AdminRoute.jsx
        │
        ├── pages/
        │   ├── Login.jsx
        │   ├── Register.jsx
        │   ├── Home.jsx
        │   ├── RoomDetails.jsx
        │   ├── Congrats.jsx
        │   ├── MyBookings.jsx
        │   └── Admin.jsx
        │
        ├── services/
        │   ├── authService.js
        │   ├── roomService.js
        │   ├── searchService.js
        │   └── bookingService.js
        │
        ├── context/
        │   └── AuthContext.jsx
        │
        ├── utils/
        │   ├── api.js
        │   └── constants.js
        │
        └── styles/
            └── App.css

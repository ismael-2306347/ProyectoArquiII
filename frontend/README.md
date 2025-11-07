# Grand Prix Hotel - Frontend

Frontend moderno y funcional para el sistema de gestión hotelera Grand Prix, construido con React, TypeScript, TailwindCSS y Vite.

## Características

- Autenticación de usuarios con JWT
- Exploración de habitaciones con filtros avanzados
- Sistema de reservas intuitivo
- Panel de gestión de reservas personales
- Diseño responsive y moderno
- Interfaz de usuario moderna y elegante

## Tecnologías

- **React 18** - Biblioteca de UI
- **TypeScript** - Tipado estático
- **Vite** - Build tool y dev server
- **TailwindCSS** - Framework de CSS utility-first
- **React Router** - Enrutamiento
- **Axios** - Cliente HTTP
- **date-fns** - Manejo de fechas
- **Lucide React** - Iconos

## Requisitos Previos

- Node.js 18 o superior
- npm o yarn
- Backend ejecutándose (microservicios en puertos 8080, 8081, 8082)

## Instalación

1. Instalar dependencias:

```bash
cd frontend
npm install
```

2. Configurar el proxy (ya está configurado en `vite.config.ts`):

El proyecto está configurado para hacer proxy a:
- `http://localhost:8080` - users-api
- `http://localhost:8081` - rooms-api
- `http://localhost:8082` - reservations-api

## Ejecución

### Modo Desarrollo

```bash
npm run dev
```

La aplicación estará disponible en `http://localhost:3000`

### Build para Producción

```bash
npm run build
```

Los archivos optimizados se generarán en la carpeta `dist/`

### Preview de Producción

```bash
npm run preview
```

## Estructura del Proyecto

```
frontend/
├── src/
│   ├── components/          # Componentes reutilizables
│   │   ├── auth/           # Componentes de autenticación
│   │   ├── layout/         # Componentes de layout (Navbar, Layout)
│   │   └── ui/             # Componentes UI base (Button, Input, Card)
│   ├── context/            # Context API (AuthContext)
│   ├── lib/                # Configuraciones (axios)
│   ├── pages/              # Páginas de la aplicación
│   │   ├── Home.tsx
│   │   ├── Login.tsx
│   │   ├── Register.tsx
│   │   ├── Rooms.tsx
│   │   ├── ReserveRoom.tsx
│   │   └── MyReservations.tsx
│   ├── services/           # Servicios de API
│   │   ├── authService.ts
│   │   ├── roomService.ts
│   │   └── reservationService.ts
│   ├── types/              # Definiciones de TypeScript
│   ├── App.tsx             # Componente principal
│   ├── main.tsx            # Punto de entrada
│   └── index.css           # Estilos globales
├── public/                 # Archivos estáticos
├── package.json
├── vite.config.ts          # Configuración de Vite
├── tailwind.config.js      # Configuración de TailwindCSS
└── tsconfig.json           # Configuración de TypeScript
```

## Rutas de la Aplicación

### Rutas Públicas
- `/login` - Página de inicio de sesión
- `/register` - Página de registro

### Rutas Protegidas (requieren autenticación)
- `/` - Página principal
- `/rooms` - Explorar habitaciones
- `/rooms/:roomId/reserve` - Reservar habitación
- `/my-reservations` - Ver mis reservas

## Funcionalidades

### Autenticación
- Registro de usuarios con validación
- Inicio de sesión con usuario/email y contraseña
- Almacenamiento de token JWT en localStorage
- Protección de rutas con ProtectedRoute
- Logout

### Habitaciones
- Listar todas las habitaciones disponibles
- Filtros por:
  - Tipo de habitación
  - Rango de precios
  - Piso
  - Amenidades (WiFi, AC, TV, Minibar)
- Vista detallada de cada habitación
- Información completa de amenidades

### Reservas
- Selección de fechas de entrada y salida
- Cálculo automático de precio total
- Validación de fechas
- Confirmación de reserva
- Visualización de reservas activas

## Servicios de API

### authService
- `login(credentials)` - Iniciar sesión
- `register(userData)` - Registrar usuario
- `getUser(id)` - Obtener usuario por ID
- `logout()` - Cerrar sesión

### roomService
- `getAllRooms(filters)` - Obtener todas las habitaciones
- `getAvailableRooms(filters)` - Obtener habitaciones disponibles
- `getRoomById(id)` - Obtener habitación por ID
- `createRoom(data)` - Crear habitación (admin)
- `updateRoom(id, data)` - Actualizar habitación (admin)
- `deleteRoom(id)` - Eliminar habitación (admin)

### reservationService
- `createReservation(data)` - Crear reserva
- `getReservationById(id)` - Obtener reserva por ID
- `cancelReservation(id, reason)` - Cancelar reserva

## Variables de Entorno

Si necesitas cambiar las URLs de los servicios, puedes crear un archivo `.env`:

```env
VITE_USERS_API_URL=http://localhost:8080
VITE_ROOMS_API_URL=http://localhost:8081
VITE_RESERVATIONS_API_URL=http://localhost:8082
```

## Flujo de Uso

1. **Registro/Login**: El usuario se registra o inicia sesión
2. **Explorar Habitaciones**: Navega por las habitaciones disponibles con filtros
3. **Reservar**: Selecciona una habitación y elige fechas
4. **Confirmar**: Revisa el precio total y confirma la reserva
5. **Gestionar**: Ve y gestiona sus reservas activas

## Estilos y Diseño

El diseño sigue un estilo moderno y elegante con:
- Paleta de colores azul (primary-600 como color principal)
- Componentes con sombras suaves
- Bordes redondeados
- Tipografía clara y legible
- Diseño responsive (mobile-first)
- Animaciones sutiles

## Solución de Problemas

### Error de CORS
Si encuentras errores de CORS, asegúrate de que el backend tenga CORS habilitado para `http://localhost:3000`

### Error de conexión a APIs
Verifica que los tres microservicios estén ejecutándose:
```bash
# Desde la raíz del proyecto
docker-compose up
```

### Errores de compilación de TypeScript
```bash
# Limpia y reinstala dependencias
rm -rf node_modules package-lock.json
npm install
```

## Próximas Mejoras

- [ ] Implementar búsqueda de habitaciones por fechas disponibles
- [ ] Agregar panel de administración para gestionar habitaciones
- [ ] Implementar notificaciones en tiempo real con WebSockets
- [ ] Agregar galería de imágenes de habitaciones
- [ ] Implementar sistema de reviews y calificaciones
- [ ] Agregar exportación de reservas a PDF

## Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia

Este proyecto es parte del curso de Arquitectura de Software II.

## Contacto

Para dudas o sugerencias sobre el frontend, por favor abre un issue en el repositorio.

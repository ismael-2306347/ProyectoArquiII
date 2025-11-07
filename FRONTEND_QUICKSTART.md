# Inicio Rápido - Frontend Grand Prix Hotel

Este documento contiene las instrucciones para poner en marcha el frontend de Grand Prix Hotel.

## Prerrequisitos

Asegúrate de tener el backend ejecutándose primero:

```bash
# En la raíz del proyecto
docker-compose up
```

Esto levantará los 3 microservicios:
- users-api en puerto 8080
- rooms-api en puerto 8081
- reservations-api en puerto 8082

## Instalación del Frontend

1. Navega a la carpeta del frontend:

```bash
cd frontend
```

2. Instala las dependencias:

```bash
npm install
```

Si encuentras errores durante la instalación, intenta:

```bash
npm install --legacy-peer-deps
```

## Ejecutar el Frontend

Una vez instaladas las dependencias, ejecuta:

```bash
npm run dev
```

El frontend estará disponible en: **http://localhost:3000**

## Primeros Pasos en la Aplicación

### 1. Registrar un Usuario

- Abre http://localhost:3000
- Serás redirigido a `/login`
- Haz clic en "Regístrate aquí"
- Completa el formulario de registro:
  - Nombre y apellido
  - Nombre de usuario (único)
  - Email (único)
  - Contraseña (mínimo 6 caracteres)
- Haz clic en "Registrarse"
- Serás automáticamente autenticado y redirigido al Home

### 2. Explorar Habitaciones

- En el navbar, haz clic en "Habitaciones"
- Verás un listado de habitaciones disponibles
- Puedes filtrar por:
  - Tipo de habitación
  - Precio mínimo y máximo
  - Piso
  - Amenidades (WiFi, AC, TV, Minibar)

### 3. Hacer una Reserva

- En la página de habitaciones, haz clic en "Reservar" en cualquier habitación
- Selecciona la fecha de entrada y salida
- Verás el cálculo automático del precio total
- Haz clic en "Confirmar Reserva"
- Serás redirigido a "Mis Reservas"

### 4. Ver tus Reservas

- En el navbar, haz clic en "Mis Reservas"
- Verás todas tus reservas activas
- (Nota: actualmente el endpoint de listar reservas por usuario puede no estar implementado en el backend)

## Estructura del Código

```
frontend/
├── src/
│   ├── components/
│   │   ├── auth/           # ProtectedRoute
│   │   ├── layout/         # Navbar, Layout
│   │   └── ui/             # Button, Input, Card
│   ├── context/            # AuthContext (manejo de autenticación)
│   ├── lib/                # Configuración de axios
│   ├── pages/              # Páginas de la app
│   ├── services/           # Servicios para llamar a las APIs
│   ├── types/              # Tipos de TypeScript
│   ├── App.tsx             # Rutas principales
│   └── main.tsx            # Entry point
```

## Características Principales

### Autenticación
- JWT almacenado en localStorage
- Rutas protegidas con ProtectedRoute
- Logout limpia el estado

### Diseño
- TailwindCSS para estilos
- Componentes reutilizables (Button, Input, Card)
- Diseño responsive
- Tema azul (TUTE style)

### API Integration
- Axios configurado con interceptors
- Proxy de Vite para las 3 APIs
- Manejo de errores centralizado

## Solución de Problemas

### Error: "Cannot connect to API"

Verifica que el backend esté corriendo:
```bash
docker-compose ps
```

Deberías ver 3 servicios activos: users-api, rooms-api, reservations-api

### Error: "Module not found"

Reinstala las dependencias:
```bash
rm -rf node_modules package-lock.json
npm install
```

### Error de CORS

Si ves errores de CORS en la consola del navegador, asegúrate de que los servicios backend tengan CORS habilitado para http://localhost:3000

Puedes agregar esto en tus controladores de Go (Gin):

```go
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

### Las habitaciones no se cargan

Verifica que tengas habitaciones en la base de datos. Puedes crearlas usando:

```bash
curl -X POST http://localhost:8081/api/v1/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "number": "101",
    "type": "double",
    "price": 100.00,
    "description": "Habitación doble con vista al mar",
    "capacity": 2,
    "floor": 1,
    "has_wifi": true,
    "has_ac": true,
    "has_tv": true,
    "has_minibar": true
  }'
```

## Próximos Pasos

1. Crear más habitaciones en el backend
2. Probar el flujo completo de reservas
3. Verificar que las reservas se almacenan correctamente en MongoDB
4. Personalizar los estilos y colores según tus preferencias

## Documentación Adicional

Para más detalles, consulta:
- [frontend/README.md](frontend/README.md) - Documentación completa del frontend
- [README.md](README.md) - Documentación general del proyecto

## Stack Tecnológico Frontend

- React 18
- TypeScript
- Vite
- TailwindCSS
- React Router
- Axios
- date-fns
- Lucide React (iconos)

Disfruta desarrollando con TUTE Hotel Frontend! 🏨

-- Crear base de datos para users-api
CREATE DATABASE IF NOT EXISTS usersdb;

-- Crear base de datos para rooms-api
CREATE DATABASE IF NOT EXISTS roomsdb;

-- Otorgar permisos al usuario 'user'
GRANT ALL PRIVILEGES ON usersdb.* TO 'user'@'%';
GRANT ALL PRIVILEGES ON roomsdb.* TO 'user'@'%';

FLUSH PRIVILEGES;

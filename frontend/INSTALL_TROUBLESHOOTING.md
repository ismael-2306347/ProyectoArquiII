# Solución de Problemas de Instalación - Frontend

## Problema: Error "patch-package" durante npm install

Este error ocurre porque hay archivos bloqueados por otro proceso (VSCode, antivirus, etc.).

## Solución 1: Limpiar e Instalar (RECOMENDADO)

### Pasos:

1. **Cierra COMPLETAMENTE VSCode** y cualquier otra terminal abierta

2. **Abre PowerShell como ADMINISTRADOR**
   - Presiona `Win + X`
   - Selecciona "Windows PowerShell (Administrador)" o "Terminal (Admin)"

3. **Navega a la carpeta del frontend:**
   ```powershell
   cd C:\Users\manub\OneDrive\Escritorio\ProyectoArquiII\frontend
   ```

4. **Elimina node_modules y archivos de lock:**
   ```powershell
   # Opción 1: Con PowerShell
   Remove-Item -Recurse -Force node_modules -ErrorAction SilentlyContinue
   Remove-Item -Force package-lock.json -ErrorAction SilentlyContinue

   # Opción 2: Con CMD (si la opción 1 falla)
   cmd /c "rd /s /q node_modules"
   cmd /c "del /f package-lock.json"
   ```

5. **Instala las dependencias:**
   ```powershell
   npm install
   ```

6. **Si aún falla, prueba con:**
   ```powershell
   npm install --legacy-peer-deps
   ```

## Solución 2: Usar Yarn (Alternativa)

Si npm sigue fallando, prueba con yarn:

```powershell
# 1. Instalar yarn globalmente
npm install -g yarn

# 2. Ir a la carpeta del frontend
cd C:\Users\manub\OneDrive\Escritorio\ProyectoArquiII\frontend

# 3. Instalar con yarn
yarn install
```

## Solución 3: Desactivar Antivirus Temporalmente

A veces Windows Defender o antivirus de terceros bloquean archivos:

1. Desactiva temporalmente el antivirus
2. Sigue los pasos de la Solución 1
3. Reactiva el antivirus después de la instalación

## Solución 4: Cambiar Ubicación del Proyecto

OneDrive puede causar problemas de sincronización:

1. Copia el proyecto a una ubicación local (ej: `C:\Projects\ProyectoArquiII`)
2. Sigue los pasos de la Solución 1 en la nueva ubicación

## Verificar que la Instalación Funcionó

Después de la instalación exitosa, verifica:

```powershell
# Ver que node_modules existe
ls node_modules

# Intentar ejecutar el servidor de desarrollo
npm run dev
```

Si `npm run dev` inicia correctamente, verás algo como:

```
  VITE v5.0.8  ready in 500 ms

  ➜  Local:   http://localhost:3000/
  ➜  Network: use --host to expose
  ➜  press h + enter to show help
```

## Errores Comunes Durante la Ejecución

### Error: "Cannot find module '@/...'

Si ves errores de módulos no encontrados, reinicia el servidor:

```powershell
# Detén el servidor (Ctrl + C)
# Vuelve a iniciarlo
npm run dev
```

### Error: "Failed to resolve import"

Asegúrate de que TypeScript esté configurado correctamente. El archivo `tsconfig.app.json` ya debería tener los path aliases configurados.

### Error de CORS

Si al hacer requests al backend ves errores de CORS, verifica:

1. Que el backend esté corriendo:
   ```bash
   docker-compose ps
   ```

2. Que el proxy esté configurado en `vite.config.ts` (ya está configurado)

## Archivos Importantes Creados

- `.npmrc` - Configuración de npm con `legacy-peer-deps`
- `package.json` - Dependencias actualizadas y compatibles
- `vite.config.ts` - Proxy configurado para las APIs

## Si Nada Funciona

Como última opción, puedes:

1. Eliminar la carpeta `frontend` completa
2. Recrearla con:
   ```bash
   npm create vite@latest frontend -- --template react-ts
   ```
3. Copiar SOLO el contenido de `src/` desde el backup
4. Copiar los archivos de configuración (tailwind.config.js, vite.config.ts, etc.)
5. Reinstalar dependencias

## Contacto

Si el problema persiste, busca ayuda proporcionando:
- El log completo del error (`npm-debug.log`)
- Tu versión de Node.js: `node --version`
- Tu versión de npm: `npm --version`
- Tu sistema operativo

---

**Última actualización:** 2025-11-07

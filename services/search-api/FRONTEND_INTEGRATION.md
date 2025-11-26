# Integración del Frontend con Search API

Esta guía muestra cómo integrar el frontend React con la search-api para implementar búsqueda de habitaciones.

## 1. Servicio de Búsqueda (TypeScript)

Crea un servicio para interactuar con search-api:

```typescript
// src/services/searchService.ts
import axios from 'axios';

const SEARCH_API_URL = import.meta.env.VITE_SEARCH_API_URL || 'http://localhost:8083';

export interface SearchFilters {
  q?: string;
  type?: 'single' | 'double' | 'suite' | 'deluxe' | 'standard';
  status?: 'available' | 'occupied' | 'maintenance' | 'reserved';
  floor?: number;
  minPrice?: number;
  maxPrice?: number;
  hasWifi?: boolean;
  hasAC?: boolean;
  hasTV?: boolean;
  hasMinibar?: boolean;
  sort?: string;
  page?: number;
  limit?: number;
}

export interface RoomSearchResult {
  id: string;
  number: string;
  type: string;
  status: string;
  price: number;
  capacity: number;
  floor: number;
  has_wifi: boolean;
  has_ac: boolean;
  has_tv: boolean;
  has_minibar: boolean;
}

export interface SearchResponse {
  page: number;
  limit: number;
  total: number;
  results: RoomSearchResult[];
}

class SearchService {
  async searchRooms(filters: SearchFilters): Promise<SearchResponse> {
    try {
      const params = new URLSearchParams();

      // Agregar filtros a params
      if (filters.q) params.append('q', filters.q);
      if (filters.type) params.append('type', filters.type);
      if (filters.status) params.append('status', filters.status);
      if (filters.floor !== undefined) params.append('floor', filters.floor.toString());
      if (filters.minPrice !== undefined) params.append('min_price', filters.minPrice.toString());
      if (filters.maxPrice !== undefined) params.append('max_price', filters.maxPrice.toString());
      if (filters.hasWifi !== undefined) params.append('has_wifi', filters.hasWifi.toString());
      if (filters.hasAC !== undefined) params.append('has_ac', filters.hasAC.toString());
      if (filters.hasTV !== undefined) params.append('has_tv', filters.hasTV.toString());
      if (filters.hasMinibar !== undefined) params.append('has_minibar', filters.hasMinibar.toString());
      if (filters.sort) params.append('sort', filters.sort);
      if (filters.page) params.append('page', filters.page.toString());
      if (filters.limit) params.append('limit', filters.limit.toString());

      const response = await axios.get<SearchResponse>(
        `${SEARCH_API_URL}/api/search/rooms?${params.toString()}`
      );

      return response.data;
    } catch (error) {
      console.error('Error searching rooms:', error);
      throw error;
    }
  }
}

export default new SearchService();
```

## 2. Hook de React para Búsqueda

Crea un custom hook para manejar la lógica de búsqueda:

```typescript
// src/hooks/useRoomSearch.ts
import { useState, useCallback } from 'react';
import searchService, { SearchFilters, SearchResponse } from '../services/searchService';

export const useRoomSearch = () => {
  const [results, setResults] = useState<SearchResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const search = useCallback(async (filters: SearchFilters) => {
    setLoading(true);
    setError(null);

    try {
      const response = await searchService.searchRooms(filters);
      setResults(response);
    } catch (err) {
      setError('Error al buscar habitaciones. Por favor intenta de nuevo.');
      console.error('Search error:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  return { results, loading, error, search };
};
```

## 3. Componente de Búsqueda

Crea un componente de búsqueda con filtros:

```tsx
// src/components/RoomSearch.tsx
import React, { useState } from 'react';
import { useRoomSearch } from '../hooks/useRoomSearch';
import { SearchFilters } from '../services/searchService';

const RoomSearch: React.FC = () => {
  const { results, loading, error, search } = useRoomSearch();

  const [filters, setFilters] = useState<SearchFilters>({
    status: 'available',
    page: 1,
    limit: 10,
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    search(filters);
  };

  const handleFilterChange = (key: keyof SearchFilters, value: any) => {
    setFilters(prev => ({ ...prev, [key]: value, page: 1 }));
  };

  const handlePageChange = (newPage: number) => {
    setFilters(prev => ({ ...prev, page: newPage }));
    search({ ...filters, page: newPage });
  };

  return (
    <div className="room-search">
      <h2 className="text-2xl font-bold mb-6">Buscar Habitaciones</h2>

      {/* Formulario de Filtros */}
      <form onSubmit={handleSearch} className="bg-white p-6 rounded-lg shadow-md mb-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">

          {/* Búsqueda de texto */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Búsqueda
            </label>
            <input
              type="text"
              placeholder="Número, tipo, descripción..."
              className="w-full px-3 py-2 border rounded-md"
              value={filters.q || ''}
              onChange={(e) => handleFilterChange('q', e.target.value)}
            />
          </div>

          {/* Tipo de habitación */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Tipo
            </label>
            <select
              className="w-full px-3 py-2 border rounded-md"
              value={filters.type || ''}
              onChange={(e) => handleFilterChange('type', e.target.value || undefined)}
            >
              <option value="">Todos</option>
              <option value="single">Individual</option>
              <option value="double">Doble</option>
              <option value="suite">Suite</option>
              <option value="deluxe">Deluxe</option>
              <option value="standard">Estándar</option>
            </select>
          </div>

          {/* Estado */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Estado
            </label>
            <select
              className="w-full px-3 py-2 border rounded-md"
              value={filters.status || ''}
              onChange={(e) => handleFilterChange('status', e.target.value || undefined)}
            >
              <option value="">Todos</option>
              <option value="available">Disponible</option>
              <option value="occupied">Ocupada</option>
              <option value="maintenance">Mantenimiento</option>
              <option value="reserved">Reservada</option>
            </select>
          </div>

          {/* Rango de precio */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Precio Mínimo
            </label>
            <input
              type="number"
              placeholder="0"
              className="w-full px-3 py-2 border rounded-md"
              value={filters.minPrice || ''}
              onChange={(e) => handleFilterChange('minPrice', e.target.value ? parseFloat(e.target.value) : undefined)}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Precio Máximo
            </label>
            <input
              type="number"
              placeholder="9999"
              className="w-full px-3 py-2 border rounded-md"
              value={filters.maxPrice || ''}
              onChange={(e) => handleFilterChange('maxPrice', e.target.value ? parseFloat(e.target.value) : undefined)}
            />
          </div>

          {/* Piso */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Piso
            </label>
            <input
              type="number"
              placeholder="Cualquiera"
              className="w-full px-3 py-2 border rounded-md"
              value={filters.floor || ''}
              onChange={(e) => handleFilterChange('floor', e.target.value ? parseInt(e.target.value) : undefined)}
            />
          </div>
        </div>

        {/* Amenities (checkboxes) */}
        <div className="mt-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Servicios
          </label>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            <label className="flex items-center space-x-2">
              <input
                type="checkbox"
                checked={filters.hasWifi || false}
                onChange={(e) => handleFilterChange('hasWifi', e.target.checked || undefined)}
              />
              <span>WiFi</span>
            </label>
            <label className="flex items-center space-x-2">
              <input
                type="checkbox"
                checked={filters.hasAC || false}
                onChange={(e) => handleFilterChange('hasAC', e.target.checked || undefined)}
              />
              <span>Aire Acondicionado</span>
            </label>
            <label className="flex items-center space-x-2">
              <input
                type="checkbox"
                checked={filters.hasTV || false}
                onChange={(e) => handleFilterChange('hasTV', e.target.checked || undefined)}
              />
              <span>TV</span>
            </label>
            <label className="flex items-center space-x-2">
              <input
                type="checkbox"
                checked={filters.hasMinibar || false}
                onChange={(e) => handleFilterChange('hasMinibar', e.target.checked || undefined)}
              />
              <span>Minibar</span>
            </label>
          </div>
        </div>

        {/* Ordenamiento */}
        <div className="mt-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Ordenar por
          </label>
          <select
            className="w-full md:w-64 px-3 py-2 border rounded-md"
            value={filters.sort || ''}
            onChange={(e) => handleFilterChange('sort', e.target.value || undefined)}
          >
            <option value="">Predeterminado</option>
            <option value="price">Precio (menor a mayor)</option>
            <option value="-price">Precio (mayor a menor)</option>
            <option value="floor">Piso (menor a mayor)</option>
            <option value="-floor">Piso (mayor a menor)</option>
            <option value="capacity">Capacidad (menor a mayor)</option>
            <option value="-capacity">Capacidad (mayor a menor)</option>
          </select>
        </div>

        <button
          type="submit"
          className="mt-4 px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          disabled={loading}
        >
          {loading ? 'Buscando...' : 'Buscar'}
        </button>
      </form>

      {/* Mensajes de error */}
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      {/* Resultados */}
      {results && (
        <div>
          <div className="mb-4 text-gray-600">
            Mostrando {results.results.length} de {results.total} resultados
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {results.results.map((room) => (
              <div key={room.id} className="bg-white p-4 rounded-lg shadow-md">
                <div className="flex justify-between items-start mb-2">
                  <h3 className="text-lg font-bold">Habitación {room.number}</h3>
                  <span className={`px-2 py-1 rounded text-xs ${
                    room.status === 'available' ? 'bg-green-100 text-green-800' :
                    room.status === 'occupied' ? 'bg-red-100 text-red-800' :
                    'bg-gray-100 text-gray-800'
                  }`}>
                    {room.status === 'available' ? 'Disponible' :
                     room.status === 'occupied' ? 'Ocupada' :
                     room.status === 'maintenance' ? 'Mantenimiento' : 'Reservada'}
                  </span>
                </div>

                <div className="text-sm text-gray-600 mb-3">
                  <div>Tipo: <span className="font-medium">{room.type}</span></div>
                  <div>Piso: {room.floor}</div>
                  <div>Capacidad: {room.capacity} personas</div>
                  <div className="text-lg font-bold text-blue-600 mt-2">
                    ${room.price.toFixed(2)} / noche
                  </div>
                </div>

                <div className="flex gap-2 flex-wrap">
                  {room.has_wifi && <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">WiFi</span>}
                  {room.has_ac && <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">AC</span>}
                  {room.has_tv && <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">TV</span>}
                  {room.has_minibar && <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">Minibar</span>}
                </div>

                {room.status === 'available' && (
                  <button className="mt-4 w-full px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700">
                    Reservar
                  </button>
                )}
              </div>
            ))}
          </div>

          {/* Paginación */}
          {results.total > results.limit && (
            <div className="mt-6 flex justify-center gap-2">
              <button
                onClick={() => handlePageChange(results.page - 1)}
                disabled={results.page === 1}
                className="px-4 py-2 border rounded disabled:opacity-50"
              >
                Anterior
              </button>

              <span className="px-4 py-2">
                Página {results.page} de {Math.ceil(results.total / results.limit)}
              </span>

              <button
                onClick={() => handlePageChange(results.page + 1)}
                disabled={results.page >= Math.ceil(results.total / results.limit)}
                className="px-4 py-2 border rounded disabled:opacity-50"
              >
                Siguiente
              </button>
            </div>
          )}
        </div>
      )}

      {/* Estado inicial */}
      {!results && !loading && (
        <div className="text-center text-gray-500 py-12">
          Usa los filtros de arriba para buscar habitaciones
        </div>
      )}
    </div>
  );
};

export default RoomSearch;
```

## 4. Configuración de Variables de Entorno

```bash
# .env
VITE_SEARCH_API_URL=http://localhost:8083
```

Para producción:
```bash
VITE_SEARCH_API_URL=https://api.tuhotel.com/search
```

## 5. Agregar Ruta en el Router

```typescript
// src/App.tsx o router config
import RoomSearch from './components/RoomSearch';

// Dentro de tus rutas:
<Route path="/search" element={<RoomSearch />} />
```

## 6. Optimizaciones Adicionales

### Debounce para búsqueda de texto

```typescript
// src/hooks/useDebounce.ts
import { useEffect, useState } from 'react';

export function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}
```

Uso en el componente:
```typescript
const [searchQuery, setSearchQuery] = useState('');
const debouncedQuery = useDebounce(searchQuery, 500);

useEffect(() => {
  if (debouncedQuery) {
    search({ ...filters, q: debouncedQuery });
  }
}, [debouncedQuery]);
```

### Cache en el frontend

```typescript
// Usar React Query para cache automático
import { useQuery } from '@tanstack/react-query';

export const useRoomSearchQuery = (filters: SearchFilters) => {
  return useQuery({
    queryKey: ['rooms', filters],
    queryFn: () => searchService.searchRooms(filters),
    staleTime: 60000, // 1 minuto
  });
};
```

## 7. Ejemplo de Uso Completo

```tsx
// src/pages/SearchPage.tsx
import React, { useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import RoomSearch from '../components/RoomSearch';

const SearchPage: React.FC = () => {
  const [searchParams] = useSearchParams();

  // Puedes pre-cargar filtros desde URL params
  useEffect(() => {
    const type = searchParams.get('type');
    const minPrice = searchParams.get('min_price');
    // etc...
  }, [searchParams]);

  return (
    <div className="container mx-auto px-4 py-8">
      <RoomSearch />
    </div>
  );
};

export default SearchPage;
```

## 8. Testing

```typescript
// src/services/__tests__/searchService.test.ts
import { describe, it, expect, vi } from 'vitest';
import axios from 'axios';
import searchService from '../searchService';

vi.mock('axios');

describe('SearchService', () => {
  it('should search rooms with filters', async () => {
    const mockResponse = {
      data: {
        page: 1,
        limit: 10,
        total: 1,
        results: [
          {
            id: '1',
            number: '101',
            type: 'double',
            status: 'available',
            price: 120.50,
            capacity: 2,
            floor: 1,
            has_wifi: true,
            has_ac: true,
            has_tv: true,
            has_minibar: false,
          },
        ],
      },
    };

    vi.mocked(axios.get).mockResolvedValue(mockResponse);

    const result = await searchService.searchRooms({ status: 'available' });

    expect(result.results).toHaveLength(1);
    expect(result.results[0].number).toBe('101');
  });
});
```

## 9. Integración con Sistema de Reservas

Cuando el usuario haga clic en "Reservar", puedes redirigirlo al flujo de reservas:

```tsx
const handleReserve = (roomId: string) => {
  navigate(`/reservations/new?room_id=${roomId}`);
};
```

O abrir un modal:

```tsx
const [selectedRoom, setSelectedRoom] = useState<RoomSearchResult | null>(null);

// En el botón de reservar:
<button onClick={() => setSelectedRoom(room)}>
  Reservar
</button>

// Modal:
{selectedRoom && (
  <ReservationModal
    room={selectedRoom}
    onClose={() => setSelectedRoom(null)}
  />
)}
```

## 10. Características Avanzadas

### Autocompletado

```tsx
const [suggestions, setSuggestions] = useState<string[]>([]);

const handleAutocomplete = async (query: string) => {
  if (query.length < 2) return;

  // Buscar con límite pequeño solo para autocompletar
  const results = await searchService.searchRooms({ q: query, limit: 5 });
  const numbers = results.results.map(r => r.number);
  setSuggestions(numbers);
};
```

### Filtros guardados

```tsx
const saveFilters = () => {
  localStorage.setItem('savedSearch', JSON.stringify(filters));
};

const loadFilters = () => {
  const saved = localStorage.getItem('savedSearch');
  if (saved) {
    setFilters(JSON.parse(saved));
  }
};
```

### Compartir búsqueda

```tsx
const shareSearch = () => {
  const params = new URLSearchParams(filters as any);
  const url = `${window.location.origin}/search?${params.toString()}`;
  navigator.clipboard.writeText(url);
};
```

---

Con esta integración, tu frontend tendrá búsqueda de habitaciones rápida, filtrable y en tiempo real!

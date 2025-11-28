import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Layout } from '@/components/layout/Layout';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { useAuth } from '@/context/AuthContext';
import type { Room, RoomFilter, RoomType } from '@/types';
import { Wifi, Tv, Wind, Coffee, Users, MapPin, DollarSign } from 'lucide-react';
import { roomService } from '@/services/roomService';

type RoomSearchFilter = RoomFilter & {
  q?: string;
};

export function Rooms() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [rooms, setRooms] = useState<Room[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [filters, setFilters] = useState<RoomSearchFilter>({
    page: 1,
    limit: 12,
    q: '',
  });

  useEffect(() => {
    fetchRooms();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filters]);

  const fetchRooms = async () => {
    setIsLoading(true);
    setError('');
    try {
      // 👉 Ahora pegamos contra rooms-api (getAvailableRooms)
      const response = await roomService.getAvailableRooms(filters);
      setRooms(response.rooms || []);
    } catch (err: any) {
      const message =
        err?.response?.data?.error || 'Error al cargar las habitaciones';
      setError(message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleFilterChange = (key: keyof RoomSearchFilter, value: any) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value || undefined,
      page: 1,
    }));
  };

  const handleReserve = (room: Room) => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    navigate(`/rooms/${room.id}/reserve`);
  };

  const getRoomTypeLabel = (type: RoomType) => {
    const labels: Record<RoomType, string> = {
      single: 'Individual',
      double: 'Doble',
      suite: 'Suite',
      deluxe: 'De Lujo',
      standard: 'Estándar',
    };

    return labels[type];
  };

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <h1 className="text-3xl font-bold text-gray-900">
            Habitaciones Disponibles
          </h1>
        </div>

        {/* Filters */}
        <Card>
          <CardHeader>
            <CardTitle>Filtros</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Buscar
                </label>
                <Input
                  type="text"
                  placeholder="Número, tipo, descripción..."
                  value={filters.q || ''}
                  onChange={(e) =>
                    handleFilterChange('q', e.target.value || undefined)
                  }
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Tipo de Habitación
                </label>
                <select
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                  value={filters.type || ''}
                  onChange={(e) => handleFilterChange('type', e.target.value)}
                >
                  <option value="">Todos</option>
                  <option value="single">Individual</option>
                  <option value="double">Doble</option>
                  <option value="suite">Suite</option>
                  <option value="deluxe">De Lujo</option>
                  <option value="standard">Estándar</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Precio Mínimo
                </label>
                <Input
                  type="number"
                  placeholder="$0"
                  value={filters.min_price || ''}
                  onChange={(e) =>
                    handleFilterChange(
                      'min_price',
                      e.target.value ? parseFloat(e.target.value) : undefined
                    )
                  }
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Precio Máximo
                </label>
                <Input
                  type="number"
                  placeholder="$999"
                  value={filters.max_price || ''}
                  onChange={(e) =>
                    handleFilterChange(
                      'max_price',
                      e.target.value ? parseFloat(e.target.value) : undefined
                    )
                  }
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Piso
                </label>
                <Input
                  type="number"
                  placeholder="Todos"
                  value={filters.floor || ''}
                  onChange={(e) =>
                    handleFilterChange(
                      'floor',
                      e.target.value ? parseInt(e.target.value, 10) : undefined
                    )
                  }
                />
              </div>
            </div>

            <div className="mt-4 flex gap-4 flex-wrap">
              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  className="rounded text-primary-600 focus:ring-primary-500"
                  checked={filters.has_wifi || false}
                  onChange={(e) =>
                    handleFilterChange(
                      'has_wifi',
                      e.target.checked || undefined
                    )
                  }
                />
                <span className="text-sm text-gray-700">WiFi</span>
              </label>

              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  className="rounded text-primary-600 focus:ring-primary-500"
                  checked={filters.has_ac || false}
                  onChange={(e) =>
                    handleFilterChange(
                      'has_ac',
                      e.target.checked || undefined
                    )
                  }
                />
                <span className="text-sm text-gray-700">
                  Aire Acondicionado
                </span>
              </label>

              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  className="rounded text-primary-600 focus:ring-primary-500"
                  checked={filters.has_tv || false}
                  onChange={(e) =>
                    handleFilterChange(
                      'has_tv',
                      e.target.checked || undefined
                    )
                  }
                />
                <span className="text-sm text-gray-700">TV</span>
              </label>

              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  className="rounded text-primary-600 focus:ring-primary-500"
                  checked={filters.has_minibar || false}
                  onChange={(e) =>
                    handleFilterChange(
                      'has_minibar',
                      e.target.checked || undefined
                    )
                  }
                />
                <span className="text-sm text-gray-700">Minibar</span>
              </label>
            </div>
          </CardContent>
        </Card>

        {/* Rooms Grid */}
        {isLoading ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
            <p className="mt-4 text-gray-600">Cargando habitaciones...</p>
          </div>
        ) : error ? (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            {error}
          </div>
        ) : rooms.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-gray-600">
                No se encontraron habitaciones con los filtros seleccionados.
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {rooms.map((room) => (
              <Card
                key={room.id}
                className="hover:shadow-lg transition-shadow"
              >
                <CardContent className="p-0">
                  <div className="h-48 bg-gradient-to-br from-primary-400 to-primary-600 rounded-t-lg flex items-center justify-center">
                    <div className="text-center text-white">
                      <div className="text-5xl font-bold">#{room.number}</div>
                      <div className="text-lg mt-2">
                        {getRoomTypeLabel(room.type)}
                      </div>
                    </div>
                  </div>

                  <div className="p-6 space-y-4">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2 text-gray-600">
                        <Users className="w-4 h-4" />
                        <span className="text-sm">
                          {room.capacity} personas
                        </span>
                      </div>
                      <div className="flex items-center space-x-2 text-gray-600">
                        <MapPin className="w-4 h-4" />
                        <span className="text-sm">Piso {room.floor}</span>
                      </div>
                    </div>

                    <p className="text-gray-600 text-sm line-clamp-2">
                      {room.description || 'Habitación cómoda y acogedora.'}
                    </p>

                    <div className="flex gap-2 flex-wrap">
                      {room.has_wifi && (
                        <div className="flex items-center space-x-1 text-xs bg-gray-100 px-2 py-1 rounded">
                          <Wifi className="w-3 h-3" />
                          <span>WiFi</span>
                        </div>
                      )}
                      {room.has_tv && (
                        <div className="flex items-center space-x-1 text-xs bg-gray-100 px-2 py-1 rounded">
                          <Tv className="w-3 h-3" />
                          <span>TV</span>
                        </div>
                      )}
                      {room.has_ac && (
                        <div className="flex items-center space-x-1 text-xs bg-gray-100 px-2 py-1 rounded">
                          <Wind className="w-3 h-3" />
                          <span>A/C</span>
                        </div>
                      )}
                      {room.has_minibar && (
                        <div className="flex items-center space-x-1 text-xs bg-gray-100 px-2 py-1 rounded">
                          <Coffee className="w-3 h-3" />
                          <span>Minibar</span>
                        </div>
                      )}
                    </div>

                    <div className="border-t pt-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-1">
                          <DollarSign className="w-5 h-5 text-primary-600" />
                          <span className="text-2xl font-bold text-gray-900">
                            {room.price.toFixed(2)}
                          </span>
                          <span className="text-gray-600 text-sm">
                            /noche
                          </span>
                        </div>
                        <Button
                          variant="primary"
                          size="sm"
                          onClick={() => handleReserve(room)}
                        >
                          Reservar
                        </Button>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </Layout>
  );
}

import React, { useState, useEffect } from 'react';
import { Layout } from '@/components/layout/Layout';
import { Card, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Calendar, MapPin, DollarSign, X, AlertCircle, Loader2 } from 'lucide-react';
import { reservationService } from '@/services/reservationService';
import type { Reservation } from '@/types';
import { format } from 'date-fns';
import { es } from 'date-fns/locale';
import { Rooms } from './Rooms';
import { roomService } from '@/services/roomService';
import { useAuth } from '@/context/AuthContext';

interface ReservationWithDetails extends Reservation {
  roomNumber?: string;
  roomType?: string;
  floor?: number;
  guest_name?: string;
  guest_email?: string;
  guest_phone?: string;
}

export function MyReservations() {
  const { user } = useAuth();
  const [reservations, setReservations] = useState<ReservationWithDetails[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [cancellingId, setCancellingId] = useState<string | null>(null);
  
  useEffect(() => {
    if (user?.id) {
      fetchReservations();
    }
  }, [user]);

  const fetchReservations = async () => {
    if (!user?.id) return;

    setLoading(true);
    setError(null);
    try {
      const data = await reservationService.getMyReservations(user.id.toString());
      setReservations(data.reservations);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Error al cargar las reservas');
      console.error('Error fetching reservations:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleCancelReservation = async (reservationId: string) => {
    if (!confirm('¿Estás seguro de que deseas cancelar esta reserva?')) {
      return;
    }

    const reason = prompt('Por favor, indica el motivo de la cancelación:');
    if (!reason) {
      return;
    }

    setCancellingId(reservationId);

    try {
      await reservationService.cancelReservation(reservationId, reason);
      
      // Actualizar la lista de reservas
      await fetchReservations();
      
      alert('Reserva cancelada exitosamente');
    } catch (err: any) {
      alert(err.response?.data?.error || 'Error al cancelar la reserva');
      console.error('Error cancelling reservation:', err);
    } finally {
      setCancellingId(null);
    }
  };

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { bg: string; text: string; label: string }> = {
      confirmed: { bg: 'bg-green-100', text: 'text-green-700', label: 'Confirmada' },
      pending: { bg: 'bg-yellow-100', text: 'text-yellow-700', label: 'Pendiente' },
      cancelled: { bg: 'bg-red-100', text: 'text-red-700', label: 'Cancelada' },
      completed: { bg: 'bg-gray-100', text: 'text-gray-700', label: 'Completada' },
    };

    const config = statusConfig[status] || statusConfig.pending;

    return (
      <span className={`px-3 py-1 ${config.bg} ${config.text} rounded-full text-sm font-medium`}>
        {config.label}
      </span>
    );
  };

  const formatDate = (dateString: string) => {
    try {
      return format(new Date(dateString), "d 'de' MMMM, yyyy", { locale: es });
    } catch {
      return dateString;
    }
  };

  const calculateNights = (checkIn: string, checkOut: string) => {
    try {
      const checkInDate = new Date(checkIn);
      const checkOutDate = new Date(checkOut);
      const diffTime = Math.abs(checkOutDate.getTime() - checkInDate.getTime());
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
      return diffDays;
    } catch {
      return 0;
    }
  };

  if (!user) {
    return (
      <Layout>
        <Card>
          <CardContent className="py-12">
            <div className="text-center text-gray-600">
              <AlertCircle className="w-16 h-16 mx-auto mb-4 text-gray-400" />
              <p className="text-lg font-medium mb-2">Debes iniciar sesión</p>
              <p className="text-sm">
                Para ver tus reservas, por favor inicia sesión.
              </p>
            </div>
          </CardContent>
        </Card>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">Mis Reservas</h1>
          <Button variant="outline" onClick={fetchReservations} disabled={loading}>
            {loading ? (
              <>
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                Cargando...
              </>
            ) : (
              'Actualizar'
            )}
          </Button>
        </div>

        {error && (
          <Card>
            <CardContent className="py-6">
              <div className="flex items-center space-x-3 text-red-600">
                <AlertCircle className="w-5 h-5" />
                <p>{error}</p>
              </div>
            </CardContent>
          </Card>
        )}

        {loading ? (
          <Card>
            <CardContent className="py-12">
              <div className="text-center text-gray-600">
                <Loader2 className="w-16 h-16 mx-auto mb-4 text-gray-400 animate-spin" />
                <p className="text-lg font-medium">Cargando reservas...</p>
              </div>
            </CardContent>
          </Card>
        ) : reservations.length === 0 ? (
          <Card>
            <CardContent className="py-12">
              <div className="text-center text-gray-600">
                <Calendar className="w-16 h-16 mx-auto mb-4 text-gray-400" />
                <p className="text-lg font-medium mb-2">No tienes reservas</p>
                <p className="text-sm">
                  Cuando realices una reserva, aparecerá aquí.
                </p>
              </div>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-4">
            {reservations.map((reservation) => {
              const nights = calculateNights(reservation.start_date, reservation.end_date);
              const canCancel = String(reservation.status) === 'confirmed' || String(reservation.status) === 'pending';
              const roomPrice = Number(roomService.getRoomsPrice(reservation.room_id.toString())) || 0;
              const total_price = nights * roomPrice;
              const guestName = String(reservationService.getGuestName(reservation.id));
              const guestEmail = String(reservationService.getGuestEmail(reservation.id));
              return (
                <Card key={reservation.id}>
                  <CardContent className="p-6">
                    <div className="flex justify-between items-start">
                      <div className="space-y-3 flex-1">
                        <div className="flex items-center space-x-4">
                          <div className="text-2xl font-bold text-primary-600">
                            Habitación #{reservation.room_id}
                          </div>
                          {getStatusBadge(reservation.status)}
                        </div>

                          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 text-sm text-gray-600">
                          <div className="flex items-start space-x-2">
                            <Calendar className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <div>
                              <div className="font-medium text-gray-700">Fechas</div>
                              <div>{formatDate(reservation.start_date)}</div>
                              <div>{formatDate(reservation.end_date)}</div>
                              <div className="text-xs text-gray-500 mt-1">
                                {nights} {nights === 1 ? 'noche' : 'noches'}
                              </div>
                            </div>
                          </div>
                          <div className="flex items-start space-x-2">
                            <DollarSign className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <div>
                              <div className="font-medium text-gray-700">Precio Total</div>
                              <div className="text-lg font-bold text-primary-600">
                                {'$' + total_price.toFixed(2)}
                              </div>
                            </div>
                          </div>
                          <div className="flex items-start space-x-2">
                            <MapPin className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <div>
                              <div className="font-medium text-gray-700">Huésped</div>
                              {guestName}
                              <div className="text-xs text-gray-500"></div>
                              {guestEmail}
                            </div>
                          </div>
                          </div>
                        </div>
                        <div className="text-xs text-gray-400 mt-2">
                          ID de reserva: {reservation.id}
                        </div>
                      </div>

                      {canCancel && (
                        <Button
                          variant="danger"
                          size="sm"
                          onClick={() => handleCancelReservation(reservation.id)}
                          disabled={cancellingId === reservation.id}
                        >
                          {cancellingId === reservation.id ? (
                            <>
                              <Loader2 className="w-4 h-4 mr-1 animate-spin" />
                              Cancelando...
                            </>
                          ) : (
                            <>
                              <X className="w-4 h-4 mr-1" />
                              Cancelar
                            </>
                          )}
                        </Button>
                      )}
                  </CardContent>
                </Card>
              );
            })}
          </div>
        )}

        {!loading && reservations.length > 0 && (
          <div className="text-center text-sm text-gray-500 py-4">
            Total de reservas: {reservations.length}
          </div>
        )}
      </div>
    </Layout>
  );
}  
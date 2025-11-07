import api from '@/lib/axios';
import type {
  Reservation,
  CreateReservationRequest,
  CreateReservationResponse,
  CancelReservationRequest
} from '@/types';

export const reservationService = {
  async createReservation(reservationData: CreateReservationRequest): Promise<CreateReservationResponse> {
    const response = await api.post<CreateReservationResponse>('/api/reservations', reservationData);
    return response.data;
  },

  async getReservationById(id: string): Promise<{ reservation: Reservation }> {
    const response = await api.get<{ reservation: Reservation }>(`/api/reservations/${id}`);
    return response.data;
  },

  async cancelReservation(id: string, reason: string): Promise<void> {
    const cancelData: CancelReservationRequest = { reason };
    await api.delete(`/api/reservations/${id}`, { data: cancelData });
  },
};

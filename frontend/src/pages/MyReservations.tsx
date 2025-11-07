import React, { useState } from 'react';
import { Layout } from '@/components/layout/Layout';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Calendar, MapPin, DollarSign, X } from 'lucide-react';

export function MyReservations() {
  const [reservations] = useState<any[]>([]);

  // Note: This is a placeholder. In a real app, you would fetch reservations from the API
  // using the user's ID and display them here.

  return (
    <Layout>
      <div className="space-y-6">
        <h1 className="text-3xl font-bold text-gray-900">Mis Reservas</h1>

        <Card>
          <CardContent className="py-12">
            <div className="text-center text-gray-600">
              <Calendar className="w-16 h-16 mx-auto mb-4 text-gray-400" />
              <p className="text-lg font-medium mb-2">No tienes reservas activas</p>
              <p className="text-sm">
                Cuando realices una reserva, aparecerá aquí.
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Example of what a reservation would look like */}
        {reservations.length > 0 && (
          <div className="space-y-4">
            {reservations.map((reservation) => (
              <Card key={reservation.id}>
                <CardContent className="p-6">
                  <div className="flex justify-between items-start">
                    <div className="space-y-2">
                      <div className="flex items-center space-x-4">
                        <div className="text-2xl font-bold text-primary-600">
                          #101
                        </div>
                        <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium">
                          Activa
                        </span>
                      </div>

                      <div className="flex items-center space-x-6 text-sm text-gray-600">
                        <div className="flex items-center space-x-2">
                          <Calendar className="w-4 h-4" />
                          <span>2025-11-15 - 2025-11-20</span>
                        </div>
                        <div className="flex items-center space-x-2">
                          <MapPin className="w-4 h-4" />
                          <span>Piso 3</span>
                        </div>
                        <div className="flex items-center space-x-2">
                          <DollarSign className="w-4 h-4" />
                          <span>$500.00 total</span>
                        </div>
                      </div>
                    </div>

                    <Button variant="danger" size="sm">
                      <X className="w-4 h-4 mr-1" />
                      Cancelar
                    </Button>
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

import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8000', // AJUSTA AL PUERTO REAL DE TU API
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {

    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');

      // 🚫 YA NO redirige fuera de React
      // 🔥 Se deja que ProtectedRoute maneje la redirección
    }

    return Promise.reject(error);
  }
);

export default api;

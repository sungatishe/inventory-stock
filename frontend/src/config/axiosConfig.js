import axios from 'axios';
import Cookies from 'js-cookie';

// Создаём инстанс axios
const apiClient = axios.create({
    baseURL: 'http://127.0.0.1:8080',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Добавляем интерцептор для автоматической подстановки токена
apiClient.interceptors.request.use(
    (config) => {
        const token = Cookies.get('auth_token');
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

export default apiClient;

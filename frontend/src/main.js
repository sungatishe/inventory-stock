import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // Импорт маршрутов
import 'axios'; // для использования конфигурации axios


createApp(App)
    .use(router) // Подключение маршрутов
    .mount('#app');

import { createRouter, createWebHistory } from 'vue-router';
import AuthPage from '../components/Auth.vue';
import ItemsPage from '../components/Items.vue';
import Cookies from 'js-cookie';

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: AuthPage
    },
    {
        path: '/items',
        name: 'Items',
        component: ItemsPage,
        meta: { requiresAuth: true }, // Маршрут защищён
    },
    {
        path: '/:pathMatch(.*)*',
        redirect: '/login'
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

// Глобальный перехватчик для проверки аутентификации
router.beforeEach((to, from, next) => {
    const token = Cookies.get('auth_token');

    if (to.meta.requiresAuth && !token) {
        // Если требуется аутентификация, а токен отсутствует - редирект на login
        next({ name: 'Login' });
    } else if (to.name === 'Login' && token) {
        // Если пользователь уже залогинен и пытается попасть на login - редирект на items
        next({ name: 'Items' });
    } else {
        next();
    }
});

export default router;

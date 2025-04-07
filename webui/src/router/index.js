import { createRouter, createWebHistory } from 'vue-router';
import LoginView from '../views/LoginView.vue';
import ChatView from '../views/ChatView.vue';

const routes = [
	{ path: '/login', name: 'Login', component: LoginView },
	{ path: '/chat', name: 'Chat', component: ChatView }
];


const router = createRouter({
	history: createWebHistory(),
	routes,
});

export default router;

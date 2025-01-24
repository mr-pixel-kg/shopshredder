import {createRouter, createWebHistory} from 'vue-router'

import StartView from "./views/StartView.vue";
import LoginView from "./views/LoginView.vue";
import AdminView from "./views/AdminView.vue";
import AdminSandboxEnvironments from "./views/admin/AdminSandboxEnvironmentsView.vue";
import AdminSandboxImages from "./views/admin/AdminSandboxImagesView.vue";

const routes = [
    { path: '/', component: StartView },
    { path: '/login', component: LoginView },
    {
        path: '/admin',
        component: AdminView,
        redirect: '/admin/sandbox-environments',
        children: [
            { path: 'sandbox-environments', component: AdminSandboxEnvironments },
            { path: 'sandbox-images', component: AdminSandboxImages },
        ],
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router
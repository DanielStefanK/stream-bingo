import { createRouter, createWebHistory } from 'vue-router'

import LoginPage from './pages/LoginPage.vue'
import OAuthCallbackPage from './pages/OAuthCallbackPage.vue'
import DashboardPage from './pages/DashboardPage.vue'
import { useUserStore } from './store/userStore'

const routes = [
  { path: '/', component: DashboardPage, meta: { authType: 'loggedin' } },
  { path: '/login', component: LoginPage, meta: { authType: 'anonym' } },
  {
    path: '/oauth-callback/:provider',
    component: OAuthCallbackPage,
    meta: { authType: 'anonym' },
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  try {
    const userStore = useUserStore()
    const user = userStore.userLoad
    const authType = to.meta.authType
    if (authType === 'loggedin' && !user) {
      next('/login')
    } else if (authType === 'anonym' && user) {
      next('/')
    } else {
      next()
    }
  } catch (error) {
    console.log(error)
    next('/login')
  }
})

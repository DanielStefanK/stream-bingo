import { createRouter, createWebHistory } from 'vue-router'

import LoginPage from './pages/LoginPage.vue'
import OAuthCallbackPage from './pages/OAuthCallbackPage.vue'
import DashboardPage from './pages/DashboardPage.vue'
import { useUserStore } from './store/userStore'
import { useToast } from 'primevue'

const routes = [
  {
    path: '/',
    component: DashboardPage,
    meta: { authType: 'loggedin', title: 'Dashboard' },
  },
  { path: '/login', component: LoginPage, meta: { authType: 'anonym' } },
  {
    path: '/oauth-callback/:provider',
    component: OAuthCallbackPage,
    meta: { authType: 'anonym' },
  },
  {
    path: '/admin/user/list',
    component: () => import('./pages/UserListPage.vue'),
    meta: { authType: 'admin', title: 'List Users', subTitle: 'Admin Console' },
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, _, next) => {
  try {
    const toast = useToast()

    const userStore = useUserStore()

    let user = userStore.user

    if (userStore.state !== 'loaded') {
      user = await userStore.userLoad
    }
    const authType = to.meta.authType
    if (authType === 'loggedin' && !user) {
      console.log('not logged in')
      userStore.user = null
      userStore.state = 'loaded'
      next('/login')
    } else if (authType === 'anonym' && user) {
      console.log('already logged in')
      next('/')
    } else if (authType === 'admin' && (!user || !user.admin)) {
      console.log('not admin')
      toast.add({
        severity: 'error',
        summary: 'Unauthorized',
        detail: 'You are not authorized to access this page.',
        life: 5000,
      })
      next('/')
    } else {
      console.log('proceed')
      next()
    }
  } catch (error) {
    console.log('error in router')
    console.log(error)
    next('/login')
  }
})

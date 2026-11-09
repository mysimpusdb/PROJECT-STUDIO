import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import AppLayout from '../components/layout/AppLayout.vue'
import Home from '../views/Home.vue'
import Projects from '../views/Projects.vue'
import ProjectDetail from '../views/ProjectDetail.vue'
import Login from '../views/Login.vue'
import { CheckAuth } from '../../wailsjs/go/services/AuthService'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    component: AppLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: Home
      },
      {
        path: 'projects',
        name: 'Projects',
        component: Projects
      },
      {
        path: 'projects/:id',
        name: 'ProjectDetail',
        component: ProjectDetail
      }
      // Add other routes here (projects, tasks, etc)
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach(async (to, from) => {
  // Try to check authentication state via Wails binding
  let isLoggedIn = false
  try {
    const authState = await CheckAuth()
    isLoggedIn = authState.is_logged_in
  } catch (err) {
    console.error("Auth check failed:", err)
  }

  if (to.meta.requiresAuth && !isLoggedIn) {
    // If route requires auth and not logged in, redirect to login
    return '/login'
  } else if (to.path === '/login' && isLoggedIn) {
    // If logged in and trying to access login page, redirect to home
    return '/'
  }
  
  return true
})

export default router

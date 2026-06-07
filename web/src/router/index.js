import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    redirect: '/polls',
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/classes',
    name: 'ClassList',
    component: () => import('@/views/class/ClassListView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/classes/:id',
    name: 'ClassDetail',
    component: () => import('@/views/class/ClassDetailView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/timetable',
    name: 'Timetable',
    component: () => import('@/views/timetable/TimetableView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/timetable/edit',
    name: 'TimetableEdit',
    component: () => import('@/views/timetable/TimetableEditView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/polls',
    name: 'PollList',
    component: () => import('@/views/poll/PollListView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/polls/create',
    name: 'PollCreate',
    component: () => import('@/views/poll/PollCreateView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/polls/:id',
    name: 'PollDetail',
    component: () => import('@/views/poll/PollDetailView.vue'),
    meta: { requiresAuth: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫：检查登录状态
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (!to.meta.requiresAuth && token && to.name === 'Login') {
    // 已登录用户访问登录页，重定向到首页
    next({ name: 'Home' })
  } else {
    next()
  }
})

export default router

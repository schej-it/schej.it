import Vue from 'vue'
import VueRouter from 'vue-router'
import Landing from '@/views/Landing'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'landing',
    component: Landing,
  },
  {
    path: '/home',
    name: 'home',
    component: () => import('@/views/Home.vue')
  },
  {
    path: '/e/:eventId',
    name: 'event',
    component: () => import('@/views/Event.vue'),
    props: true,
  },
  {
    path: '/auth',
    name: 'auth',
    component: () => import('@/views/Auth.vue')
  },


  {
    path: '/privacy-policy',
    name: 'privacy-policy',
    component: () => import('@/views/PrivacyPolicy.vue'),
  },
  {
    path: '/donut',
    name: 'donut',
    redirect: '/e/6417b47a0c6fc139f870a47d',
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router

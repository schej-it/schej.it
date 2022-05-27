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
    path: '/schedule',
    name: 'schedule',
    component: () => import('@/views/Schedule.vue')
  },
  {
    path: '/friends',
    name: 'friends',
    component: () => import('@/views/Friends.vue')
  },
  {
    path: '/friend-schedule',
    name: 'friend-schedule',
    component: () => import('@/views/FriendSchedule.vue'),
  },
  {
    path: '/sign-in',
    name: 'sign-in',
    component: () => import('@/views/SignIn.vue')
  },
  {
    path: '/j/:eventId',
    name: 'join',
    component: () => import('@/views/Join.vue'),
    props: true,
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
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router

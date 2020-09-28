import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home'
import store from '../store'

Vue.use(VueRouter)

const routes = [
  {
    path: '/auth',
    name: 'Auth',
    component: () => import('../views/Auth'),
    meta: { tree: 0 }
  },
  {
    path: '/',
    redirect: '/data',
    component: Home,
    children: [
      {
        path: 'data',
        name: 'Data',
        component: () => import('../views/Manager/Data')
      },
      {
        path: 'user',
        name: 'Growth',
        component: () => import('../views/Manager/User')
      },
      {
        path: 'messages',
        name: 'Message',
        component: () => import('../views/Manager/Message')
      },
      {
        path: 'setting',
        name: 'Setting',
        component: () => import('../views/Manager/Setting')
      }
    ]
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router

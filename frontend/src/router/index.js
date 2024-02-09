import { createRouter, createWebHistory } from 'vue-router'
import axios from 'axios'
import HomeView from '../views/HomeView.vue'
// import DirectionsView from '../views/DirectionsView.vue'
import LocationView from '../views/LocationView.vue'
import LoginView from '../views/LoginView.vue'
import MapView from '../views/MapView.vue'

const routes = [
  {
    path: '/',
    name: 'login',
    component: LoginView
  },
  {
    path: '/home',
    name: 'home',
    component: HomeView
  },
  // {
  //   path: '/directions',
  //   name: 'directions',
  //   component: DirectionsView
  // },
  {
    path: '/map',
    name: 'map',
    component: MapView
  },
  {
    path: '/location',
    name: 'location',
    component: LocationView
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

// eslint-disable-next-line
router.beforeEach((to, from) => {
  if (to.name === 'login') {
    return true
  }

  if (!localStorage.getItem('token')) {
    return {
      name: 'login'
    }
  }

  checkTokenAuthenticity()
})

const checkTokenAuthenticity = () => {
  // TODO: change this .env
  axios.get('http://localhost:6969/api/auth', {
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`
    }
  })
    // eslint-disable-next-line
    .then((response) => {})
    // eslint-disable-next-line
    .catch((error) => {
      localStorage.removeItem('token')
      router.push({
        name: 'login'
      })
    })
}

export default router

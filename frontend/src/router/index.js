import { createRouter, createWebHistory } from 'vue-router'
import axios from 'axios'
import CusHomeView from '../views/customer/CusHomeView.vue'
// import DirectionsView from '../views/DirectionsView.vue'
import CusLocationView from '../views/customer/CusLocationView.vue'
import CusLoginView from '../views/customer/CusLoginView.vue'
import CusMapView from '../views/customer/CusMapView.vue'
import ChooseAppView from '../views/ChooseAppView.vue'
import DriverLoginView from '../views/driver/DriverLoginView.vue'

const routes = [
  {
    path: '/',
    name: 'choose-app',
    component: ChooseAppView
  },
  {
    path: '/driver-login',
    name: 'driver-login',
    component: DriverLoginView
  },
  {
    path: '/cus-login',
    name: 'cus-login',
    component: CusLoginView
  },
  {
    path: '/cus-home',
    name: 'cus-home',
    component: CusHomeView
  },
  // {
  //   path: '/directions',
  //   name: 'directions',
  //   component: DirectionsView
  // },
  {
    path: '/cus-map',
    name: 'cus-map',
    component: CusMapView
  },
  {
    path: '/cus-location',
    name: 'cus-location',
    component: CusLocationView
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

// eslint-disable-next-line
router.beforeEach((to, from) => {
  if (from.name === 'choose-app' || to.name === 'cus-login' || to.name === 'driver-login') {
    return true
  }

  if (from.name === 'cus-login') {
    if (!localStorage.getItem('cus-token')) {
      return {
        name: 'cus-login'
      }
    }
  
    checkTokenAuthenticity()
  }

  if (from.name === 'driver-login') {
    if (!localStorage.getItem('driver-token')) {
      return {
        name: 'driver-login'
      }
    }
  
    checkDriverTokenAuthenticity()
  }
})

const checkTokenAuthenticity = () => {
  // TODO: change this .env
  axios.get('http://localhost:6969/api/auth', {
    headers: {
      Authorization: `Bearer ${localStorage.getItem('cus-token')}`
    }
  })
    // eslint-disable-next-line
    .then((response) => {})
    // eslint-disable-next-line
    .catch((error) => {
      localStorage.removeItem('cus-token')
      router.push({
        name: 'cus-login'
      })
    })
}

const checkDriverTokenAuthenticity = () => {
  // TODO: change this .env
  axios.get('http://localhost:6969/api/driver/auth', {
    headers: {
      Authorization: `Bearer ${localStorage.getItem('driver-token')}`
    }
  })
    // eslint-disable-next-line
    .then((response) => {})
    // eslint-disable-next-line
    .catch((error) => {
      localStorage.removeItem('driver-token')
      router.push({
        name: 'driver-login'
      })
    })
}

export default router

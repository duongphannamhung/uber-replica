import { createRouter, createWebHistory } from 'vue-router'
import axios from 'axios'
import CusHomeView from '../views/customer/CusHomeView.vue'
// import DirectionsView from '../views/DirectionsView.vue'
import CusLocationView from '../views/customer/CusLocationView.vue'
import CusLoginView from '../views/customer/CusLoginView.vue'
import CusMapView from '../views/customer/CusMapView.vue'
import ChooseAppView from '../views/ChooseAppView.vue'
import DriverLoginView from '../views/driver/DriverLoginView.vue'
import DriverHomeView from '../views/driver/DriverHomeView.vue'
import CusFindingDriver from '../views/customer/CusFindingDriver.vue'
import DriverDriveToCus from '../views/driver/DriverDriveToCus.vue'
import CusWaitingDriverArrive from '../views/customer/CusWaitingDriverArrive.vue'

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
    path: '/cus-finding-driver',
    name: 'cus-finding-driver',
    component: CusFindingDriver
  },
  {
    path: '/cus-waiting-driver-arrive',
    name: 'cus-waiting-driver-arrive',
    component: CusWaitingDriverArrive
  },
  {
    path: '/cus-message-in-trip',
    name: 'cus-message-in-trip',
    component: () => import("../views/customer/CustomerMessageInTrip.vue"),
  }, 
  {
    path: "/driver-message-in-trip",
    name: "driver-message-in-trip",
    component: () => import("../views/driver/DriverMessageInTrip.vue"),
  },
  {
    path: '/driver-drive-to-cus',
    name: 'driver-drive-to-cus',
    component: DriverDriveToCus
  },
  {
    path: '/driver-home',
    name: 'driver-home',
    component: DriverHomeView
  },
  {
    path: '/cus-in-trip',
    name: 'cus-in-trip',
    component: () => import("../views/customer/CusInTrip.vue"), 
  },
  {
    path: '/driver-in-trip',
    name: 'driver-in-trip',
    component: () => import("../views/driver/DriverInTrip.vue"),
  },
  {
    path: '/driver-complete-trip',
    name: 'driver-complete-trip',
    component: () => import("../views/driver/DriverCompleteTrip.vue"),
  },
  {
    path: '/cus-complete-trip',
    name: 'cus-complete-trip',
    component: () => import("../views/customer/CusCompleteTrip.vue"),
  },
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
    .then((response) => {
    })
      
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

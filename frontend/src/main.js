import App from './App.vue'
import './registerServiceWorker'
import router from './router'
// vue add tailwind
import './assets/tailwind.css'
// https://www.npmjs.com/package/vue-material-design-icons

import axios from 'axios'
axios.defaults.baseURL = process.env.VUE_APP_API_URL

import { createApp, provide, h } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

import VueGoogleMaps from '@fawmi/vue-google-maps'
import { websocketStore } from '@/store/websocket-store'
// import DriverToggle from '@/components/DriverToggle.vue'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

const app = createApp({
    setup() {
        provide('websocketStore', websocketStore)
    },
    render() {
        return h(App);
    }
})

app.use(router)
app.use(pinia)

app.use(VueGoogleMaps, {
    load: {
        key: process.env.VUE_APP_GOOGLE_MAP_API_KEY,
        libraries: 'places'
    }
})

app.mount('#app')
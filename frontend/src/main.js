import { createApp } from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
// vue add tailwind
import './assets/tailwind.css'
// https://www.npmjs.com/package/vue-material-design-icons

import axios from 'axios'
axios.defaults.baseURL = process.env.VUE_APP_API_URL

import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

import VueGoogleMaps from '@fawmi/vue-google-maps'
// import DriverToggle from '@/components/DriverToggle.vue'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

const app = createApp(App)
// const app = createApp(App, {
//     data() {
//         return {
//             status: true,
//             data: [],
//             interval:null
//         };
//     },
//     methods: {
//         callEverySecond() {
//             console.log("called every 5 seconds")
//             // axios.get("https://randomuser.me/api/?results=5").then(
//             // (response) =>
//             //     // console.log(response)
//             //     (this.data = response.data.results)
//             // )
//           }
//     },
//     mounted() {
//         this.callEverySecond()
//     },
//     created() {
//         this.interval = setInterval(() => {
//             this.callEverySecond()
//         }, 3000)
//     },
//     beforeUnmount() {
//         clearInterval(this.intervalId)
//     }
// })

app.use(router)
app.use(pinia)

app.use(VueGoogleMaps, {
    load: {
        key: process.env.VUE_APP_GOOGLE_MAP_API_KEY,
        libraries: 'places'
    }
})

app.mount('#app')
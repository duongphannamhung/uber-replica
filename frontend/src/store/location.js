import { reactive } from 'vue'
import { defineStore } from 'pinia'

const getUserLocation = async () => {
    return new Promise((res, rej) => {
        navigator.geolocation.getCurrentPosition(res, rej)
    })
}

export const useLocationStore = defineStore('location', () => {
    const destination = reactive({
        name: '',
        address: '',
        geometry: {
            lat: null,
            lng: null
        }
    })

    // const current = reactive({
    //     geometry: {
    //         lat: null,
    //         lng: null
    //     }
    // })

    const departure = reactive({
        display_name: '',
        geometry: {
            lat: null,
            lng: null
        }
    })

    const updateDestination = async () => {
        if (destination.geometry.lat && destination.geometry.lng && destination.address) return
        destination.geometry = {
            lat: localStorage.getItem('destination_lat'),
            lng: localStorage.getItem('destination_lng')
        }
        destination.address = localStorage.getItem('destination_address')
        destination.name = localStorage.getItem('destination_name')
    }

    const updateCurrentLocation = async () => {
        if (departure.geometry.lat && departure.geometry.lng && departure.display_name) return

        if (localStorage.getItem('departure_lat') && localStorage.getItem('departure_lng') && localStorage.getItem('departure_address') && localStorage.getItem('departure_name')) {
            departure.geometry = {
                lat: localStorage.getItem('departure_lat'),
                lng: localStorage.getItem('departure_lng')
            }
            departure.display_name = localStorage.getItem('departure_address')
            departure.name = localStorage.getItem('departure_name')
            return
        }

        const userLocation = await getUserLocation()
        departure.geometry = {
            lat: userLocation.coords.latitude,
            lng: userLocation.coords.longitude
        }
        const response = await fetch(`https://nominatim.openstreetmap.org/reverse?lat=${departure.geometry.lat}&lon=${departure.geometry.lng}&format=json`, {
            headers: {
              'User-agent': navigator.userAgent,
            }
          });
        response.json().then(data => {
            departure.display_name = data.display_name
        })

        localStorage.setItem('departure_lat', departure.geometry.lat)
        localStorage.setItem('departure_lng', departure.geometry.lng)
        localStorage.setItem('departure_address', departure.display_name)
        localStorage.setItem('departure_name', departure.name)
    }

    const reset = () => {
        destination.name = ''
        destination.address = ''
        destination.geometry.lat = null
        destination.geometry.lng = null

        departure.display_name = ''
        departure.geometry.lat = null
        departure.geometry.lng = null
    }
    
    return { destination, departure, updateCurrentLocation, reset, updateDestination }
})

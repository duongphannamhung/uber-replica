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

    const updateCurrentLocation = async () => {
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
    
    return { destination, departure, updateCurrentLocation, reset }
})

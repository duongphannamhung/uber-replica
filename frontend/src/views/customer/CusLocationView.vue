<template>
    <div class="pt-16">
        <h1 class="text-3xl font-semibold mb-4 text-center">Bạn muốn đi đâu?</h1>
        <form action="#" @submit.prevent="">
            <div class="overflow-hidden shadow sm:rounded-md max-w-sm mx-auto text-left">
                <div class="bg-white px-2 py-5 sm:p-6">
                <div class="flex items-center space-x-2">
                    <img src="img/logo/start_point.png" alt="Location Icon" class="w-6 h-6">
                    <div v-if="!showLocationAutocomplete" @click="showLocationAutocomplete = true" class="mt-1 bg-gray-200 block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-black focus:outline-none">Địa điểm hiện tại</div>
                    <GMapAutocomplete
                    v-else
                    placeholder="Nơi đi"
                    @place_changed="handleDepartureChanged"
                    class="mt-1 block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-black focus:outline-none">
                    </GMapAutocomplete>
                </div>
                <div class="flex items-center space-x-2">
                    <img src="img/logo/destination_point.png" alt="Destination Icon" class="w-6 h-6">
                    <GMapAutocomplete
                    placeholder="Điểm đến"
                    @place_changed="handleLocationChanged"
                    class="mt-1 block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-black focus:outline-none">
                    </GMapAutocomplete>
                </div>
                </div>
                <div class="bg-gray-50 px-4 py-3 text-right sm:px-6">
                    <button
                        @click.prevent="handleSelectLocation"
                        type="button"
                        class="inline-flex justify-center rounded-md border border-transparent bg-black py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-gray-600 focus:outline-none">
                        Find A Ride
                    </button>
                </div>
            </div>
        </form>
    </div>
</template>
<script setup>
import { useLocationStore } from '@/store/location'
import { useRouter } from 'vue-router'
import { ref } from 'vue'

const location = useLocationStore()
const router = useRouter()

const showLocationAutocomplete = ref(false)

const handleDepartureChanged = async (e) => {
    console.log('handleDepartureChanged', e)
    location.$patch({
        departure: {
            geometry: {
                lat: e.geometry.location.lat(),
                lng: e.geometry.location.lng()
            }
        }
    })
    const response = await fetch(`https://nominatim.openstreetmap.org/reverse?lat=${e.geometry.location.lat()}&lon=${e.geometry.location.lng()}&format=json`, {
            headers: {
              'User-agent': navigator.userAgent,
            }
          });
        response.json().then(data => {
            location.$patch({
                departure: {
                    display_name: data.display_name
                }
            })
        })
    }

const handleLocationChanged = (e) => {
    // console.log('handleLocationChanged', e)
    location.$patch({
        destination: {
            name: e.name,
            address: e.formatted_address,
            geometry: {
                lat: e.geometry.location.lat(),
                lng: e.geometry.location.lng()
            }
        }
    })
}

const handleSelectLocation = () => {
    if (location.destination.name !== '') {
        router.push({
            name: 'cus-map'
        })
    }
}
</script>
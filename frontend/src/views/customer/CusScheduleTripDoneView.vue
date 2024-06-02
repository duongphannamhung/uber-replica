<template>
    <div class="pt-16">
        <h1 class="text-3xl font-semibold mb-4 text-center">Chuyến đi đã lên lịch</h1>
        <div 
            id="BackBtn" 
            class="absolute z-50 rounded-full bg-white p-1 top-8 left-4"
            @click="goBack()"
            >
            <ArrowLeftIcon :size="40" />
        </div>
        <form action="#" @submit.prevent="">
            <div class="overflow-hidden shadow sm:rounded-md max-w-md mx-auto text-left">
                <div class="bg-white px-2 py-5 sm:p-6">
                    <div class="flex items-center">
                        <img src="img/logo/start_point.png" alt="Start Icon" class="mr-2 w-6 h-6">
                        <p class="font-bold text-lg">{{ beautifulizeAddress(location.departure.display_name) }}</p>
                        </div>
                        <p class="text-lg ml-8">{{ tailAddress(location.departure.display_name) }}</p>

                    <div class="border-b border-gray-400 my-2"></div>
                    <div class="flex items-center">
                        <img src="img/logo/destination_point.png" alt="Destination Icon" class="mr-2 w-6 h-6">
                        <p class="font-bold text-lg">{{ beautifulizeAddress(location.destination.address) }}</p>
                        </div>
                        <p class="text-lg ml-8">{{ tailAddress(location.destination.address) }}</p>
                <div>
                    <div class="border-b border-gray-400 my-2"></div>
                <div class="flex items-center space-x-2">
                    <img src="img/logo/calendar.png" alt="Calendar Icon" class="w-6 h-6">
                    <p class="font-bold text-lg">{{ schedule_time }}</p>
                </div>
                <div class="flex items-center space-x-2 mt-2">
                    <img src="img/logo/bus.png" alt="Transportation Icon" class="w-6 h-6">
                    <p class="font-bold text-lg">{{ vehicle_name }}</p>
                </div>
                </div>
                </div>
                <div class="bg-gray-50 px-4 py-3 sm:px-6 flex items-center">
                <p class="font-bold text-xl text-blue-600 mr-2">Thời gian còn lại:</p>
                <div v-if="timeRemaining > 0" class="text-xl font-bold text-blue-600">
                    <span>{{ Math.floor(timeRemaining / 60) }}:</span>
                    <span>{{ ('0' + (timeRemaining % 60)).slice(-2) }}</span>
                </div>
                <button
                    @click.prevent="handleCancelSchedule"
                    type="button"
                    class="inline-flex justify-center rounded-md border border-transparent bg-red-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-gray-600 focus:outline-none ml-auto">
                    Hủy lịch
                </button>
                </div>
            </div>
        </form>
    </div>
</template>
<script setup>
import { useLocationStore } from '@/store/location'
import { useRouter } from 'vue-router'
import ArrowLeftIcon from 'vue-material-design-icons/ArrowLeft.vue';
import { onMounted, ref } from 'vue'
import axios from 'axios'

const location = useLocationStore()
const router = useRouter()
const schedule_time = ref(null)
const vehicle_name = ref(null)
const timeRemaining = ref(null)
const distance = ref({text: '', value: null})
const duration = ref({text: '', value: null})

const goBack = () => {
    router.push({
      name : 'cus-home'
    })
  }

const calculateTimeRemaining = () => {
    const currentDateTime = new Date();
    const selectedDateTime = new Date(localStorage.getItem('scheduleDateTime'));
    const diff = selectedDateTime - currentDateTime;
    return Math.floor(diff / 1000);
}

const startCountdown = () => {
      const intervalId = setInterval( async () => {
        if (timeRemaining.value === null) {
            timeRemaining.value = calculateTimeRemaining()
        }
        
        if (timeRemaining.value > 0) {
            timeRemaining.value -= 1;
        } else {
          clearInterval(intervalId); 

          let vehicleName = localStorage.getItem('vehicleName')
          let fare = ((distance.value.value * 9000) * localStorage.getItem('multiplier') / 1000000).toFixed(0) * 1000;
          let tripRequest = {
            user_id: localStorage.getItem('current_user_id'),
            user_phone: localStorage.getItem('current_user_phone'),
            vehicle: vehicleNameToType(vehicleName),
            // eslint-disable-next-line
            departure_point : new google.maps.LatLng(location.departure.geometry),
            departure_name: location.departure.display_name,
            // eslint-disable-next-line
            destination_point : new google.maps.LatLng(location.destination.geometry),
            destination_name: location.destination.address,
            fare: fare
        }

        // if (vehicleName == "UrepBike") {
        await axios.post('create-trip', tripRequest, {
            headers: {
                Authorization: `Bearer ${localStorage.getItem('cus-token')}`
            }
        })
            .then(async (response) => {
                await getDistance()

                localStorage.setItem('current_trip_id', response.data.trip_id)

                // let fare = 

                console.log(`fare ${fare}, distance: ${distance.value.value}, multiplier: ${localStorage.getItem('multiplier')}`)

                // await axios.post('http://localhost:6969/api/driver/update-trip-fare', {
                //     trip_id: response.data.trip_id,
                //     fare: fare
                // }, {
                //     headers: {
                //     Authorization: `Bearer ${localStorage.getItem('cus-token')}`
                //     }
                // })

                localStorage.removeItem('scheduleDateTime')
                localStorage.removeItem('vehicleName')
                router.push({
                    name : 'cus-finding-driver'
                })
            })
            .catch((error) => {
                console.error(error)
                alert(error.response.data.message)
                router.push({
                    name : 'cus-home'
                })
            })
        }
      }, 1000);
}

const vehicleNameToType = (vehicleName) => {
    if (vehicleName == 'UrepBike') {
        return 1
    } else if (vehicleName == 'UrepCar') {
        return 2
    } else if (vehicleName == 'UrepCar4') {
        return 3
    } else if (vehicleName == 'UrepCar7') {
        return 4
    } else {
        return -1
    }
}

const getDistance = async () => {
    // eslint-disable-next-line
    let departurePoint = new google.maps.LatLng(location.departure.geometry)
    // eslint-disable-next-line
    let destinationPoint = new google.maps.LatLng(location.destination.geometry)
    let res = await axios.get('distance/' + departurePoint + '/' + destinationPoint, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('cus-token')}`
      }
    })

    distance.value.text = res.data.distance_text
    distance.value.value = res.data.distance
    duration.value.text = res.data.duration_text
    duration.value.value = res.data.duration
  }

onMounted(async () => {
    await location.updateDestination()
    await location.updateCurrentLocation()

    schedule_time.value = localStorage.getItem('scheduleDateTime')
    vehicle_name.value = localStorage.getItem('vehicleName')
    startCountdown()

    await sleep(3000)
    if (timeRemaining.value <= 0) {
        localStorage.removeItem('scheduleDateTime')
        localStorage.removeItem('vehicleName')
        router.push({
            name : 'cus-schedule-trip'
        })
    }
})

var sleepSetTimeout_ctrl;
function sleep(ms) {
        clearInterval(sleepSetTimeout_ctrl);
        return new Promise(resolve => sleepSetTimeout_ctrl = setTimeout(resolve, ms));
    }

const beautifulizeAddress = (address) => {
        if (!address) return ''
        let list_address = address.split(',')
        if (list_address.length < 2) {
            return address
        }
        else {
            return list_address[0] + ', ' + list_address[1] + ', ' + list_address[2]
        }
    }

    const tailAddress = (address) => {
        if (!address) return ''
        let list_address = address.split(',')
        if (list_address.length < 2) {
            return address
        }
        else {
            if (Number.isInteger(Number(list_address[list_address.length - 2]))) {
                return list_address[list_address.length - 3] + ', ' + list_address[list_address.length - 1]
            }
            else {
                return list_address[list_address.length - 2] + ', ' + list_address[list_address.length - 1]
            }
        }
    }

const handleCancelSchedule = () => {
    localStorage.removeItem('scheduleDateTime')
    localStorage.removeItem('vehicleName')
    router.push({
        name : 'cus-schedule-trip'
    })
}

</script>
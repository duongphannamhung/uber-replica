<template>
    <div>
      <GMapMap
          :zoom="4" :center="location.departure.geometry"
          :options="{
            minZoom: 3,
            maxZoom : 17,
            fullscreenControl : false,
            zoomControl : false,
            streetViewControl : false,
            mapTypeControl : false
          }"
          ref="gMap"
          style="width: 100%; height: 65vh;">
      </GMapMap>

      <div id="CustomerInfo" class="w-full">
      <div class="w-full h-2 border-t"></div>
      <div class="w-full text-center border-t-2 p-1.5 text-gray-700 text-lg font-semibold">
        Di chuyển đến khách hàng
      </div>
      <div class="border-b"></div>
      <div class="flex justify-between items-center mt-3">
        <div>
          <p class="font-bold text-lg ml-5">{{ beautifulizeAddress(customer_info.address) }}</p>
          <p class="ml-5">{{ customer_info.name }}</p>
          <div style="background-color: lightgoldenrodyellow; padding: 4px; width: 175px; margin-left: 20px;">
            <p style="font-weight: 600;">Tiền mặt: {{ convertPriceToVND(customer_info.fare) }}</p>
        </div>
        </div>
        <div>
          <img :src="customer_info.image" alt="Customer Image" class="rounded-full w-20 h-20 mr-10">
        </div>
      </div>
    </div>  

      <div 
          class="
            flex 
            items-center 
            justify-center 
            bg-white
            py-6 
            px-4
            w-full 
            absolute 
            bottom-0 
            shadow-inner
          " 
        >
      <div class="flex justify-around w-full">
            <button  
                class="
                border-4 border-black
                bg-white
                text-2xl 
                text-black
                py-2 
                px-4 
                rounded-sm
                w-1/2
                mr-2
                " @click="goToMessage()"
            >
                Tin Nhắn
            </button>

            <button 
                v-if="driverArrived"
                class="
                bg-black 
                text-2xl 
                text-white
                py-2
                px-4 
                rounded-sm
                w-1/2
                ml-2
                "
                @click="startTrip()"
            >
                Bắt Đầu Chuyến
            </button>
            <button 
                v-else
                class="
                bg-gray-400
                text-2xl 
                text-white
                py-2
                px-4 
                rounded-sm
                w-1/2
                ml-2
                "
            >
                Bắt Đầu Chuyến
            </button>
          </div>
      </div>
    </div>
  </template>
  
  <script setup>
    import axios from 'axios';
    import { onMounted, ref, reactive } from 'vue'
    import { useRouter } from 'vue-router'
    // import ArrowLeftIcon from 'vue-material-design-icons/ArrowLeft.vue';
    // import { useDirectionStore } from '@/store/direction-store';
    import { useLocationStore } from '@/store/location'
    // import { useTripStore } from '@/stores/trip'

    const driverArrived = ref(false);
    const router = useRouter()
    const location = useLocationStore()
    const customer_location = reactive({
        display_name: '',
        geometry: {
            lat: null,
            lng: null
        }
    })
    // const trip = useTripStore()
  
    const gMap = ref(null)
    let driver_come_interval = null

    const customer_info = ref({
      name: 'Dương Phan Nam Hưng',
      address: '', // TODO: get real info
      image: 'https://www.w3schools.com/howto/img_avatar.png',
      fare: 0,
    })

    const convertPriceToVND = (price) => {
      let k = 0
      price = price.toString()
      let result = ''
      for (let i = price.length - 1; i >= 0; i--) {
        if (k % 3 == 0 && k != 0) {
          result = price[i] + '.' + result
        } else {
          result = price[i] + result
        }
        k++
      }
      return result + ' đ'
    }

    // const direction = useDirectionStore()
  
    onMounted(async () => {
        await axios.get(`http://localhost:6969/api/trip/${localStorage.getItem('current_trip_id')}`, {
            headers: {
                Authorization: `Bearer ${localStorage.getItem('driver-token')}`
            }
        }).then((resp) => {
            customer_info.value.fare = resp.data.fare
            customer_location.geometry.lat = resp.data.departure_latitude
            customer_location.geometry.lng = resp.data.departure_longitude
            customer_info.value.address = resp.data.departure_name
        }).catch((error) => {
            console.error(error)
            alert(error.response.data.message)
        })

        driver_come_interval = setInterval(updateDriverLocation, 1000);
        
        // lets get the driver current location
        await location.updateDestination()
        await location.updateCurrentLocation()
  
        while (!gMap.value) {
          await sleep(1000);
        }

        // draw a path on the map
        gMap.value.$mapPromise.then((mapObject) => {
            // eslint-disable-next-line
            let departurePoint = new google.maps.LatLng(location.departure.geometry),
            // eslint-disable-next-line
                destinationPoint = new google.maps.LatLng(customer_location.geometry),
            // eslint-disable-next-line
                directionsService = new google.maps.DirectionsService,
            // eslint-disable-next-line
                directionsDisplay = new google.maps.DirectionsRenderer({
                    map: mapObject
                })
  
            directionsService.route({
                origin: departurePoint,
                destination: destinationPoint,
                avoidTolls: false,
                avoidHighways: false,
            // eslint-disable-next-line
                travelMode: google.maps.TravelMode.DRIVING
            }, (res, status) => {
            // eslint-disable-next-line
                if (status === google.maps.DirectionsStatus.OK) {
                    directionsDisplay.setDirections(res)
                } else {
                    console.error(status)
                }
            })
        })
    })

    // const latLng = ref({ start: { lat: null, lng: null }, end: { lat: null, lng: null } })
  
    // onMounted(async () => {
    //   if (!direction.pickup || !direction.destination) { router.push('/') }
    //   setTimeout(() => { initMap() }, 50)
    // })
  
    const goToMessage = async () => {
        router.push({
            name : 'driver-message-in-trip'
        })
    }

    const updateDriverLocation = async () => {
    if (getDistanceFromLatLonInKm(location.departure.geometry.lat, location.departure.geometry.lng, customer_location.geometry.lat, customer_location.geometry.lng) < 100) {
      clearInterval(driver_come_interval)

      // await sleep(7000);
      driverArrived.value = true;
    }
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

  const startTrip = async () => {
    localStorage.setItem('driver_arrived', 'true')
    router.push({
        name : 'driver-in-trip'
    })
  }

  var sleepSetTimeout_ctrl;

  function sleep(ms) {
      clearInterval(sleepSetTimeout_ctrl);
      return new Promise(resolve => sleepSetTimeout_ctrl = setTimeout(resolve, ms));
  }

  function getDistanceFromLatLonInKm(lat1,lon1,lat2,lon2) {
    var R = 6371; // Radius of the earth in km
    var dLat = deg2rad(lat2-lat1);  // deg2rad below
    var dLon = deg2rad(lon2-lon1); 
    var a = 
      Math.sin(dLat/2) * Math.sin(dLat/2) +
      Math.cos(deg2rad(lat1)) * Math.cos(deg2rad(lat2)) * 
      Math.sin(dLon/2) * Math.sin(dLon/2)
      ; 
    var c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a)); 
    var d = R * c; // Distance in km
    return d * 1000;
  }

  function deg2rad(deg) {
    return deg * (Math.PI/180)
  }

  </script>
  
  <style lang="scss">
  
    #map {
      width:100%;
      height: 45vh;
      top: 0px;
      left: 0px;
    }
  
    .gm-style-cc {
      display: none;
    }
  
    #VehicleSelection {
      .scrollSection {
        height: calc(56.6vh - 100px); 
        position: absolute; 
        overflow-y: auto; 
        width: 100%
      }
  
      .bg-custom-gray {
        background-color: rgb(237, 237, 237);
      }
    }
  </style>
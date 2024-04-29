<template>
    <div>
  
      <!-- <div 
        id="BackBtn" 
        class="absolute z-50 rounded-full bg-white p-1 top-8 left-4"
        @click="goBack()"
      >
        <ArrowLeftIcon :size="40" />
      </div> -->
      
      <!-- <div id="map" > -->
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
            Di chuyển đến điểm đến
        </div>
        <div class="border-b"></div>
        <div class="flex justify-between items-center mt-3">
          <div>
            <p class="font-bold text-lg ml-5">{{ beautifulizeAddress(destination.address) }}</p>
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
                Nhắn tin
            </button>
            <button  
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
                @click="completeTrip()"
            >
                Hoàn thành chuyến
            </button>
            </div>
        </div>
    </div>
  </template>
  
  <script setup>
    import axios from 'axios';
    import { onMounted, ref, reactive, inject } from 'vue'
    import { useRouter } from 'vue-router'
    import { useLocationStore } from '@/store/location'
  
    // import mapStyles from '../mapStyles'
    const router = useRouter()
    const location = useLocationStore()
    const destination = ref({address: '', geometry: { lat: null, lng: null}})
    // const trip = useTripStore()
    const conn = inject('websocketStore').conn
  
    const gMap = ref(null)
    const customer_location = reactive({
        // name: '',
        // address: '',
        geometry: {
            lat: null,
            lng: null
        }
    })

    var sleepSetTimeout_ctrl;

    const customer_info = ref({
      name: 'Dương Phan Nam Hưng',
      address: '',
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

    function sleep(ms) {
        clearInterval(sleepSetTimeout_ctrl);
        return new Promise(resolve => sleepSetTimeout_ctrl = setTimeout(resolve, ms));
    }
    // const direction = useDirectionStore()
  
    onMounted(async () => {
        await location.updateCurrentLocation()
        
        await axios.get(`http://localhost:6969/api/trip/${localStorage.getItem('current_trip_id')}`, {
            headers: {
                Authorization: `Bearer ${localStorage.getItem('driver-token')}`
            }
        }).then((resp) => {
            customer_location.geometry.lat = resp.data.origin_latitude
            customer_location.geometry.lng = resp.data.origin_longitude
            destination.value.address = resp.data.destination_name
            destination.value.geometry.lat = resp.data.destination_latitude
            destination.value.geometry.lng = resp.data.destination_longitude
            customer_info.value.address = resp.data.departure_name
            customer_info.value.fare = resp.data.fare
        }).catch((error) => {
            console.error(error)
            alert(error.response.data.message)
        }) 
  
        while (!gMap.value) {
          await sleep(1000);
        }

        // draw a path on the map
        gMap.value.$mapPromise.then((mapObject) => {
            // eslint-disable-next-line
            let departurePoint = new google.maps.LatLng(location.departure.geometry),
            // eslint-disable-next-line
                destinationPoint = new google.maps.LatLng(destination.value.geometry),
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

    const goToMessage = async () => {
        router.push({
            name : 'driver-message-in-trip'
        })
    }

    const completeTrip = async () => {
        if (conn === null) {
          alert("Ws conn null. Can't completed trip")
          return
        }
        
        conn.send('finishline-@123!(*234kh219871233hadsfh')
        router.push({
            name : 'driver-complete-trip'
        })
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
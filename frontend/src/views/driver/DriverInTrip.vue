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
          :zoom="4" :center="location.current.geometry"
          :options="{
            minZoom: 3,
            maxZoom : 17,
            fullscreenControl : false,
            zoomControl : false,
            streetViewControl : false,
            mapTypeControl : false
          }"
          ref="gMap"
          style="width: 100%; height: 256px;">
      </GMapMap>
  
      <div id="CustomerInfo" class=" w-full">
        <div class="w-full h-2 border-t"></div>
        <div class="w-full text-center border-t-2 p-1.5 text-gray-700 text-lg font-semibold">
          Đang trên đường đến {{ destination.address }}
        </div>
        <!-- <div class="scrollSection">
          <div 
            v-for="(item, index) in items" 
            :key="index" 
            :style="{ backgroundColor: selectedItemIndex === index ? 'rgb(237, 237, 237)' : 'white' }" 
            @click="selectItem(index)" class="flex items-center px-4 py-5">
              <img width="75" :src="item.image">
              <div class="w-full ml-3">
                <div class="flex items-center justify-between">
                  <div class="text-2xl mb-1">{{ item.name }}</div>
                  <div class="text-xl">{{ calculatePrice(item.priceMultiplier, distance.value) }}</div>
                </div>
                <div class="text-gray-500">{{ item.description }}</div>
              </div> 
          </div>
        </div> -->
  
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
                bg-black 
                text-2xl 
                text-white
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
        }).catch((error) => {
            console.error(error)
            alert(error.response.data.message)
        }) 
  
        // draw a path on the map
        gMap.value.$mapPromise.then((mapObject) => {
            // eslint-disable-next-line
            let currentPoint = new google.maps.LatLng(location.current.geometry),
            // eslint-disable-next-line
                destinationPoint = new google.maps.LatLng(destination.value.geometry),
            // eslint-disable-next-line
                directionsService = new google.maps.DirectionsService,
            // eslint-disable-next-line
                directionsDisplay = new google.maps.DirectionsRenderer({
                    map: mapObject
                })
  
            directionsService.route({
                origin: currentPoint,
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
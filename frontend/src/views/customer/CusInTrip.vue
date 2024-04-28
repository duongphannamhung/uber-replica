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
          style="width: 100%; height: 65vh;">
      </GMapMap>

    <div id="DriverInfo" class="w-full">
      <div class="w-full h-2 border-t"></div>
      <div class="w-full text-center border-t-2 p-1.5 text-gray-700 text-lg font-semibold">
        Bạn đang đang trên đường đến {{ location.destination.name }}
      </div>
      <div class="border-b"></div>
      <div class="flex justify-between items-center mt-3">
        <div>
          <p class="font-bold text-lg ml-5">{{ driver_info.name }}</p>          
          <p class="ml-5">{{ driver_info.carLabel }} | {{ driver_info.carName }}</p>
          <div style="background-color: lightgray; padding: 4px; width: 125px; margin-left: 20px;">
            <p style="font-weight: 600;">{{ driver_info.carNumber }}</p>
          </div>
        </div>
        <div>
          <img :src="driver_info.image" alt="Driver Image" class="rounded-full w-20 h-20 mr-10">
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
                " 
            >
                Báo cáo
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
                @click="goToMessage()"
            >
                Tin nhắn
            </button>
          </div>
    </div>

    </div>
  </template>
  
  <script setup>
    // import axios from 'axios';
    import { onMounted, ref, watchEffect, inject } from 'vue'
    import { useRouter } from 'vue-router'
    import { useLocationStore } from '@/store/location'
  
    // import mapStyles from '../mapStyles'
    const m = ref(null); 
    const router = useRouter()
    const location = useLocationStore()
    const conn = inject('websocketStore').conn;
    // const trip = useTripStore()
  
    const gMap = ref(null)
    
    const driver_info = ref({
      name: 'DƯƠNG PHAN NAM HƯNG',
      carNumber: '59C1-123.45',
      carName: 'SH',
      carLabel: 'Honda',
      image: 'img/logo/driver.jpg'
    })

    var sleepSetTimeout_ctrl;

    function sleep(ms) {
        clearInterval(sleepSetTimeout_ctrl);
        return new Promise(resolve => sleepSetTimeout_ctrl = setTimeout(resolve, ms));
    }
    // const direction = useDirectionStore()
  
    watchEffect(() => {
      if (conn === null) {
        alert("Ws conn null. Can't completed trip")
      }
      conn.onmessage = (message) => {
        m.value = JSON.parse(message.data);

        if (m.value && m.value.content == 'finishline-@123!(*234kh219871233hadsfh') {
            conn.close();
            router.push({
                name: 'cus-complete-trip'
            })
        }
    }
    conn.onclose = () => {}
    conn.onerror = () => {}
    conn.onopen = () => {}
    })

    onMounted(async () => {

        // lets get the users current location
        await location.updateCurrentLocation()

        while (!gMap.value) {
          await sleep(1000);
        }
        // draw a path on the map
        gMap.value.$mapPromise.then((mapObject) => {
            // eslint-disable-next-line
            let currentPoint = new google.maps.LatLng(location.current.geometry),
            // eslint-disable-next-line
                destinationPoint = new google.maps.LatLng(location.destination.geometry),
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
          name : 'cus-message-in-trip'
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
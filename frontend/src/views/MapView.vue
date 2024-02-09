<template>
  <div>

    <div 
      id="BackBtn" 
      class="absolute z-50 rounded-full bg-white p-1 top-8 left-4"
      @click="goBack()"
    >
      <ArrowLeftIcon :size="40" />
    </div>
    
    <!-- <div id="map" > -->
    <GMapMap v-if="location.destination.name !== ''" 
        :zoom="4" :center="location.destination.geometry"
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

    <div id="VehicleSelection" class=" w-full">
      <div class="w-full h-2 border-t"></div>
      <div class="w-full text-center border-t-2 p-1.5 text-gray-700 text-lg font-semibold">
        Distance - {{ distance.text }}
      </div>

      <div class="scrollSection">

        <div class="bg-custom-gray">
          <div class="flex items-center px-4 py-5">
            <img width="75" src="img/uber/ride.png">
            <div class="w-full ml-3">
              <div class="flex items-center justify-between">
                <div class="text-2xl mb-1">UrepCar</div>
                <div class="text-xl">{{ calculatePrice(1.36, distance.value) }}</div>
              </div>
              <div class="text-gray-500">4-seater</div>
            </div>
          </div>
        </div>

        <div>
          <div class="flex items-center px-4 py-5">
            <img width="75" src="img/uber/bike.png">
            <div class="w-full ml-3">
              <div class="flex items-center justify-between">
                <div class="text-2xl mb-1">UrepBike</div>
                <div class="text-xl">{{ calculatePrice(0.75, distance.value) }}</div>
              </div>
              <div class="text-gray-500">Motor scooter</div>
            </div>
          </div> 
        </div>

        <div>
          <div class="flex items-center px-4 py-5">
            <img width="75" src="img/uber/uberxl.png">
            <div class="w-full ml-3">
              <div class="flex items-center justify-between">
                <div class="text-2xl mb-1">UrepCar 7</div>
                <div class="text-xl">{{ calculatePrice(1.51, distance.value) }}</div>
              </div>
              <div class="text-gray-500">7-seater</div>
            </div>
          </div>
        </div>

        <div>
          <div class="flex items-center px-4 py-5">
            <!-- <button @click="chooseDiv('UberXL')" -->
            <img width="75" src="img/uber/comfort.png">
            <div class="w-full ml-3">
              <div class="flex items-center justify-between">
                <div class="text-2xl mb-1">UrepCar Plus</div>
                <div class="text-xl">{{ calculatePrice(1.73, distance.value) }}</div>
              </div>
              <div class="text-gray-500">Professional service</div>
            </div>
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
        <button 
          @click="handleConfirmTrip"
          class="
            bg-black 
            text-2xl 
            text-white
            py-4 
            px-4 
            rounded-sm
            w-full
          "
        >
          Confirm UberX
        </button>
      </div>

    </div>
  </div>
</template>

<script setup>
  import axios from 'axios';
  import { onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import ArrowLeftIcon from 'vue-material-design-icons/ArrowLeft.vue';
  // import { useDirectionStore } from '@/store/direction-store';
  import { useLocationStore } from '@/store/location'
  // import { useTripStore } from '@/stores/trip'

  // import mapStyles from '../mapStyles'

  const router = useRouter()
  const location = useLocationStore()
  // const trip = useTripStore()

  const gMap = ref(null)
  
  // const direction = useDirectionStore()

  onMounted(async () => {
      // does the user have a location set?
      if (location.destination.name === '') {
          router.push({
              name: 'location'
          })
      }
      // lets get the users current location
      await location.updateCurrentLocation()

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
                  getDistance()
              } else {
                  console.error(status)
              }
          })
      })
  })

  const distance = ref({text: '', value: null})
  const duration = ref({text: '', value: null})
  // const latLng = ref({ start: { lat: null, lng: null }, end: { lat: null, lng: null } })

  // onMounted(async () => {
  //   if (!direction.pickup || !direction.destination) { router.push('/') }
  //   setTimeout(() => { initMap() }, 50)
  // })

  const goBack = () => {
    router.push('/location')
    // direction.pickup = ''
    // direction.destination = ''
  }

  const getDistance = async () => {
    // eslint-disable-next-line
    let currentPoint = new google.maps.LatLng(location.current.geometry)
    // eslint-disable-next-line
    let destinationPoint = new google.maps.LatLng(location.destination.geometry)
    let res = await axios.get('distance/' + currentPoint + '/' + destinationPoint, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      }
    })

    distance.value.text = res.data.distance_text
    distance.value.value = res.data.distance
    duration.value.text = res.data.duration_text
    duration.value.value = res.data.duration
  }

  const calculatePrice = (multiplier, distance) => {
    let res = (distance * 9000) * multiplier / 1000000
    if (res) {
      return convertPriceToVND(res.toFixed(0) * 1000)
    }
  }
  
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
      height: calc(50vh - 120px); 
      position: absolute; 
      overflow-y: auto; 
      width: 100%
    }

    .bg-custom-gray {
      background-color: rgb(237, 237, 237);
    }
  }
</style>
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
        " @click="handleConfirmTrip()"
      >
        <button  
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
          Book
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

  // const backgroundColor = ref({selected:'gray', unselected:'white'});
  const selectedItemIndex = ref(0);

  const items = [
    {
      name: 'UrepCar',
      description: '4-seater',
      priceMultiplier: 1.36,
      image: 'img/uber/car-4.png'
    },
    {
      name: 'UrepBike',
      description: 'Motor scooter',
      priceMultiplier: 0.75,
      image: 'img/uber/bike.png'
    },
    {
      name: 'UrepCar 7',
      description: '7-seater',
      priceMultiplier: 1.51,
      image: 'img/uber/uberxl.png'
    },
    {
      name: 'UrepCar Plus',
      description: 'Professional service',
      priceMultiplier: 1.73,
      image: 'img/uber/comfort.png'
    }
  ];

  const selectItem = (index) => {
    selectedItemIndex.value = index;
  }

  const router = useRouter()
  const location = useLocationStore()
  // const trip = useTripStore()

  const gMap = ref(null)

  var sleepSetTimeout_ctrl;

  function sleep(ms) {
      clearInterval(sleepSetTimeout_ctrl);
      return new Promise(resolve => sleepSetTimeout_ctrl = setTimeout(resolve, ms));
  }
  // const direction = useDirectionStore()

  onMounted(async () => {
      // does the user have a location set?
      if (location.destination.name === '') {
          router.push({
              name: 'cus-location'
          })
      }
      // lets get the users current location
      await location.updateCurrentLocation()

      while (!gMap.value) {
          await sleep(1000);
      }
      
      // draw a path on the map
      gMap.value.$mapPromise.then((mapObject) => {
          // eslint-disable-next-line
          let departurePoint = new google.maps.LatLng(location.departure.geometry),
          // eslint-disable-next-line
              destinationPoint = new google.maps.LatLng(location.destination.geometry),
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
    router.push({
      name : 'cus-location'
    })
    // direction.pickup = ''
    // direction.destination = ''
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

  const handleConfirmTrip = async () => {
    let item = items[selectedItemIndex.value]
    let vehicleName = item.name; 
    let tripRequest = {
      user_id: localStorage.getItem('current_user_id'),
      user_phone: localStorage.getItem('current_user_phone'),
      vehicle: vehicleName,
      // eslint-disable-next-line
      departure_point : new google.maps.LatLng(location.departure.geometry),
      departure_name: location.departure.display_name,
      // eslint-disable-next-line
      destination_point : new google.maps.LatLng(location.destination.geometry),
      destination_name: location.destination.name,
    }

    if (vehicleName == "UrepBike") {
      await axios.post('trip/bike', tripRequest, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('cus-token')}`
        }
    })
        .then(async (response) => {
          localStorage.setItem('current_trip_id', response.data.trip_id)

          let fare = ((distance.value.value * 9000) * item.priceMultiplier / 1000000).toFixed(0) * 1000
          console.log(`fare ${fare}, distance: ${distance.value.value}, multiplier: ${item.priceMultiplier}`)

          await axios.post('http://localhost:6969/api/driver/update-trip-fare', {
              trip_id: response.data.trip_id,
              fare: fare
            }, {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('cus-token')}`
              }
            })

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
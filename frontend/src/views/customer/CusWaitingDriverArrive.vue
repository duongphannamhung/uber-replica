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

    <div id="DriverInfo" class="w-full">
      <div class="w-full h-2 border-t"></div>
      <div class="w-full text-center border-t-2 p-1.5 text-gray-700 text-lg font-semibold">
        Tài xế đang đến đón bạn
      </div>
      <div class="border-b"></div>
      <div class="flex justify-between items-center mt-3">
        <div>
          <p class="font-bold text-lg ml-5">{{ driver_info.name.toUpperCase() }}</p>          
          <p class="ml-5">{{ driver_info.vehicleLabel }} | {{ driver_info.vehicleModel }} | {{ driver_info.vehicleColor }}</p>
          <div style="background-color: lightgray; padding: 4px; width: 125px; margin-left: 20px;">
            <p style="font-weight: 600;">{{ driver_info.vehicleLicensePlate }}</p>
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

  const router = useRouter()
  const location = useLocationStore()
  const driver_location = reactive({
      // name: '',
      // address: '',
      geometry: {
          lat: null,
          lng: null
      }
  })

  const driver_info = ref({
    name: '',
    vehicleLicensePlate: '',
    vehicleModel: '',
    vehicleLabel: '',
    vehicleColor: '',
    image: 'img/logo/driver.jpg'
  })

  // const trip = useTripStore()

  const gMap = ref(null)
  let driver_come_interval = null
  const driverArrived = ref(false); 
  // const direction = useDirectionStore()

  onMounted(async () => {
    await axios.get(`http://localhost:6969/api/trip/get-driver-info/${localStorage.getItem('current_trip_id')}`, {
          headers: {
              Authorization: `Bearer ${localStorage.getItem('cus-token')}`
          }
      }).then((resp) => {
          driver_info.value.name = resp.data.name
          driver_info.value.vehicleColor = resp.data.color
          driver_info.value.vehicleLabel = resp.data.label
          driver_info.value.vehicleModel = resp.data.model
          driver_info.value.vehicleLicensePlate = resp.data.license_plate

          localStorage.setItem('current_driver_vehicle_model', resp.data.model) 
          localStorage.setItem('current_driver_name', resp.data.name) 
          localStorage.setItem('current_vehicle_license_plate', resp.data.license_plate) 
          localStorage.setItem('current_driver_vehicle_color', resp.data.color) 
          localStorage.setItem('current_driver_vehicle_label', resp.data.label) 
      }).catch((error) => {
          console.error(error)
          alert(error.response.data.message)
      })

      await axios.get(`http://localhost:6969/api/trip/${localStorage.getItem('current_trip_id')}`, {
          headers: {
              Authorization: `Bearer ${localStorage.getItem('cus-token')}`
          }
      }).then((resp) => {
          driver_location.geometry.lat = resp.data.driver_location_latitude.Float64
          driver_location.geometry.lng = resp.data.driver_location_longitude.Float64
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
          let departurePoint =  new google.maps.LatLng(driver_location.geometry),
          // eslint-disable-next-line
              destinationPoint = new google.maps.LatLng(location.departure.geometry),
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
          name : 'cus-message-in-trip'
      })
  }

  const updateDriverLocation = async () => {
    if (getDistanceFromLatLonInKm(location.departure.geometry.lat, location.departure.geometry.lng, driver_location.geometry.lat, driver_location.geometry.lng) < 100) {
      clearInterval(driver_come_interval)

      // await sleep(3000);
      driverArrived.value = true; 
    }
  }

  const startTrip = async () => {
    localStorage.setItem('driver_arrived', 'true')
    router.push({
        name : 'cus-in-trip'
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
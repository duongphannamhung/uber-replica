<template>
  <div>
    <GMapMap v-if='checked'
        :zoom="15" 
        :center="location.current.geometry"
        :options="{
          minZoom: 3,
          maxZoom : 17,
          fullscreenControl : false,
          zoomControl : false,
          streetViewControl : false,
          mapTypeControl : false
        }"
        ref="gMap"
        style="width: 100%; height: 512px;">
        <GMapMarker :position="location.current.geometry"/>
    </GMapMap>

    <GMapMap v-else
        :zoom="15" 
        :center="location.current.geometry"
        :options="{
          minZoom: 3,
          maxZoom : 17,
          fullscreenControl : false,
          zoomControl : false,
          streetViewControl : false,
          mapTypeControl : false,
          styles: [
            {
              featureType: 'all',
              stylers: [
                { saturation: -100 },
                { lightness: 25 }
              ]
            }
          ]
        }"
        ref="gMap"
        style="width: 100%; height: 512px;">
    </GMapMap>

    <div class="toggle-container">
          <DriverToggle @click="toggle()"/>
    </div>
  </div>
</template>
  
  <script setup>
    import axios from 'axios';
    import { onMounted, ref } from 'vue'
    import DriverToggle from '@/components/DriverToggle.vue';
    import { useLocationStore } from '@/store/location'
import router from '@/router';

    // import { useTripStore } from '@/stores/trip'
  
    const checked = ref(false)
    const location = useLocationStore()
    let intervalId = null;
    let intervalGetStatus = null;

    const isActive = () => {
      if (checked.value) {
        return 2
      } else {
        return 1
      }
    }

    const toggle = () => {
        checked.value = !checked.value;
        if (checked.value) {
          intervalId = setInterval(updateEngagement, 1000)
          intervalGetStatus = setInterval(getDriverStatus, 1000)
        } else {
          clearInterval(intervalId);
          clearInterval(intervalGetStatus);
          intervalId = null;
          intervalGetStatus = null;

          updateEngagement();
        }
    }

    const getEngagement = () => {
      return {
        driver_id: localStorage.getItem('current_driver_id'),
        driver_phone: localStorage.getItem('current_driver_phone'),
        status: isActive(),
        lat : location.current.geometry.lat,
        lng : location.current.geometry.lng,
        geo_id : 1 // update geo_id later
      }
    }

    const updateEngagement = async () => {
      await axios.post('http://localhost:6969/api/driver/update-engagement', getEngagement(), {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('driver-token')}`
        }
    })
      // .then((response) => {
      // })
      .catch((error) => {
          console.error(error)
          alert(error.response.data.message)
      })
  }

  const getDriverStatus = async () => {
    await axios.get('http://localhost:6969/api/driver/current-status?driver_id=' + localStorage.getItem('current_driver_id') , {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('driver-token')}`
        }
    })
      .then((response) => {
        if (response.data.status !== 1 && response.data.status !== 2) {
          clearInterval(intervalId);
          clearInterval(intervalGetStatus);
          intervalId = null;
          intervalGetStatus = null;
          router.push({
            name: 'driver-drive-to-cus'
          })
        }
      })
      .catch((error) => {
          console.error(error)
          alert(error.response.data.message)
      })
  }

    onMounted(async () => {
      await location.updateCurrentLocation()
    })
  </script>

<style scoped>
.toggle-container {
  position: absolute;
  top: 10px;
  right: 10px;
}
</style>
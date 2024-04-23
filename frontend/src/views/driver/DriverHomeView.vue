<template>
  <div>
    <GMapMap v-if='isToggleOn'
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
          <DriverToggle v-model="isToggleOn" @click="toggle()"/>
    </div>
  </div>
</template>
  
  <script setup>
    import axios from 'axios';
    import { onMounted, ref } from 'vue'
    // import { onBeforeRouteUpdate } from 'vue-router'
    import DriverToggle from '@/components/DriverToggle.vue';
    import { useLocationStore } from '@/store/location'
    import { websocketStore } from '@/store/websocket-store'
    import router from '@/router';

    // import { useTripStore } from '@/stores/trip'
    // const checked = ref(false)
    const location = useLocationStore()
    let intervalId = null;
    let intervalGetStatus = null;
    let isToggleOn = ref(false)

    // if (localStorage.getItem('isToggleOn') === 'true') {
    //   localStorage.setItem('isToggleOn', 'true')
    //   intervalId = setInterval(() => {
    //         updateEngagement(getEngagement());
    //       }, 1000)
    //   intervalGetStatus = setInterval(getDriverStatus, 1000)
    // } else {
    //   isToggleOn.value = false
    // }



  //   onBeforeRouteEnter((to, from, next) => {
  //   if (localStorage.getItem('after_trip') === true) {
  //     checked.value = true;
  //     localStorage.setItem('after_trip', null)
  //   }
  //   next();
  // });

  //   onBeforeRouteUpdate((to, from, next) => {
  //     if (localStorage.getItem('after_trip') === true) {
  //       isToggleOn.value = true
  //       localStorage.setItem('after_trip', null)
  //   }
  //   next();
  // })

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

  const updateEngagement = async (engagementInfo) => {
      await axios.post('http://localhost:6969/api/driver/update-engagement', engagementInfo, {
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

          getTripInfo(response.data.trip_id)
        }
      })
      .catch((error) => {
          console.error(error)
          alert(error.response.data.message)
      })
  }

    const isActive = () => {
      if (isToggleOn.value) {
        return 2
      } else {
        return 1
      }
    }

    const toggle = () => {
      isToggleOn.value = !isToggleOn.value;
        if (isToggleOn.value) {
          localStorage.setItem('isToggleOn', 'true')
          intervalId = setInterval(() => {
            updateEngagement(getEngagement());
          }, 1000)
          intervalGetStatus = setInterval(getDriverStatus, 1000)
        } else {
          localStorage.setItem('isToggleOn', 'false')
          clearInterval(intervalId);
          clearInterval(intervalGetStatus);
          intervalId = null;
          intervalGetStatus = null;

          updateEngagement(getEngagement());
        }
    }

    const getTripInfo = async (trip_id) => {
    await axios.get(`http://localhost:6969/api/trip/${trip_id}`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('driver-token')}`
      }
    }).then((resp) => {
      const ws = new WebSocket(`ws://localhost:6969/ws/join-room/${resp.data.trip_id}?userId=${resp.data.driver_id}&phoneNumber="${localStorage.getItem('current_driver_phone')}"&isCustomer=false`);
      websocketStore.setConn(ws);
      if (websocketStore.conn && websocketStore.conn.OPEN) {
        router.push({
          name: 'driver-drive-to-cus'
        })
      }
    }).catch((error) => {
      console.error(error)
      alert(error.response.data.message)
    })
  }

    if (localStorage.getItem('after_trip') == 'true') {
      toggle()
    }
    // window.addEventListener('beforeunload', function (e) {
      // if (intervalId) {
      //   clearInterval(intervalId);
      //   clearInterval(intervalGetStatus);
      //   intervalId = null;
      //   intervalGetStatus = null;
      //   updateEngagement();
      // }

  onMounted(async () => {
    await location.updateCurrentLocation()

    window.addEventListener('beforeunload', function () {
      updateEngagement({
        driver_id: localStorage.getItem('current_driver_id'),
        driver_phone: localStorage.getItem('current_driver_phone'),
        status: 1,
        lat : location.current.geometry.lat,
        lng : location.current.geometry.lng,
        geo_id : 1 // update geo_id later
      }) 
    })
  })
  </script>

<style scoped>
.toggle-container {
  position: absolute;
  top: 10px;
  right: 10px;
}
</style>
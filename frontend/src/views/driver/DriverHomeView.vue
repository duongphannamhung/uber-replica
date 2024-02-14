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
    // import axios from 'axios';
    import { onMounted, ref } from 'vue'
    import DriverToggle from '@/components/DriverToggle.vue';
    import { useLocationStore } from '@/store/location'

    // import { useTripStore } from '@/stores/trip'
  
    const checked = ref(false)
    const location = useLocationStore()

    const toggle = () => {
        checked.value = !checked.value;
        console.log(checked.value)
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
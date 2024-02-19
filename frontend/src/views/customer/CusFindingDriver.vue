<!-- <template>

    <div>
        <p>Đang tìm kiếm tài xế cho chuyến đi của bạn..</p>
    </div>
    <div 
      id="BackBtn" 
      class="absolute z-50 rounded-full bg-white p-1 top-8 left-4"
      @click="goBack()"
    >
      <ArrowLeftIcon :size="40" />
    </div> -->




    <!-- <div class="flex justify-between">
        <button
            @click="goBack"
            class="rounded-md border border-transparent bg-black py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-gray-600 focus:outline-none">
            Trở về
        </button>
    </div> -->

    <!-- <div v-else>
      <div>Here's the data!</div>
      <pre>{{ data.toString() }}</pre>
    </div> -->




  <!-- </template>
  
  <script setup>
    import { useRouter } from 'vue-router'
 -->




<!-- //   import { ref , onMounted } from 'vue';
  
//   export default {
//     setup() {
//     const data = ref(null);
  
//       onMounted(() => {
//             data.value = null;
//             console.log("here")
// //         // Replace this `fetch` call with whatever your endpoint call may be.
// //         fetch('./endpoint')
// //           .then(resp => resp.json())
// //           .then(json => data.value = json);
//       });
  
//       return { data };
//     }
//   };
    const router = useRouter()

    const goBack = () => {
        router.push({
            name: 'cus-login'
        })
    } -->








<template>
    <div class="pt-16 flex flex-col items-center justify-start h-screen">
        <img src="img/uber/urep-user-find-driver.gif" alt="Loading..." />
        <h1 v-if="i === 0" class="mt-4 text-3xl font-medium mb-4">Đang tìm kiếm tài xế</h1>
        <h1 v-else-if="i === 1" class="mt-4 text-3xl font-medium mb-4">Đang tìm kiếm tài xế.</h1> 
        <h1 v-else-if="i === 2" class="mt-4 text-3xl font-medium mb-4">Đang tìm kiếm tài xế..</h1> 
        <h1 v-else class="mt-4 text-3xl font-medium mb-4">Đang tìm kiếm tài xế...</h1>
        <button
            @click="goBack"
            class="mt-4 rounded-md border border-transparent bg-black py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-gray-600 focus:outline-none">
            Hủy và trở về
        </button>
    </div>
</template>
<script setup>
import { onMounted, onUnmounted, ref } from 'vue';
import { useRouter } from 'vue-router'
import axios from 'axios';
import { websocketStore } from '@/store/websocket-store'

const router = useRouter()
const goBack = () => {
    router.push({
        name: 'cus-map'
    })
}

const i = ref(0);
let _interval = null;
let find_driver_interval = null;

onMounted(() => {
    _interval = setInterval(updateSecond, 1000);
    find_driver_interval = setInterval(findDriver, 1000);
})

onUnmounted(() => {
    clearInterval(_interval)
    clearInterval(find_driver_interval)
})

const updateSecond = () => {
    i.value += 1;
    i.value = i.value % 4;    
}

const findDriver = async () => {
    await axios.get('http://localhost:6969/api/trip/find-driver?trip_id=' + localStorage.getItem('current_trip_id') , {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('cus-token')}`
        }
    })
        .then(async (response) => {
            if (response.data.find_done) {
                await axios.post('http://localhost:6969/ws/create-room', {
                    id: localStorage.getItem('current_trip_id'),
                    customer_id: localStorage.getItem('current_user_id'),
                    driver_id: response.data.driver_id.toString()
                }).then((resp) => {
                    const ws = new WebSocket(`ws://localhost:6969/ws/join-room/${resp.data.id}?userId=${resp.data.customer_id}&phoneNumber="${localStorage.getItem('current_user_phone')}"&isCustomer=true`);
                    websocketStore.setConn(ws);
                    if (websocketStore.conn && websocketStore.conn.OPEN) {
                        router.push({
                                name : 'cus-waiting-driver-arrive'
                            })
                        }
                    }
                ).catch((error) => {
                    console.error(error)
                    alert(error.response.data.message)
                    goBack()
                })
            }
        })
        .catch((error) => {
            console.error(error)
            alert(error.response.data.message)
            goBack()
        })
}


</script>
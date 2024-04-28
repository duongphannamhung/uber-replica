<template>
    <div class="w-full flex flex-col justify-between border-t-2 p-1.5 text-gray-700 text-lg font-semibold">
        <div class="flex-1 text-center">
            Hoàn thành chuyến #{{ metadata.trip_id }}
        </div>
        <div class="flex-1 text-center">
            Đến: {{ metadata.destination_name }}
        </div>
        <div class="flex-1 text-center">
            <div>
                Cước phí: {{ metadata.fare }} VND
            </div>
            <div>
                Thời gian: {{ metadata.trip_created_at }}
            </div>
        </div>
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
                w-1/4
                mr-2
                " @click="goToMessage()"
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
                w-3/4
                ml-2
                "
                @click="backHome()"
            >
                Hoàn thành chuyến
            </button>
            </div>
    </div>
</template>
<script setup>
    import axios from 'axios';
    import { onMounted, ref } from 'vue'
    import { useRouter } from 'vue-router'

    const metadata = ref({
        trip_id: null,
        fare: null,
        destination_name: null,
        trip_created_at: null
    })
    const router = useRouter()

    onMounted(async () => {
        await axios.get(`http://localhost:6969/api/trip/${localStorage.getItem('current_trip_id')}`, {
          headers: {
              Authorization: `Bearer ${localStorage.getItem('cus-token')}`
          }
      }).then((resp) => {
            metadata.value.trip_id = resp.data.trip_id
            metadata.value.fare = resp.data.fare
            metadata.value.destination_name = resp.data.destination_name
            metadata.value.trip_created_at = resp.data.trip_created_at
      }).catch((error) => {
          console.error(error)
          alert(error.response.data.message)
      })
    })

    const backHome = async () => {
        localStorage.removeItem('current_trip_id')
        localStorage.removeItem('driver_arrived')
        router.push({
            name : 'cus-home'
        })
    }
</script>
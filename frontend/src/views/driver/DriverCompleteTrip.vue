<template>
    <div class="w-full text-center border-t-2 p-1.5 text-gray-700 text-lg font-semibold">
        Thông tin chuyến đi
    </div>
    <div id="DriverInfo" class="w-full">
        <div class="w-full h-1 border-t border-t bg-gray-300"></div>
            <div class="border-b"></div>
            <div class="flex justify-between items-center mt-3">
            <div>
                <p class="font-bold text-lg ml-5">{{ customer_info.address }}</p>
                <p class="ml-5">{{ customer_info.name }}</p>
                <div style="background-color: white; padding: 4px; width: 175px; margin-left: 20px;">
            </div>
            </div>
            <div>
                <img :src="customer_info.image" alt="Customer Image" class="rounded-full w-20 h-20 mr-10">
            </div>
            </div>
            <div class="border-b border-gray-200 my-2"></div>
            <button 
                    class="
                    bg-white
                    text-lg
                    text-black
                    rounded-sm
                    w-full
                    ml-2
                    mb-2
                    " 
                >
                    ★ Đánh giá khách hàng
                </button>
    </div>

    <div id="Payment" class="w-full mb-3">
        <div class="w-full h-1 border-t border-t bg-gray-300"></div>
            <div class="border-b"></div>
            <div class="flex flex-col">
            <p class="font-bold text-lg mt-2 ml-5">Payment</p>
            <div class="flex justify-between items-center mx-5 mt-1">
                <p>Cước chuyến đi</p> 
                <p>{{ convertPriceToVND(metadata.fare) }}</p>
            </div>
            <div class="flex justify-between items-center mx-5 mt-1">
                <p>Giảm giá</p> 
                <p>- 0 đ</p>
            </div>
            <div class="flex justify-between items-center mx-5 mt-1">
                <p>Bảo hiểm chuyến đi</p> 
                <p>1.000 đ</p>
            </div>
            </div>
            <div class="border-b border-gray-200 my-2 ml-5 mr-5"></div>
            <div class="font-bold text-l flex justify-between items-center mx-5">
                <p>Tiền mặt phải thu</p> 
                <p>{{ convertPriceToVND(metadata.fare + 1000) }}</p>
            </div>
    </div>

    <div id="TripInfo" class="w-full mb-5">
        <div class="w-full h-1 border-t border-t bg-gray-300"></div>
            <div class="border-b"></div>
            <div class="flex justify-between items-center mt-3 mb-3">
            <div>
                <p class="font-bold text-lg ml-5">Trip ID: {{ metadata.trip_id }}</p>          
                <p class="ml-5">{{ date_trip_info.date }} | {{ date_trip_info.time }}</p>
            </div>
        </div>

            <div class="flex flex-col">
                <div class="flex flex-col space-y-5 p-5">
                <div class="flex flex-col">
                <div class="flex items-center">
                    <img src="img/logo/start_point.png" alt="Start Icon" class="mr-2 w-6 h-6">
                    <p class="font-bold text-lg">Lorem ipsum</p>
                    </div>
                    <p class="text-lg ml-8">Quận 1, Hồ Chí Minh</p>
                </div>
                <div class="border-b border-gray-400 my-2"></div>
                <div class="flex flex-col">
                <div class="flex items-center">
                    <img src="img/logo/destination_point.png" alt="Destination Icon" class="mr-2 w-6 h-6">
                    <p class="font-bold text-lg">{{ metadata.destination_name }}</p>
                    </div>
                    <p class="text-lg ml-8">Quận 1, Hồ Chí Minh</p>
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

    const customer_info = ref({
      name: 'Dương Phan Nam Hưng',
      address: '7/42 Thành Thái, Q.10, HCM', // TODO: get real info
      image: 'https://www.w3schools.com/howto/img_avatar.png',
      fare: 0,
    })

    
    const date_trip_info = ref({
        date: null,
        time: null
    })

    const timestampToFormat = (_datetime) => {
        if (_datetime == null) return ''
        let dt = _datetime.toString()
        date_trip_info.value.date = dt.split(' ')[0]
        let _time = dt.split(' ')[1]
        date_trip_info.value.time = _time.split(':')[0] + ':' + _time.split(':')[1]
    }

    const convertPriceToVND = (price) => {
    if (price == null) return '0 đ'
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

    onMounted(async () => {
        await axios.post(`http://localhost:6969/api/driver/finish-engagement`, {
            driver_phone: localStorage.getItem('current_driver_phone'), driver_id: localStorage.getItem('current_driver_id'), trip_id: localStorage.getItem('current_trip_id')
        }, {
        headers: {
                Authorization: `Bearer ${localStorage.getItem('driver-token')}`
            }
        }).then((resp) => {
            metadata.value.trip_id = resp.data.trip_id
            metadata.value.fare = resp.data.fare
            metadata.value.destination_name = resp.data.destination_name
            metadata.value.trip_created_at = resp.data.trip_created_at
            timestampToFormat(metadata.value.trip_created_at)
        }).catch((error) => {
            console.error(error)
            alert(error.response.data.message)
        }) 
    })

    const backHome = async () => {
        localStorage.removeItem('driver_arrived')
        localStorage.removeItem('current_trip_id')
        localStorage.setItem('after_trip', true)
        router.push({
            name : 'driver-home'
        })
    }
</script>
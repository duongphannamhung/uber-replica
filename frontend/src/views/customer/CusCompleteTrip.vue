<template>
    <div class="w-full text-center border-t-2 p-1.5 text-gray-700 text-lg font-semibold">
        Thông tin chuyến đi
    </div>
    <div id="DriverInfo" class="w-full">
        <div class="w-full h-1 border-t border-t bg-gray-300"></div>
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
                    ★ Đánh giá tài xế
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
                <p>Tiền mặt phải trả</p> 
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
                    <p class="font-bold text-lg">{{ beautifulizeAddress(metadata.departure_name) }}</p>
                    </div>
                    <p class="text-lg ml-8">{{ tailAddress(metadata.departure_name) }}</p>
                </div>
                <div class="border-b border-gray-400 my-2"></div>
                <div class="flex flex-col">
                <div class="flex items-center">
                    <img src="img/logo/destination_point.png" alt="Destination Icon" class="mr-2 w-6 h-6">
                    <p class="font-bold text-lg">{{ beautifulizeAddress(metadata.destination_name) }}</p>
                    </div>
                    <p class="text-lg ml-8">{{ tailAddress(metadata.destination_name) }}</p>
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

    const driver_info = ref({
      name: 'DƯƠNG PHAN NAM HƯNG',
      carNumber: '59C1-123.45',
      carName: 'SH',
      carLabel: 'Honda',
      image: 'img/logo/driver.jpg'
    })
    
    const date_trip_info = ref({
        date: null,
        time: null
    })

    const convertPriceToVND = (price) => {
        if (price == null) return ''
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

    const metadata = ref({
        trip_id: null,
        fare: null,
        destination_name: null,
        departure_name: null,
        trip_created_at: null
    })
    const router = useRouter()

    const timestampToFormat = (_datetime) => {
        if (_datetime == null) return ''
        let dt = _datetime.toString()
        date_trip_info.value.date = dt.split(' ')[0]
        let _time = dt.split(' ')[1]
        date_trip_info.value.time = _time.split(':')[0] + ':' + _time.split(':')[1]
    }

    onMounted(async () => {
        await axios.get(`http://localhost:6969/api/trip/${localStorage.getItem('current_trip_id')}`, {
          headers: {
              Authorization: `Bearer ${localStorage.getItem('cus-token')}`
          }
      }).then((resp) => {
            metadata.value.trip_id = resp.data.trip_id
            metadata.value.fare = resp.data.fare
            metadata.value.destination_name = resp.data.destination_name
            metadata.value.departure_name = resp.data.departure_name
            metadata.value.trip_created_at = resp.data.trip_created_at
            timestampToFormat(metadata.value.trip_created_at)
      }).catch((error) => {
          console.error(error)
          alert(error.response.data.message)
      })
    })

    const beautifulizeAddress = (address) => {
        if (!address) return ''
        let list_address = address.split(',')
        if (list_address.length < 2) {
        return address
        }
        else {
        return list_address[0] + ', ' + list_address[1] + ', ' + list_address[2]
        }
    }

    const tailAddress = (address) => {
        if (!address) return ''
        let list_address = address.split(',')
        if (list_address.length < 2) {
            return address
        }
        else {
            if (Number.isInteger(Number(list_address[list_address.length - 2]))) {
                return list_address[list_address.length - 3] + ', ' + list_address[list_address.length - 1]
            }
            else {
                return list_address[list_address.length - 2] + ', ' + list_address[list_address.length - 1]
            }
        }
    }

    const backHome = async () => {
        localStorage.removeItem('current_trip_id')
        localStorage.removeItem('driver_arrived')
        router.push({
            name : 'cus-home'
        })
    }
</script>
<template>
    <div class="pt-16">
        <h1 class="text-3xl font-semibold mb-4 text-center">Đăng ký tài xế mới</h1>
        <form action="#" @submit.prevent="handleRegister">
            <div class="overflow-hidden shadow sm:rounded-md max-w-sm mx-auto text-left">
                <div class="bg-white px-4 py-5 sm:p-6">
                    <div>
                        <label for="driverName" class="block text-sm font-medium text-gray-700">Họ và tên:</label>
                        <input type="text" v-model="driver.name" name="driverName" id="driverName"
                            class="mt-1 block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-black focus:outline-none">
                    </div>
                    <div class="mt-4">
                        <label for="vehicleType" class="block text-sm font-medium text-gray-700">Dịch vụ:</label>
                        <select v-model="driver.vehicleType" name="vehicleType" id="vehicleType"
                            class="mt-1 block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-black focus:outline-none">
                            <option value="1">Xe máy: UrepBike</option>
                            <option value="2">Ô tô 4 chỗ: UrepCar</option>
                            <option value="3">Ô tô 7 chỗ: UrepCar 7</option>
                            <option value="4">Ô tô cao cấp: UrepCar Plus</option>
                        </select>
                    </div>
                    <div class="mt-4">
                        <label for="vehicleLabel" class="block text-sm font-medium text-gray-700">Hãng xe:</label>
                        <input type="text" v-model="driver.vehicleLabel" name="vehicleLabel" id="vehicleLabel"
                            class="mt-1 block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-black focus:outline-none">
                    </div>
                    <div class="mt-4">
                        <label for="vehicleModel" class="block text-sm font-medium text-gray-700">Dòng xe:</label>
                        <input type="text" v-model="driver.vehicleModel" name="vehicleModel" id="vehicleModel"
                            class="mt-1 block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-black focus:outline-none">
                    </div>
                    <div class="mt-4">
                        <label for="vehicleColor" class="block text-sm font-medium text-gray-700">Màu xe:</label>
                        <input type="text" v-model="driver.vehicleColor" name="vehicleColor" id="vehicleColor"
                            class="mt-1 block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-black focus:outline-none">
                    </div>
                    <div class="mt-4">
                        <label for="vehicleLicensePlate" class="block text-sm font-medium text-gray-700">Biển số xe:</label>
                        <input type="text" v-model="driver.vehicleLicensePlate" name="vehicleLicensePlate" id="vehicleLicensePlate"
                            class="mt-1 block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-black focus:outline-none">
                    </div>
                </div>
                <div class="bg-gray-50 px-4 py-3 text-right sm:px-6">
                    <button type="submit"
                        class="inline-flex justify-center rounded-md border border-transparent bg-black py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-gray-600 focus:outline-none">Đăng ký</button>
                </div>
            </div>
        </form>
    </div>
</template>

<script setup>
import axios from 'axios';
import { useRouter } from 'vue-router'
import { ref } from 'vue'

const router = useRouter()
const driver = ref({
        name: '',
        vehicleType: '',
        vehicleLabel: '',
        vehicleModel: '',
        vehicleColor: '',
        vehicleLicensePlate: ''
    })

const handleRegister = () => {
    if (isFormValid()) {
// Handle the onboarding process here
    } else {
        alert('Hãy nhập đủ thông tin.')
        return
    }
    // Handle the onboarding process here
    axios.post('http://localhost:6969/api/driver/register', {
        driver_id: parseInt(localStorage.getItem('current_driver_id')),
        name: driver.value.name,
        vehicle_type: parseInt(driver.value.vehicleType),
        vehicle_label: driver.value.vehicleLabel,
        vehicle_model: driver.value.vehicleModel,
        vehicle_color: driver.value.vehicleColor,
        vehicle_plate: driver.value.vehicleLicensePlate
    })
        .then(() => {
            router.push({
                name: 'driver-home'
            })
        })
        .catch((error) => {
            console.error(error)
            alert(error.response.data.message)
        })
}

const isFormValid = () => {
    return Object.values(driver).every(value => value !== '')
}

</script>
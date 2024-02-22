<template>
    <div class="pt-16">
        <h1 class="text-3xl font-semibold mb-4 text-center">Urep</h1>
        <div class="overflow-hidden shadow sm:rounded-md max-w-sm mx-auto text-left">
            <div class="bg-white px-4 py-5 sm:p-6">
                <div class="flex justify-between">
                    <button
                        @click="handleLoginCustomer"
                        class="rounded-md border border-transparent bg-black py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-gray-600 focus:outline-none">
                        Customer App
                    </button>
                    <button
                        @click="handleLoginDriver"
                        class="rounded-md border border-transparent bg-black py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-gray-600 focus:outline-none">
                        Driver App
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
const handleLoginDriver = () => {
    router.push({
        name: 'driver-login'
    })
}
const handleLoginCustomer = () => {
    router.push({
        name: 'cus-login'
    })
}
</script>
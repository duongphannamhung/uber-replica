<template>
    <div class="pt-16">
        <h1 class="text-3xl font-semibold mb-4">[Customer] Nhập số điện thoại của bạn</h1>
        <form v-if="!waitingOnVerification" action="#" @submit.prevent="handleLogin">
            <div class="overflow-hidden shadow sm:rounded-md max-w-sm mx-auto text-left">
                <div class="bg-white px-4 py-5 sm:p-6">
                    <div>
                        <input type="text" v-maska data-maska="(+84)-###-###-###" v-model="credentials.phone" name="phone" id="phone" placeholder="(+84)-987-123-456"
                            class="mt-1 block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-black focus:outline-none">
                    </div>
                </div>
                <div class="bg-gray-50 px-4 py-3 text-right sm:px-6">
                    <button type="submit" @submit.prevent="handleLogin"
                        class="inline-flex justify-center rounded-md border border-transparent bg-black py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-gray-600 focus:outline-none">Continue</button>
                </div>
            </div>
        </form>
        <form v-else action="#" @submit.prevent="handleVerification">
            <div class="overflow-hidden shadow sm:rounded-md max-w-sm mx-auto text-left">
                <div class="bg-white px-4 py-5 sm:p-6">
                    <div>
                        <input type="text" v-maska data-maska="######" v-model="credentials.login_code" name="login_code" id="login_code" placeholder="Nhập OTP"
                            class="mt-1 block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-black focus:outline-none">
                    </div>
                </div>
                <div class="bg-gray-50 px-4 py-3 text-right sm:px-6">
                    <button type="submit" @submit.prevent="handleVerification"
                        class="inline-flex justify-center rounded-md border border-transparent bg-black py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-gray-600 focus:outline-none">Verify</button>
                </div>
            </div>
        </form>
    </div>
</template>
<script setup>
import { vMaska } from 'maska'
// eslint-disable-next-line
import { reactive, ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router';

const router = useRouter()

const credentials = reactive({
    phone: null,
    login_code: null
})

const waitingOnVerification = ref(false)

onMounted(() => {
    if (localStorage.getItem('cus-token')) {
        router.push({
            name: 'cus-home'
        })
    }
})

const getFormattedCredentials = () => {
    return {
        phone: credentials.phone.replaceAll('(+84)','+84').replaceAll(' ', '').replace('(', '').replace(')', '').replaceAll('-', ''),
        login_code: credentials.login_code
    }
}

const handleLogin = () => {
    // TODO: change this to .env
    axios.post('http://localhost:6969/api/login-phone', getFormattedCredentials())
        .then((response) => {
            console.log(response.data)
            // waitingOnVerification.value = false
            waitingOnVerification.value = true
        })
        .catch((error) => {
            console.error(error)
            alert(error.response.data.message)
        })
}

const handleVerification = () => {
    // TODO : change this to .env
    axios.post('http://localhost:6969/api/login-phone/verify', getFormattedCredentials())
        .then((response) => {
            localStorage.setItem('cus-token', response.data.access_token)
            localStorage.setItem('current_user_phone', response.data.user.phone)
            localStorage.setItem('current_user_id', response.data.user.id)
            router.push({
                name: 'cus-home'
            })
        })
        .catch((error) => {
            console.error(error)
            alert(error.response.data.message)
        })
}
</script>
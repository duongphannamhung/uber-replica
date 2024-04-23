<template>
    <div>
        <h2 class="mb-4"></h2>
        <Toggle v-model="isToggleOn" class="blue-shadow"/>
    </div>
</template>
          
<script>
import Toggle from '@vueform/toggle'
export default {
    components: {
        Toggle,
    },
    data() {
        return {
            isToggleOn: null,
        }
    },
    methods: {
        checkCondition() {
        if (localStorage.getItem('after_trip') == 'true') {
            this.isToggleOn = true;
        } else {
            this.isToggleOn = false;
        }
    },
    },
    mounted() {
        this.checkCondition()
    }
}
</script>

<style scoped>
.blue-shadow {
    border-radius: 5%;
    box-shadow: 0 0 5px 5px lightgreen;
    overflow: hidden;
}
</style>

<style src="@vueform/toggle/themes/default.css"></style>
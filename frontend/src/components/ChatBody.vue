<template>
    <div id="ChatBody">
      <div
        v-for="(message, index) in data"
        :key="index"
        :class="message.is_customer == 'true' ? 'flex flex-col mt-2 w-full text-right justify-end' : 'mt-2'"
      >
        <div class="text-sm">{{ message.user_id }}</div>
        <div>
          <div
            class="px-4 py-1 rounded-md inline-block mt-1"
            :class="message.is_customer == 'true' ? 'bg-blue text-white' : 'bg-grey text-dark-secondary'"
          >
            {{ message.content }}
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      data: {
        type: Array,
        required: true
      }
    }
  };
  </script>

<!-- 
  <template>
    <div id="ChatBody">
      <div
        v-for="(message, index) in data"
        :key="index"
      >
      <div :class="'flex flex-col mt-2 w-full text-right justify-end'">
        <div>
          <div
            class="px-4 py-1 rounded-md inline-block mt-1"
            :class="'bg-blue text-white'"
          >
          </div>
        </div>
      </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      data: {
        type: Array,
        required: true
      }
    }
  };
  </script> -->
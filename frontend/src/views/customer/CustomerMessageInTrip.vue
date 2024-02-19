<template>
    <div class="flex flex-col w-full">
    <div class="flex justify-center items-center p-4 md:mx-6">
        <div 
            id="BackBtn" 
            class="fixed z-50 rounded-full bg-white p-1 top-8 left-4"
            @click="goBack()"
            >
            <ArrowLeftIcon :size="40" />
        </div>
        <h1 class="text-center">Tin nhắn</h1>
    </div>
      <div class="p-4 md:mx-6 mb-14">
        <div 
          v-for="(message, index) in messages" 
          :key="index" 
          :class="message.is_customer ? 'flex flex-col mt-2 w-full text-right justify-end' : 'mt-2'"
          >
            <div v-if="message.is_customer" class="text-sm">You</div>
            <div v-else class="text-sm">Driver id {{ message.user_id }}</div>
          <div>
          <div
            class="px-4 py-1 rounded-md inline-block mt-1"
            :class="message.is_customer ? 'bg-blue text-black bg-blue-200' : 'bg-grey text-dark-secondary bg-gray-200'"
          >
            {{ message.content }}
          </div>
        </div>
        </div>
      </div>
      <div class="fixed bottom-0 mt-4 w-full">
        <div class="flex md:flex-row px-4 py-2 bg-grey md:mx-4 rounded-md">
          <div class="flex w-full mr-4 rounded-md border border-blue">
            <textarea
              ref="textareaRef"
              placeholder="type your message here"
              class="w-full h-10 p-2 rounded-md focus:outline-none"
              style="resize: none"></textarea>  
          </div>
          <div class="flex items-center">
            <button class="p-2 rounded-md bg-blue text-black" @click="sendMessage">
              Send
            </button>
          </div>
        </div>
      </div>
    </div>
</template>

<script>
    import router from '@/router';
    import { ref, watchEffect, inject } from 'vue';
    import ArrowLeftIcon from 'vue-material-design-icons/ArrowLeft.vue';

    export default {
        components: {
            ArrowLeftIcon
        },
        setup() {
            let _message = {
                user_id: "",
                is_customer: "",
                content: ""
            }

            const m = ref(null);

            const messages = ref([]);
            const textareaRef = ref(null);
            const conn = inject('websocketStore').conn;


            // watchEffect(() => {
            //     if (conn === null) {
            //         router.push({
            //             name: 'cus-map'
            //             // name: 'cus-finding-driver'
            //         })
            //     }
            //     // const roomId = conn.url.split('/')[5]
            // })

            const sendMessage = () => {
                if (!textareaRef.value || !textareaRef.value.value) {
                    return;
                }
                if (conn === null) {
                    router.push({
                        name: 'cus-map'
                        // name: 'cus-finding-driver'
                    })
                }

                conn.send(textareaRef.value.value);
                textareaRef.value.value = '';
            }

            watchEffect(() => {
                if (conn === null) {
                    router.push({
                        name: 'cus-map'
                        // name: 'cus-finding-driver'
                    })
                }

                conn.onmessage = (message) => {
                    m.value = JSON.parse(message.data);

                    _message = {
                        user_id: m.value.user_id,
                        is_customer: m.value.is_customer,
                        content: m.value.content
                    }

                    if (m.value.content == 'finishline-@123!(*234kh219871233hadsfh') {
                        conn.close();
                        router.push({
                            name: 'cus-complete-trip'
                        })
                    }

                    messages.value.push(_message);
                    // for (let i = 0; i < messages.value.length; i++) {
                    //     console.log(`test loop test ${messages.value[i].user_id}`)
                    //     console.log(`test loop test ${messages.value[i].is_customer}`)
                    //     console.log(`test loop test ${messages.value[i].content}`)
                    // }
                }
                conn.onclose = () => {}
                conn.onerror = () => {}
                conn.onopen = () => {}
            })
            
            return {                   
                textareaRef,
                messages,
                conn,
                sendMessage
            }
        }
    }

</script>
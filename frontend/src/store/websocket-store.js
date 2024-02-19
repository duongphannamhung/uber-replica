// websocketStore.js
import { reactive } from 'vue';

export const websocketStore = reactive({
  conn: null,
  setConn(c) {
    this.conn = c;
  }
});
import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';

const pinia = createPinia(),
  app = createApp(App);

app.use(pinia);
app.mount('#app');

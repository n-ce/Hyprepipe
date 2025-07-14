<script setup>
import { useData } from '@/stores/player.js';

const data = useData();

defineEmits(['get-artist']);
</script>

<template>
  <div v-if="data.state.title && data.state.artist" class="wrap">
    <h1>{{ data.state.title }}</h1>
    <h3>
      <a
        :href="data.state.artistUrl"
        @click.prevent="$emit('get-artist', data.state.artistUrl)"
        >{{ data.state.artist }}</a
      >
    </h3>
  </div>
</template>

<style scoped>
h1 {
  font-weight: 500;
}

h3 {
  font-size: 1.2rem;
}
.wrap {
  text-align: center;
}
@media (min-width: 1024px) {
  .wrap {
    text-align: left;
  }
}
</style>

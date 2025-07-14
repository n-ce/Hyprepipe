<script setup>
import { useI18n } from '@/stores/misc.js';

const { t } = useI18n();

defineProps({
  name: String,
  author: {
    type: String,
    default: '',
  },
  art: {
    type: String,
    default: '/1x1.png',
  },
});

defineEmits(['open-album']);
</script>

<template>
  <div
    class="album card pop"
    tabindex="0"
    @click="$emit('open-album')"
    @keydown.enter="$emit('open-album')">
    <img class="card-bg bg-img pop-2" :src="art" loading="lazy" alt="" />

    <div class="card-text">
      <h4>{{ name?.replace('subscribers', t('artist.subs')) }}</h4>
      <i>{{ author }}</i>
    </div>
  </div>
</template>

<style scoped>
.card {
  min-height: 17rem;
  width: 15rem;
  border-radius: 0.5rem;
  padding: 1rem;
  margin: auto;
}
.card-bg {
  height: 13rem;
  width: 13rem;
}
.card-text {
  margin-top: 0.5rem;
  cursor: default;
}

body.compact .card {
  min-height: 12rem;
  width: 10rem;
}
body.compact .card-bg {
  width: 7rem;
  height: 7rem;
}
body.compact .card-text {
  margin-top: 0.25rem;
}
</style>

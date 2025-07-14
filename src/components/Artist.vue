<script setup>
import { ref } from 'vue';
import Btn from './Btn.vue';

import { useStore } from '@/scripts/util.js';

import { useResults, useArtist } from '@/stores/results.js';
import { useI18n } from '@/stores/misc.js';

const artist = useArtist(),
  results = useResults(),
  store = useStore(),
  { t } = useI18n();

const subs = JSON.parse(store.subs ? store.subs : '[]'),
  hash = artist.state.hash,
  isSub = ref(subs.includes(hash));

function Sub() {
  if (!artist.state.title) return;

  if (isSub.value) {
    subs.splice(subs.indexOf(hash), 1);
    isSub.value = false;
  } else {
    subs.push(hash);
    isSub.value = true;
  }

  store.setItem('subs', JSON.stringify(subs));
}
</script>

<template>
  <div class="us-wrap">
    <img class="bg-imgfill" :src="artist.state.thumbnails[1].url" />
    <div class="us-main">
      <h2>{{ artist.state.title }}</h2>
      <p @click="$event.target.classList.toggle('more')">
        {{ artist.state.description }}
      </p>
      <div class="us-playwrap">
        <Btn @click="results.getAlbum('?list=' + artist.state.playlistId)" />
        <span
          class="us-box subs"
          :data-active="isSub"
          :data-subs-l="t('artist.subscribers')"
          :data-subs="t('artist.subs')"
          @click="Sub">
          {{ artist.state.subscriberCount || 0 }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.us-wrap {
  min-height: 5rem;
  margin-bottom: 5rem;
}
.bg-imgfill {
  position: absolute;
  width: 100%;
  height: 100%;
  filter: blur(5px) opacity(50%);
}
.us-main {
  color: #fff;
  border-radius: 0.5rem;
  display: flex;
  flex-direction: column;
  padding: 1rem;
  box-shadow: inset 0 0 10rem #000;
  z-index: 1;
}
h2 {
  padding: 1rem;
  font-size: 2rem;
  font-weight: bold;
}
p {
  --line: 3;
  display: -webkit-box;
  -webkit-line-clamp: var(--line);
  line-clamp: var(--line);
  -webkit-box-orient: vertical;
  hyphens: auto;
  margin: 1rem;
  font-size: 1rem;
  line-height: 1.5rem;
  overflow: hidden;
}
p.more {
  --line: auto;
}
.us-playwrap {
  padding: 1rem;
}
.us-box {
  border-radius: 0.25rem;
  font-weight: bold;
  margin-left: 2rem;
  padding: 0.5rem;
  border: 0.125rem solid var(--color-foreground);
  background: var(--color-background-mute);
  color: var(--color-foreground);
  box-shadow: 0 0 1rem var(--color-background-mute);
}
.subs {
  transition:
    background-color 0.4s ease,
    color 0.4s ease;
}
.subs::after {
  content: ' ' attr(data-subs-l);
  font-weight: bold;
}
.subs:hover,
.subs[data-active='true'] {
  background: var(--color-gradient);
  color: var(--color-background);
  background-clip: border-box;
}
@media (max-width: 400px) {
  .subs::after {
    content: ' ' attr(data-subs);
  }
}
</style>

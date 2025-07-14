<script setup>
import { ref, watch } from 'vue';

import { getJsonHyp } from '@/scripts/fetch.js';
import { useData } from '@/stores/player.js';
import { useI18n } from '@/stores/misc.js';

import TextModal from './TextModal.vue';

const { t } = useI18n(),
  data = useData(),
  text = ref(''),
  source = ref(''),
  status = ref(false);

function set(id) {
  getJsonHyp('/lyrics/' + id).then(res => {
    text.value = res.text;
    source.value = res.source;
    status.value = true;
  });
}

function get() {
  status.value = false;

  const now = data.current();

  if (now.lyrics) {
    set(now.lyrics);
  } else if (data.state.url) {
    getJsonHyp(
      '/next/' + data.state.url.replace('/watch?v=', '') + '?queue=avoid',
    ).then(({ lyricsId }) => {
      if (lyricsId) {
        set(lyricsId);
        now.lyrics = lyricsId;
      }
    });
  }
}

get();

watch(
  () => data.state.url,
  () => get(),
);
</script>

<template>
  <TextModal>
    <template #content>
      <pre
        class="placeholder"
        :data-placeholder="
          data.state.urls[0]?.url && !status
            ? t('lyrics.load')
            : t('lyrics.void')
        "
        >{{ text }}</pre
      >
      <div>{{ source }}</div>
    </template>
  </TextModal>
</template>

<style scoped>
.placeholder:empty::before {
  --ico: '\f262';
}
</style>

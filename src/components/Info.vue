<script setup>
import { useI18n } from '@/stores/misc.js';

import TextModal from './TextModal.vue';

defineProps({ text: String });

const { t } = useI18n(),
  parse = d => new DOMParser().parseFromString(d, 'text/html').body.innerText;
</script>

<template>
  <TextModal>
    <template #content>
      <pre class="placeholder" :data-placeholder="t('info.no_info')">{{
        text && parse(text.replaceAll('<br>', '\n'))
      }}</pre>
    </template>
  </TextModal>
</template>

<style scoped>
.placeholder:empty::before {
  --ico: '\F3B9';
}
</style>

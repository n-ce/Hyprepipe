<script setup>
defineProps(['ico']);
defineEmits(['click']);

import { useI18n } from '@/stores/misc.js';
const { t } = useI18n();
</script>

<template>
  <button
    :class="'bi clickable bi-' + (ico ?? 'play')"
    @click="$emit('click')"
    :aria-label="ico ? t(`action.${ico}`) : 'play'"
    :title="ico ? t(`action.${ico}`) : 'play'">
    <slot name="menu"></slot>
  </button>
</template>

<style scoped>
.bi {
  height: 4rem;
  width: 4rem;
  font-size: 4rem;
  color: var(--color-text);
  padding: 0;
  line-height: 0;
  background: var(--color-gradient);
  border-radius: 50%;
  vertical-align: -1rem;
  text-align: center;
  transition: background 0.4s ease;
  margin-right: auto;
}
.bi-play {
  color: var(--color-background);
  box-shadow: 0 0 1rem var(--color-foreground);
  transition: all 0.4s ease;
}
.bi-play::before {
  padding-left: 0.2rem;
}
.bi:hover,
.bi:not(.bi-play) {
  background: transparent;
  box-shadow: none;
}
.bi:not(.bi-play) {
  font-size: 1.75rem;
  width: 3.5rem;
}
</style>

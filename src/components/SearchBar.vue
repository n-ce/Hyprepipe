<script setup>
import { ref, watch } from 'vue';
import { useNav, useI18n } from '@/stores/misc.js';

const { t } = useI18n(),
  show = ref(false),
  nav = useNav(),
  searchEl = ref(null);

function search(e) {
  nav.state.search = e.target.value;
  nav.state.page = 'home';
  e.target.blur();
}

watch(
  () => nav.state.show,
  e => {
    if (e === true) setTimeout(() => searchEl.value.focus(), 0);
  },
);
</script>

<template>
  <button
    class="bi bi-search popup-wrap"
    aria-label="Search"
    @mouseenter="nav.state.show = true"
    @mouseleave="nav.state.show = false"
    @keydown.enter="nav.show()">
    <Transition name="fade">
      <div v-show="nav.state.show" class="popup">
        <input
          type="search"
          ref="searchEl"
          name="searchEl"
          aria-label="Search Input"
          :placeholder="t('title.search') + '...'"
          @change="search"
          @keydown.enter="search"
          :value="nav.state.search" />
      </div>
    </Transition>
  </button>
</template>

<style scoped>
.popup {
  background: transparent;
  right: 0;
  top: -0.5rem;
  bottom: -0.5rem;
  padding: 0;
}
.popup input {
  color: var(--color-text);
  width: calc(100vw - 4rem);
  max-width: 600px;
  font-size: 1rem;
  border: none;
  border-radius: inherit;
  background: var(--color-background-mute);
  text-align: center;
}
</style>

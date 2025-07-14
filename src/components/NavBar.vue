<script setup>
import SearchBar from '@/components/SearchBar.vue';
import IcoHyp from '@/assets/icons/IcoHyp.vue';

import { useNav } from '@/stores/misc.js';

import { useRoute } from '@/scripts/util.js';

const emit = defineEmits(['explore']);

const nav = useNav();

function home() {
  useRoute('/');
  nav.state.page = 'home';
  emit('explore');
}

function set(page) {
  if (page == 'home') useRoute('/');
  else useRoute(`/${page}/`);

  if (page == 'home' && nav.state.page == 'home') emit('explore');

  nav.state.page = page;
}
</script>

<template>
  <nav>
    <h1
      class="bi"
      @click="home"
      @keydown.enter="home"
      tabindex="0"
      aria-label="Home">
      <IcoHyp />
    </h1>

    <div class="wrap">
      <span
        role="link"
        aria-label="Home"
        tabindex="0"
        class="nav-ico bi bi-house clickable"
        :data-active="nav.state.page == 'home'"
        @click="set('home')"
        @keydown.enter="set('home')"></span>
      <span
        role="link"
        aria-label="Explore"
        tabindex="0"
        class="nav-ico bi bi-compass clickable"
        :data-active="nav.state.page == 'explore'"
        @click="set('explore')"
        @keydown.enter="set('explore')"></span>
      <span
        role="link"
        aria-label="Charts"
        tabindex="0"
        class="nav-ico bi bi-graph-up-arrow clickable"
        :data-active="nav.state.page == 'charts'"
        @click="set('charts')"
        @keydown.enter="set('charts')"></span>
      <span
        role="link"
        aria-label="Library"
        tabindex="0"
        class="nav-ico bi bi-collection clickable"
        :data-active="nav.state.page == 'library'"
        @click="set('library')"
        @keydown.enter="set('library')"></span>
      <span
        role="link"
        aria-label="Prefs"
        tabindex="0"
        class="nav-ico bi bi-gear clickable"
        :data-active="nav.state.page == 'prefs'"
        @click="set('prefs')"
        @keydown.enter="set('prefs')"></span>
    </div>

    <div class="wrap">
      <SearchBar />
    </div>
  </nav>
</template>

<style scoped>
nav {
  width: 100%;
  margin-bottom: 2rem;
  display: flex;
  align-items: center;
}
h1.bi {
  display: flex;
  font-size: calc(1.75rem + 1.5vw);
  padding-bottom: 0;
}
h1.bi svg {
  width: calc(1.75rem + 1.5vw);
  height: calc(1.75rem + 1.5vw);
  fill: currentColor;
}
.bi {
  font-size: calc(1rem + 1.5vw);
}
.wrap {
  text-align: center;
  margin-left: auto;
}
[role='menuitem'] {
  cursor: pointer;
}
.nav-ico {
  margin: 0 0.5rem;
}
@media (min-width: 1024px) {
  .bi {
    font-size: 1.5rem;
  }
}
</style>

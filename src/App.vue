<script setup>
/* Imports */
import {
  ref,
  reactive,
  defineAsyncComponent,
  onBeforeMount,
  onMounted,
} from 'vue';

/* Components */
import NavBar from '@/components/NavBar.vue';
import Player from '@/components/Player.vue';
import StatusBar from '@/components/StatusBar.vue';
import NowPlaying from '@/components/NowPlaying.vue';
import Search from '@/components/Search.vue';
import Playlists from '@/components/Playlists.vue';
import Lyrics from '@/components/Lyrics.vue';
import Info from '@/components/Info.vue';
import AddToPlaylist from '@/components/AddToPlaylist.vue';
import Artist from '@/components/Artist.vue';

/* Async Components */
const Genres = defineAsyncComponent(() => import('@/components/Genres.vue')),
  Charts = defineAsyncComponent(() => import('@/components/Charts.vue')),
  Library = defineAsyncComponent(() => import('@/components/Library.vue')),
  Prefs = defineAsyncComponent(() => import('@/components/Prefs.vue')),
  RestorePrefs = defineAsyncComponent(
    () => import('@/components/RestorePrefs.vue'),
  );

/* Composables */
import { useStore, useUnwrap, useAutoTheme } from '@/scripts/util.js';
import { useSetupDB, useUpdatePlaylist } from '@/scripts/db.js';

/* Stores */
import { useData, usePlayer } from '@/stores/player.js';
import { useResults, useArtist } from '@/stores/results.js';
import { useNav, useI18n, useAlert } from '@/stores/misc.js';

const { t, setupLocale } = useI18n(),
  store = useStore(),
  data = useData(),
  player = usePlayer(),
  results = useResults(),
  artist = useArtist(),
  nav = useNav(),
  errs = useAlert();

const genreid = ref(''),
  path = ref(location.pathname);

/* Functions */
function parseUrl() {
  const loc = location.pathname.split('/'),
    base = loc[1].replace(location.search, '');

  path.value = location.pathname;

  switch (base) {
    case '':
      results.getExplore();
      break;
    case 'search':
      nav.state.search = useUnwrap(loc[2]);
      break;
    case 'watch':
      player.state.status = 'circle';
      data.getSong(loc[1] + location.search);
      break;
    case 'playlist':
      results.getAlbum(loc[1] + location.search);
      break;
    case 'channel':
      artist.getArtist(loc[2]);
      break;
    case 'explore':
      genreid.value = loc[2];
      nav.state.page = 'explore';
      break;
    case 'browse':
      if (loc[2].startsWith('MPRE')) results.getAlbum('?list=' + loc[2]);
      break;
    case 'charts':
    case 'library':
    case 'prefs':
    case 'restore':
      nav.state.page = base;
      break;
    default:
      console.log(loc);
      break;
  }
}

/* Keybindings */
function setupKeys() {
  window.addEventListener('keydown', e => {
    if (!e.ctrlKey || e.repeat || 'string' == typeof e.target.value) return;

    switch (e.code) {
      case 'Space':
        player.toggle('play');
        break;
      case 'Slash':
        nav.show();
        break;
      case 'KeyN':
        data.nextTrack();
        break;
      case 'KeyP':
        data.prevTrack();
        break;
    }
  });
}

function playThis(t) {
  const i = data.state.urls.indexOf(t);
  data.play(data.state.urls[i]);
}

function playList(a) {
  data.state.urls = a;
  data.play(data.state.urls[0]);
}

/* Lifestyle hooks */
onBeforeMount(() => {
  /* Set the default theme if set */
  if (store.theme)
    document.body.setAttribute('data-theme', useAutoTheme(store.theme));

  if (store.compact == 'true') document.body.classList.add('compact');

  /* Prefers Reduced Motion */
  if (
    store.prm !== 'false' &&
    (store.prm === 'true' || matchMedia('(prefers-reduced-motion)').matches)
  )
    document.body.classList.add('prm');

  /* Set the default locale if set */
  const loc =
    store.locale ||
    (navigator.language.startsWith('en') ? null : navigator.language);
  if (loc) setupLocale(loc);

  /* Set the default tab */
  if (location.pathname == '/' && store.page) nav.state.page = store.page;
});

onMounted(() => {
  /* Event Listener for change in url */
  window.addEventListener('popstate', parseUrl);

  /* Alert User on close if url is present */
  window.onbeforeunload = () => {
    if (data.state.url) return 'Are you Sure?';
  };

  /* Setup IndexedDB for storing custom playlists */
  useSetupDB();

  setupKeys();
  parseUrl();

  console.log('Mounted <App>!');
});
</script>

<template>
  <NavBar @explore="results.getExplore" />

  <template v-if="artist.state.title && nav.state.page == 'home'">
    <Artist />
  </template>

  <header v-if="!artist.state.title">
    <img v-if="data.state.art" class="art" :src="data.state.art" />

    <div class="wrapper">
      <NowPlaying @get-artist="artist.getArtist" />
    </div>
  </header>

  <main class="placeholder" :data-placeholder="t('info.search')">
    <KeepAlive>
      <Search v-if="nav.state.page == 'home'" @play-urls="playList" />
    </KeepAlive>

    <KeepAlive>
      <Genres v-if="nav.state.page == 'explore'" :id="genreid" />
    </KeepAlive>

    <KeepAlive>
      <Charts v-if="nav.state.page == 'charts'" @play-urls="playList" />
    </KeepAlive>

    <Library
      v-if="nav.state.page == 'library'"
      @play-urls="playList"
      @open-playlist="results.getAlbum" />

    <Prefs v-if="nav.state.page == 'prefs'" />

    <RestorePrefs v-if="nav.state.page == 'restore'" />
  </main>

  <Transition name="fade">
    <Playlists v-if="player.state.playlist" @playthis="playThis" />
  </Transition>

  <Transition name="fade">
    <Lyrics v-if="player.state.lyrics" />
  </Transition>

  <Transition name="fade">
    <Info v-if="player.state.info" :text="data.state.description" />
  </Transition>

  <Transition name="fade">
    <div v-if="errs.msg" class="alert">
      {{ errs.msg }}
    </div>
  </Transition>

  <AddToPlaylist />

  <StatusBar />

  <Player />
</template>

<style>
@import './assets/base.css';
@import 'bootstrap-icons/font/bootstrap-icons.css';

#app {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem;
  font-weight: normal;
  margin-bottom: 10rem;
}
main {
  display: flex;
  flex-direction: column;
}
main:empty::before {
  font-size: 3.5rem !important;
}
header {
  line-height: 1.5;
  padding-bottom: 2.5rem;
}
.art {
  display: block;
  margin: 0 auto 2rem;
  height: 175px;
  width: 175px;
}
img,
.card,
.card-bg {
  border-radius: 0.25rem;
  border-radius: 0.25rem;
}
h1 {
  text-align: center;
  padding-bottom: 1.5rem;
}
h1,
h4 {
  font-weight: bold;
}
a,
.green {
  text-decoration: none;
  color: var(--color-foreground);
  transition: 0.4s;
}
a:focus {
  outline: none;
}
a:focus-visible {
  border-radius: 0.25rem;
  outline: 0.2rem solid var(--color-border);
  transition: none;
}
.flex {
  display: flex;
  align-items: center;
  justify-content: center;
}
.center {
  margin-left: auto;
  margin-right: auto;
}
.flex .bi {
  line-height: 0;
}
.alert {
  position: fixed;
  font-weight: bold;
  letter-spacing: 0.1rem;
  max-width: 20rem;
  right: 1rem;
  bottom: 10rem;
  padding: 1rem;
  background: linear-gradient(135deg, indianred, #bf616a);
  color: var(--color-background);
  border-radius: 0.25rem;
  box-shadow: 0.3rem 0.3rem 0.4rem indianred;
}

@media (hover: hover) {
  a:hover {
    background-color: var(--color-border);
  }
}

@media (min-width: 1024px) {
  header {
    margin: auto;
    display: flex;
    place-items: center;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
  .art {
    margin: 0 2rem 0 0;
  }
}
</style>

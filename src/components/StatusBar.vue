<script setup>
import { ref, reactive } from 'vue';

import { useStore } from '@/scripts/util.js';
import {
  getAuthPlaylists,
  useAuthCreatePlaylist,
  useAuthAddToPlaylist,
} from '@/scripts/fetch.js';

import { useData, usePlayer } from '@/stores/player.js';
import { useI18n, useAlert } from '@/stores/misc.js';

const { t } = useI18n(),
  data = useData(),
  player = usePlayer(),
  store = useStore(),
  a = useAlert();

const showme = reactive({
    menu: false,
    vol: false,
  }),
  liked = ref(undefined),
  liking = ref(false);

function getFormattedTime(sec) {
  let seconds = parseInt(sec);
  let minutes = parseInt(seconds / 60);
  seconds -= minutes * 60;
  const hours = parseInt(minutes / 60);
  minutes -= hours * 60;

  if (hours === 0) return `${minutes}:${String(seconds % 60).padStart(2, '0')}`;
  else
    return `${hours}:${String(minutes).padStart(2, '0')}:${String(
      seconds % 60,
    ).padStart(2, '0')}`;
}

async function Offline() {
  if (window.offline && data.state.url) {
    window.offline
      .store(window.audioPlayer.getAssetUri(), {
        title: data.state.title,
        url: data.state.url,
        artist: data.state.artist,
        artistUrl: data.state.artistUrl,
      })
      .promise.catch(e => {
        console.error(e);
        a.add('Error: ' + e.code);
      });
  } else a.add('offline storage not found');
}

async function Like() {
  liking.value = true;

  const remote = await getAuthPlaylists();

  let fav = remote.find(i => i.name == 'Playlist - Favorites');

  if (!fav) {
    const { playlistId } = await useAuthCreatePlaylist('Favorites');

    fav = { id: playlistId };
  }

  const { message } = await useAuthAddToPlaylist(fav.id, data.state.url);

  if (message == 'ok') liked.value = data.state.url;

  liking.value = false;
}
</script>
<template>
  <div id="statusbar" class="flex">
    <div class="flex statusbar-progress-container">
      <span :aria-label="t('info.time_elapsed')">
        {{ getFormattedTime(player.state.realTime) }}
      </span>
      <div
        id="statusbar-progress"
        class="range-wrap"
        :aria-label="t('action.change_time')">
        <input
          :aria-label="t('action.change_time')"
          :value="player.state.time"
          type="range"
          name="statusbar-progress"
          max="100"
          @input="
            player.state.currentTime =
              ($event.target.value / 100) * player.state.duration
          " />
      </div>
      <span :aria-label="t('info.track_duration')">
        {{ getFormattedTime(player.state.duration) }}
      </span>
    </div>

    <div class="flex statusbar-controls-container">
      <div class="statusbar-control-item left">
        <button
          id="list-btn"
          :title="t('statusBar.currentPlaylist')"
          :aria-label="t('statusBar.currentPlaylist')"
          class="bi bi-music-note-list clickable"
          :data-active="player.state.playlist"
          @click="player.toggle('playlist')"></button>

        <button
          id="loop-btn"
          :title="t('statusBar.loop')"
          :aria-label="t('statusBar.loop')"
          class="bi clickable"
          :class="player.state.loop < 2 ? 'bi-repeat' : 'bi-repeat-1'"
          :data-active="player.state.loop > 0"
          @click="
            player.state.loop < 2
              ? player.state.loop++
              : (player.state.loop = 0)
          "></button>
      </div>
      <div class="statusbar-control-item center">
        <button
          :class="{
            hidden:
              data.state.urls[
                data.state.urls.findIndex(s => s.url == data.state.url) - 1
              ] == undefined &&
              (player.state.loop != 1 || data.state.urls.length == 0),
          }"
          id="btn-previoustrack"
          aria-label="Play previous track"
          class="bi bi-chevron-bar-left clickable"
          @click="data.prevTrack"></button>

        <button
          :class="{ hidden: player.state.duration == 0 }"
          id="btn-rewind"
          aria-label="Rewind 10 seconds"
          class="bi bi-chevron-left clickable"
          @click="
            player.state.currentTime = Math.max(
              0,
              player.state.currentTime - 10,
            )
          "></button>

        <button
          id="btn-play-pause"
          aria-label="Play or Pause"
          class="bi clickable"
          :class="'bi-' + player.state.status"
          @click="player.toggle('play')"></button>

        <button
          :class="{ hidden: player.state.duration == 0 }"
          id="btn-forward"
          aria-label="Forward 10 seconds"
          class="bi bi-chevron-right clickable"
          @click="
            player.state.currentTime = Math.min(
              player.state.duration,
              player.state.currentTime + 10,
            )
          "></button>

        <button
          :class="{
            hidden:
              data.state.urls[
                data.state.urls.findIndex(s => s.url == data.state.url) + 1
              ] == undefined &&
              (player.state.loop != 1 || data.state.urls.length == 0),
          }"
          id="btn-nexttrack"
          :aria-label="t('action.play_next_track')"
          class="bi bi-chevron-bar-right clickable"
          @click="data.nextTrack"></button>
      </div>
      <div class="statusbar-control-item right">
        <button
          id="vol-btn"
          aria-label="Volume Buttons"
          @click="
            showme.vol = !showme.vol;
            showme.menu = false;
          "
          class="popup-wrap bi bi-volume-up clickable">
          <Transition name="fade">
            <div v-if="showme.vol" id="vol" class="popup">
              <input
                id="vol-input"
                aria-label="Volume Input"
                type="range"
                :value="player.state.vol"
                max="1"
                step=".01"
                @input="player.state.vol = $event.target.value" />
            </div>
          </Transition>
        </button>

        <button
          class="bi bi-three-dots clickable"
          :aria-label="t('action.more_controls')"
          @click="
            showme.menu = !showme.menu;
            showme.vol = false;
          "></button>
        <Transition name="fade">
          <div id="menu" v-if="showme.menu" class="popup">
            <button
              id="info-btn"
              class="bi bi-info-circle clickable"
              :title="t('action.show_information_about_track')"
              :aria-label="t('action.show_information_about_track')"
              :data-active="player.state.info"
              @click="player.toggle('info')"></button>

            <button
              v-if="store.auth"
              id="like-btn"
              :title="t('action.add_track_to_favorites')"
              :aria-label="t('action.add_track_to_favorites')"
              class="bi blink clickable"
              :class="data.state.url == liked ? 'bi-heart-fill' : 'bi-heart'"
              :data-active="data.state.url == liked"
              :data-loading="liking"
              @click="Like"></button>

            <button
              id="dl-btn"
              :title="t('action.save_for_offline_playback')"
              :aria-label="t('action.save_for_offline_playback')"
              class="bi bi-download clickable"
              @click="Offline"></button>

            <button
              id="addToPlaylist"
              :title="t('playlist.add')"
              :aria-label="t('playlist.add')"
              class="bi bi-collection clickable"
              @click="
                player.toggle('add');
                showme.menu = false;
              "></button>

            <button
              id="btn-lyrics"
              :title="t('lyrics.lyrics')"
              :aria-label="t('lyrics.lyrics')"
              class="bi bi-chat-square-quote clickable"
              :data-active="player.state.lyrics"
              @click="player.toggle('lyrics')"></button>
          </div>
        </Transition>
      </div>
    </div>
  </div>
</template>

<style scoped>
#statusbar {
  flex-direction: column;
  row-gap: 0.5rem;
  position: fixed;
  left: 0;
  bottom: 0;
  right: 0;
  padding: 1.5rem 0;
  border-top: 0.25rem solid var(--color-foreground);
  background: var(--color-background);
  min-height: 15vh;
  z-index: 999;
}

.statusbar-progress-container,
.statusbar-controls-container {
  width: 80vw;
  min-width: 200px;
  max-width: 600px;
}

/* Progress */
.statusbar-progress-container {
  column-gap: 0.5rem;
}
.range-wrap {
  --fw: 0%;
  display: flex;
  align-items: center;
  position: relative;
  height: 0.25rem;
  display: inline-block;
}
.range-wrap::before {
  content: '';
  width: var(--fw);
  position: absolute;
  left: 0;
  bottom: 0;
  height: 0.25rem;
  background-color: var(--color-foreground);
  transition: width 0.4s ease;
  z-index: 2;
}
.range-wrap input[type='range'] {
  --w: 100%;
  --h: 100%;
  position: absolute;
  top: 0;
}
#statusbar-progress {
  --fw: v-bind('player.state.time + "%"');
  width: 100%;
}

/* Controls */
.statusbar-controls-container {
  display: grid;
  grid-template-columns: 30% 40% 30%;
}
.statusbar-control-item {
  display: flex;
}
.statusbar-control-item.left {
  justify-content: flex-start;
}
.statusbar-control-item.center {
  justify-content: center;
}
.statusbar-control-item.right {
  justify-content: flex-end;
}
.bi-play,
.bi-pause,
.bi-circle {
  font-size: 2.5rem !important;
}
.bi-circle {
  animation: blink infinite ease 1s;
}
.bi-volume-up,
.bi-repeat,
.bi-repeat-1 {
  font-size: 1.5rem !important;
}
#vol {
  --h: 6.5rem;
  --w: 1rem;
  display: flex;
  box-shadow: -0.5rem -0.5rem 2rem var(--color-shadow);
  transform: rotateZ(270deg) translateX(calc(calc(var(--h) / 2) - 0.5rem))
    translateY(calc(calc(var(--w) + 2rem) * -1));
}
#menu {
  left: 50%;
  box-shadow: 0.5rem 0.5rem 2rem var(--color-shadow);
}

/* Range */
input[type='range'] {
  width: var(--h);
  height: var(--w);
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  border: none;
  border-radius: 5rem;
  background: transparent;
  color: transparent;
  cursor: pointer;
  position: relative;
}
input[type='range']:hover,
input[type='range']:focus {
  outline-style: none;
}
input[type='range']:focus-visible {
  outline-style: solid;
}

/* Range Webkit */
input[type='range']::-webkit-slider-thumb {
  background-color: var(--color-foreground);
  -webkit-appearance: none;
  appearance: none;
  opacity: 1;
  height: 1rem;
  width: 1rem;
  border-radius: 50%;
  margin-top: -0.4rem;
}
input[type='range']::-webkit-slider-runnable-track {
  height: 100%;
  background-color: var(--color-border);
}
#vol input[type='range']::-webkit-slider-runnable-track {
  height: 0.1rem;
}

/* Range Moz */
input[type='range']::-moz-range-thumb {
  -moz-appearance: none;
  appearance: none;
  border: none;
  outline: none;
  background-color: var(--color-foreground);
  opacity: 1;
  height: 1rem;
  width: 1rem;
  border-radius: 50%;
  margin-top: -0.45rem;
}
input[type='range']::-moz-range-track {
  height: 100%;
  background-color: var(--color-border);
}
#vol input[type='range']::-moz-range-track {
  height: 0.1rem;
}

@media (max-width: 500px) {
  .statusbar-progress-container {
    min-width: initial;
  }
  #menu {
    left: auto;
    right: -0.5rem;
  }
  #btn-rewind,
  #btn-forward {
    display: none;
  }
}
</style>

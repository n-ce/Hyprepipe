<script setup>
import { ref, watch, onMounted } from 'vue';

import Modal from './Modal.vue';

import { getAuthPlaylists, useAuthAddToPlaylist } from '@/scripts/fetch.js';
import { useListPlaylists, useUpdatePlaylist } from '@/scripts/db.js';
import { useStore } from '@/scripts/util.js';

import { useData, usePlayer } from '@/stores/player.js';
import { useI18n } from '@/stores/misc.js';

const { t } = useI18n(),
  data = useData(),
  player = usePlayer(),
  store = useStore();

const props = defineProps({
    show: Boolean,
    song: String,
    title: String,
    artist: String,
  }),
  emit = defineEmits(['show']);

const pl = ref(''),
  list = ref([]),
  remote = ref([]),
  plRemote = ref(false);

const url = () => props.song || data.state.url,
  title = () => props.title || data.state.title,
  artist = () => props.artist || data.state.artist,
  show = {
    get is() {
      return props.song ? true : player.state.add;
    },
    set is(e) {
      props.song ? emit('show', e) : (player.state.add = e);
    },
  };

function Save() {
  if (pl.value) {
    if (plRemote.value == true && store.auth) {
      useAuthAddToPlaylist(pl.value, url());
    } else if (plRemote.value == false) {
      useUpdatePlaylist(
        pl.value,
        { url: url(), title: title(), artist: artist() },
        e => {
          if (e === true) console.log('Added Song');
        },
      );
    }
  }
}

function List() {
  show.is = true;
  useListPlaylists(res => {
    list.value = res;
  });
  getAuthPlaylists().then(res => {
    remote.value = res;
  });
}

watch(
  () => (props.song ? false : player.state.add),
  e => e === true && List(),
);

onMounted(() => {
  props.song && List();
});
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <Modal
        n="2"
        :display="show.is"
        :title="t('playlist.select')"
        @show="
          e => {
            show.is = e;
          }
        ">
        <template #content>
          <div
            v-for="i in list"
            :key="i.name"
            class="flex item"
            @click="
              pl = i.name;
              plRemote = false;
            "
            :data-active="pl == i.name && plRemote == false">
            <span>{{ i.name }}</span
            ><span class="ml-auto">{{ i.urls.length || '' }}</span>
          </div>
          <div
            v-for="i in remote"
            :key="i.id"
            class="flex item remote"
            @click="
              pl = i.id;
              plRemote = true;
            "
            :data-active="pl == i.id && plRemote == true">
            <span>{{ i.name.replace('Playlist -', '') }}</span
            ><span class="ml-auto">{{ i.videos }}</span>
          </div>
        </template>
        <template #buttons>
          <button aria-label="Cancel" @click="show.is = false">
            {{ t('action.cancel') }}
          </button>

          <button
            aria-label="Add Song"
            @click="
              Save();
              show.is = false;
            ">
            {{ t('action.add') }}
          </button>
        </template>
      </Modal>
    </Transition>
  </Teleport>
</template>

<style scoped>
.ml-auto {
  margin-left: auto;
  width: min-content;
}
.item {
  background: var(--color-background);
  border-radius: 0.5rem;
  margin: 0.5rem 0;
  transition: background-color 0.1s ease;
}
.item:hover {
  background: var(--color-background-mute);
}
.item:active {
  background: var(--color-border);
}
.item[data-active='true'] {
  color: var(--color-background);
  background: linear-gradient(135deg, cornflowerblue, #88c0d0);
}
.remote.item::before {
  content: '\F4E1';
  font-family: bootstrap-icons;
  font-size: 1.25rem;
}
</style>

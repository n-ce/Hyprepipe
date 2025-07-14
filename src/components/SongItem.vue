<script setup>
import { ref } from 'vue';

import AddToPlaylist from '@/components/AddToPlaylist.vue';

import { getJsonAuth } from '@/scripts/fetch.js';
import { useStore, useShare } from '@/scripts/util.js';
import { useUpdatePlaylist } from '@/scripts/db.js';

import { useResults, useArtist } from '@/stores/results.js';
import { useData, usePlayer } from '@/stores/player.js';
import { useI18n } from '@/stores/misc.js';

const data = useData(),
  results = useResults(),
  player = usePlayer(),
  artist = useArtist(),
  { t } = useI18n();

const props = defineProps({
    index: Number,
    author: String,
    title: String,
    channel: String,
    play: String,
    art: String,
    playlistId: String,
    offlineUri: String,
  }),
  emit = defineEmits(['open-song', 'nxt-song', 'remove']);

const show = ref(false),
  showPl = ref(false);

const openSong = el => {
    if (!el.classList.contains('ign')) emit('open-song', props.play);
  },
  addSong = () => {
    data.state.urls.push({
      url: props.play,
      title: props.title,
      thumbnails: [{ url: props.art }],
    });

    const index = data.state.urls.findIndex(s => s.url == data.state.url);

    if (
      (index == data.state.urls.length - 1 && player.state.time > 98) ||
      data.state.urls.length == 1
    )
      emit('open-song', props.play);
  },
  appendSong = () => {
    const index = data.state.urls.findIndex(s => s.url == data.state.url);

    data.state.urls.splice(index + 1, 0, {
      url: props.play,
      title: props.title,
      thumbnails: [{ url: props.art }],
    });
  },
  rm = () => emit('remove', props.index),
  Remove = () => {
    const auth = useStore().getItem('auth'),
      isRemote = results.items?.songs?.title?.startsWith('Playlist - ');

    if (!confirm('Are you sure?')) return;

    if (isRemote && auth) {
      getJsonAuth('/user/playlists/remove', {
        method: 'POST',
        headers: {
          Authorization: auth,
        },
        body: JSON.stringify({
          index: props.index,
          playlistId: props.playlistId,
        }),
      }).then(json => {
        console.log(json);

        if (!json.error) {
          if (json.message == 'ok') rm();
        } else alert(json.error);
      });
    } else if (props.offlineUri)
      window.offline && window.offline.remove(props.offlineUri).then(rm);
    else useUpdatePlaylist(props.playlistId, props.index, () => rm());
  },
  Share = () => {
    const data = {
      title: `Listen to ${props.title} by ${props.author} on Hyperpipe`,
      url: location.origin + props.play,
    };

    useShare(data);
  };
</script>
<template>
  <AddToPlaylist
    v-if="showPl"
    :song="play"
    :title="title"
    :artist="author"
    @show="e => (showPl = e)" />

  <div
    class="song card flex pop"
    tabindex="0"
    @click="openSong($event.target)"
    @keydown.enter="openSong($event.target)">
    <img
      class="pop-2 bg-img song-bg"
      :loading="index > 4 ? 'lazy' : 'eager'"
      :src="art"
      alt />

    <span class="flex content">
      <h4>
        <template v-if="results.items?.songs?.title"
          >{{ index + 1 }}.
        </template>
        {{ title }}
      </h4>
      <a
        class="ign"
        :href="channel"
        @click.prevent="artist.getArtist(channel.replace('/channel/', ''))">
        <i class="ign">{{ author ? author.replaceAll(' - Topic', '') : '' }}</i>
      </a>
      <span v-if="play === data.state?.url">
        <span class="bi-play" aria-hidden="true"></span>
        <span>{{ t('title.now_playing') }}</span>
      </span>
    </span>

    <button
      class="bi bi-three-dots-vertical popup-wrap ign"
      @mouseenter="show = true"
      @mouseleave="show = false"
      @keyup.enter="show = !show"
      :aria-label="t('action.more_controls')">
      <Transition name="fade">
        <div v-if="show" class="popup ign">
          <button
            v-if="playlistId || offlineUri"
            class="bi bi-dash-lg clickable ign"
            :aria-label="t('playlist.remove_from')"
            @click="Remove"></button>
          <button
            class="bi bi-chevron-bar-right clickable ign"
            :aria-label="t('queue.add_to_top')"
            @click="appendSong"></button>
          <button
            class="bi bi-collection clickable ign"
            :aria-label="t('playlist.add')"
            @click="showPl = true"></button>
          <button
            class="bi bi-broadcast clickable ign"
            :aria-label="t('action.play_now')"
            @click="$emit('nxt-song')"></button>
          <button
            class="bi bi-plus-lg clickable ign"
            :aria-label="t('queue.add_to_bottom')"
            @click="addSong"></button>
          <button
            class="bi bi-share clickable ign"
            :aria-label="t('action.share')"
            @click="Share"></button>
        </div>
      </Transition>
    </button>
  </div>
</template>

<style scoped>
.card {
  margin: 1rem 0;
  justify-content: initial;
  cursor: pointer;
}
.song h4 {
  overflow-wrap: anywhere;
}
span.content {
  padding: 1rem;
  flex-direction: column;
  align-items: flex-start;
  flex-grow: 1;
}
span.bi-three-dots-vertical {
  margin: 2rem;
}
button.bi-three-dots-vertical {
  border: none;
}
.popup {
  line-height: auto;
  height: auto;
  right: 0;
  bottom: 0;
  padding: 0.75rem 0.5rem;
  box-shadow: var(--shadow);
  border-radius: 0.25rem;
}
.popup span {
  padding: 0.5rem;
}
.song-bg {
  width: 120px;
  height: 120px;
}
.bi-dash-lg {
  color: indianred;
}

body.compact .card {
  margin: 0;
  padding: 0 0.5rem;
}
body.compact .song-bg {
  width: 70px;
  height: 70px;
}
body.compact .content :is(h4, a) {
  overflow: hidden;
  line-clamp: 1;
  -webkit-line-clamp: 1;
  display: -webkit-box;
  hyphens: auto;
  -webkit-box-orient: vertical;
}

body.prm .card {
  margin: 0;
  padding: 0.5rem;
}
</style>

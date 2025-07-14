<script setup>
import {
  ref,
  watch,
  computed,
  onActivated,
  onUpdated,
  onDeactivated,
} from 'vue';

import Btn from './Btn.vue';
import SongItem from './SongItem.vue';
import AlbumItem from './AlbumItem.vue';

import {
  getJsonPiped,
  getPipedQuery,
  useAuthRemovePlaylist,
} from '@/scripts/fetch.js';
import {
  useVerifyAuth,
  useRoute,
  useWrap,
  useShare,
  useShuffle,
} from '@/scripts/util.js';
import { useCreatePlaylist, useRemovePlaylist } from '@/scripts/db.js';

import { useResults, useArtist } from '@/stores/results.js';
import { useData } from '@/stores/player.js';
import { useNav, useI18n } from '@/stores/misc.js';

const { t } = useI18n(),
  results = useResults(),
  data = useData(),
  nav = useNav(),
  artist = useArtist();

const emit = defineEmits(['play-urls']),
  filters = ['music_songs', 'music_albums', 'music_artists', 'music_playlists'],
  filter = ref('music_songs'),
  isSearch = ref(/search/.test(location.pathname)),
  albumMenu = ref(false);

const plId = computed(
  (x = results.items?.songs?.items?.[0]) => x?.playlistId || !!x?.offlineUri,
);

const shuffleAdd = () => {
    const songs = results.items.songs.items.map(i => ({
      url: i.url,
      title: i.title,
      thumbnails: [{ url: i.thumbnail }],
      thumbnail: i.thumbnail,
      offlineUri: i.offlineUri,
      duration: i.duration,
    }));

    emit('play-urls', useShuffle(songs));
  },
  openSong = (song, nxt = false) => {
    if (results.items?.songs?.title && !nxt) {
      data.state.urls = results.items.songs.items.map(i => ({
        url: i.url || '/watch?v=' + song.id,
        title: i.title,
        thumbnails: [{ url: i.thumbnail }],
        thumbnail: i.thumbnail,
        offlineUri: i.offlineUri,
        duration: i.duration,
      }));

      song.url = song.url || '/watch?v=' + song.id;
      data.play(song);
    } else {
      emit('play-urls', [
        {
          url: song.url || '/watch?v=' + song.id,
          title: song.title || song.name,
          thumbnails: [
            {
              url:
                song.thumbnail ||
                song.thumbnails[1]?.url ||
                song.thumbnails[0]?.url,
            },
          ],
        },
      ]);
    }
  },
  removeSong = i => {
    results.items.songs.items.splice(i, 1);
  },
  shareAlbum = () => {
    const data = {
      title: `View ${results.items?.songs?.title} on Hyperpipe`,
      url:
        location.origin +
        (results.album.startsWith('/') ? '' : '/') +
        results.album,
    };

    useShare(data);
  },
  saveAlbum = () => {
    const urls = results.items?.songs?.items?.map(item => ({
      url: item.url,
      title: item.title,
      artist: item.uploaderName || item.artist || item.subtitle,
    }));

    let title = results.items?.songs?.title;

    if (title) {
      if (title == 'Songs')
        title += ' - ' + results.items.songs.items[0].uploaderName;

      useCreatePlaylist(title, urls, () => {
        alert(t('info.saved'));
      });
    }
  },
  removePlaylist = async () => {
    const consent = confirm('Confirm?');

    if (!plId.value || !consent) return;

    if (plId.value === true)
      window.offline &&
        Promise.all(
          (await window.offline.list()).map(i =>
            window.offline.remove(i.offlineUri),
          ),
        );
    else if (useVerifyAuth(plId.value)) {
      const { message } = await useAuthRemovePlaylist(plId.value);
      if (message != 'ok') {
        alert(message);
        return;
      }
    } else useRemovePlaylist(plId.value);

    useRoute('/library');
    nav.state.page = 'library';
  },
  getSearch = q => {
    if (q) {
      const pq = useWrap(q);

      history.pushState({}, '', `/search/${pq + getPipedQuery()}`);

      document.title = 'Search Results for ' + q;
      isSearch.value = /search/.test(location.pathname);

      getResults(pq);
    } else {
      results.resetItems();

      useRoute('/');
      document.title = 'Hyperpipe';

      console.log('No Search');
    }
  },
  loading = ref(false),
  getNextPage = async () => {
    if (
      (!isSearch.value && !results.items?.songs?.title) ||
      loading.value ||
      results.next == 'null' ||
      !results.next
    )
      return;

    if (
      window.innerHeight + window.scrollY >=
      document.body.offsetHeight - window.innerHeight
    ) {
      loading.value = true;

      let items, key;

      if (isSearch.value) {
        const f = filter.value || 'music_songs';

        const json = await getJsonPiped(
          `/nextpage/search?nextpage=${encodeURIComponent(results.next)}&q=${
            nav.state.search
          }&filter=${f}`,
        );

        key = f.split('_')[1];
        results.next = json.nextpage;

        items = json.items;
      } else {
        const json = await getJsonPiped(`/nextpage/playlists/${results.next}`);
        key = 'songs';

        items = json.relatedStreams;
      }

      results.items[key].items.push(...items);

      loading.value = false;
    }
  },
  getResults = async q => {
    results.resetItems();

    const f = filter.value || 'music_songs',
      json = await getJsonPiped(`/search?q=${q}&filter=${f}`),
      key = f.split('_')[1];

    results.next = json.nextpage;

    results.setItem(key, json);
  };

watch(
  () => nav.state.search,
  n => {
    if (n) {
      n = n.replace(location.search || '', '');

      artist.reset();
      getSearch(n);
    }
  },
);

onUpdated(() => {
  isSearch.value = /search/.test(location.pathname);
});

onActivated(() => {
  window.addEventListener('scroll', getNextPage);
});

onDeactivated(() => {
  window.removeEventListener('scroll', getNextPage);
});
</script>

<template>
  <div
    v-if="results.items.songs && results.items.songs.corrected"
    class="text-full">
    Did you mean, "<span class="caps">{{
      results.items?.songs?.suggestion
    }}</span
    >"!!
  </div>

  <div v-if="results.items?.songs?.title" class="text-full flex">
    <Btn
      @click="
        $emit(
          'play-urls',
          results.items?.songs?.items?.map(item => ({
            url: item.url,
            title: item.title,
            thumbnails: [{ url: item.thumbnail }],
            thumbnail: item.thumbnail,
            offlineUri: item.offlineUri,
            duration: item.duration,
          })),
        )
      " />

    <Btn ico="three-dots" @click="albumMenu = !albumMenu">
      <template #menu>
        <Transition name="fade">
          <div v-if="albumMenu" class="alb popup">
            <button
              class="bi bi-bookmark-plus clickable"
              @click="saveAlbum"
              :aria-label="t('playlist.save')"></button>

            <button
              class="bi bi-share clickable"
              @click="shareAlbum"
              :aria-label="t('action.share')"></button>

            <button
              class="bi bi-plus-lg clickable"
              @click="
                data.state.urls.push(
                  ...results.items.songs.items.map(i => ({
                    url: i.url,
                    title: i.title,
                    thumbnails: [{ url: i.thumbnail }],
                  })),
                )
              "
              :aria-label="t('queue.add')"></button>

            <button
              class="bi bi-shuffle clickable"
              @click="shuffleAdd"
              :aria-label="t('action.shuffle')"></button>

            <button
              v-if="plId"
              class="bi bi-trash3 clickable"
              @click="removePlaylist"
              :aria-label="t('playlist.remove')"></button>
          </div>
        </Transition>
      </template>
    </Btn>

    <span class="ml-auto">{{ results.items?.songs?.title }}</span>
  </div>

  <div v-if="isSearch" class="filters">
    <button
      v-for="f in filters"
      class="filter caps"
      :key="f"
      :data-active="f == filter"
      @click="
        filter = f;
        getSearch(nav.state.search);
      ">
      {{ t('title.' + f.split('_')[1]) }}
    </button>
  </div>

  <div
    v-if="results.items.songs && results.items.songs.items[0]"
    class="search-wrap">
    <h2 v-if="!isSearch">{{ t('title.songs') }}</h2>
    <div class="grid">
      <SongItem
        v-for="(song, index) in results.items.songs.items"
        :key="song.url || song.id"
        :index="index"
        :playlistId="song.playlistId"
        :offlineUri="song.offlineUri"
        :author="song.uploaderName || song.artist || song.subtitle"
        :title="song.title || song.name"
        :channel="
          song.uploaderUrl ||
          song.artistUrl ||
          (song.subId && '/channel/' + song.subId)
        "
        :play="song.url || '/watch?v=' + song.id"
        :art="
          song.thumbnail ||
          song.thumbnails?.[1]?.url ||
          song.thumbnails?.[0]?.url ||
          '/1x1.png'
        "
        @remove="removeSong"
        @open-song="openSong(song)"
        @nxt-song="openSong(song, true)" />
    </div>
    <a
      v-if="artist.state.playlistId"
      @click.prevent="results.getAlbum('?list=' + artist.state.playlistId)"
      class="more"
      :href="'/playlist?list=' + artist.state.playlistId"
      >{{ t('info.see_all') }}</a
    >
  </div>

  <div
    v-if="results.items.albums && results.items.albums.items[0]"
    class="search-wrap">
    <h2 v-if="!isSearch">{{ t('title.albums') }}</h2>
    <div class="grid-3">
      <AlbumItem
        v-for="album in results.items.albums.items"
        :key="album.url || album.id"
        :author="album.uploaderName || album.subtitle"
        :name="album.name || album.title"
        :art="album.thumbnail || album.thumbnails[0].url"
        @open-album="results.getAlbum(album.url || '?list=' + album.id)" />
    </div>
    <a
      v-if="results.items.albums.more?.params"
      @click.prevent="artist.getArtistNext('albums', results.items.albums.more)"
      class="more"
      >{{ t('info.see_all') }}</a
    >
  </div>

  <div
    v-if="results.items.playlists && results.items.playlists.items[0]"
    class="search-wrap">
    <h2>{{ t('title.playlists') }}</h2>
    <div class="grid-3">
      <AlbumItem
        v-for="pl in results.items.playlists.items"
        :key="pl.url"
        :author="pl.videos + ' Songs • ' + pl.uploaderName"
        :name="pl.name"
        :art="pl.thumbnail"
        @open-album="results.getAlbum(pl.url)" />
    </div>
  </div>

  <div
    v-if="results.items.singles && results.items.singles.items[0]"
    class="search-wrap">
    <h2>{{ t('title.singles') }}</h2>
    <div class="grid-3">
      <AlbumItem
        v-for="single in results.items.singles.items"
        :key="single.id"
        :author="single.subtitle"
        :name="single.title"
        :art="single.thumbnails[0].url"
        @open-album="results.getAlbum('?list=' + single.id)" />
    </div>
    <a
      v-if="results.items.singles.more?.params"
      @click.prevent="
        artist.getArtistNext('singles', results.items.singles.more)
      "
      class="more"
      >{{ t('info.see_all') }}</a
    >
  </div>

  <div
    v-if="
      (results.items.recommendedArtists &&
        results.items.recommendedArtists.items[0]) ||
      (results.items.artists && results.items.artists.items[0])
    "
    class="search-wrap">
    <h2 v-if="!isSearch">
      {{
        results.items.artists ? t('title.artists') : t('title.similar_artists')
      }}
    </h2>
    <div class="grid-3 circle">
      <AlbumItem
        v-for="a in results.items.artists
          ? results.items.artists.items
          : results.items.recommendedArtists.items"
        :key="a.id || a.url"
        :author="a.subtitle?.replace('subscribers', t('artist.subscribers'))"
        :name="a.name || a.title"
        :art="a.thumbnail || a.thumbnails[0].url"
        @open-album="
          artist.getArtist(a.id || a.url.replace('/channel/', ''))
        " />
    </div>
  </div>
</template>

<style scoped>
.search-wrap {
  margin-bottom: 2rem;
}
.search-wrap h2,
.circle {
  text-align: center;
}
:deep(.bi) {
  margin-right: 0;
}
:deep(.bi-play) {
  margin-right: 0.75rem;
}
:deep(.ml-auto) {
  margin-left: auto;
}
.alb {
  bottom: 2.5rem;
  border-radius: 0.25rem;
  box-shadow: -0.5rem 0.5rem 1rem var(--color-shadow);
}
.alb .bi {
  padding: 0.5rem;
}
.filters {
  max-width: 100%;
  width: max-content;
  margin: 0 auto 2rem;
  overflow-x: auto;
  white-space: nowrap;
}
.filter {
  max-width: 200px;
  margin: 0 0.25rem;
  padding: 0.5rem 1rem;
  font-size: 1.25rem;
  border-radius: 0.25rem;
}
.filter:hover {
  background: var(--color-background-mute);
}
.filter[data-active='true'] {
  border-bottom: 0.125rem var(--color-text) solid;
  border-radius: 0.25rem 0.25rem 0 0;
}
.text-full {
  padding: 1rem;
  font-size: 1.5rem;
  text-align: center;
  margin-bottom: 1rem;
}
.text-right {
  text-align: right;
}
.more {
  margin: 1.5rem 0.5rem;
  font-weight: bold;
  font-size: 1rem;
}
</style>

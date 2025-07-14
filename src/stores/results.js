import { ref, reactive } from 'vue';
import { defineStore } from 'pinia';

import { useNav } from '@/stores/misc.js';

import { getJsonPiped, getJsonHyp, getJsonAuth } from '@/scripts/fetch.js';
import { useVerifyAuth, useRoute, useTitle } from '@/scripts/util.js';

export const useResults = defineStore('results', () => {
  const items = ref({}),
    search = ref(''),
    chartsId = ref(''),
    album = ref(''),
    next = ref('');

  function setItem(key, val) {
    items.value[key] = val;
  }

  function resetItems() {
    next.value = undefined;
    album.value = undefined;

    useArtist().reset();

    for (let i in items.value) items.value[i] = undefined;
  }

  async function getExplore() {
    const json = await getJsonHyp('/explore');

    resetItems();

    chartsId.value = json.chartsId;

    setItem('songs', { items: json.trending });
    setItem('albums', { items: json.albums_and_singles });
  }

  async function getAlbum(e) {
    const hash = new URLSearchParams(e.substring(e.indexOf('?'))).get('list'),
      isAuth = useVerifyAuth(hash),
      path =
        '/playlists/' +
        (hash.startsWith('MPRE')
          ? (await getJsonHyp('/album/' + hash))?.id
          : hash),
      json = isAuth ? await getJsonAuth(path) : await getJsonPiped(path);

    resetItems();

    album.value = e;

    if (isAuth)
      json.relatedStreams = json.relatedStreams.map(i => {
        i.playlistId = hash;
        return i;
      });

    setItem('songs', {
      items: json.relatedStreams,
      title: json.name,
    });

    useTitle(json.name);
    useRoute('/playlist?list=' + hash);
    useNav().state.page = 'home';
    document.body.scrollIntoView();

    next.value =
      json.nextpage && json.nextpage != 'null'
        ? hash + '?nextpage=' + encodeURIComponent(json.nextpage)
        : null;
  }

  return {
    items,
    search,
    chartsId,
    next,
    album,
    setItem,
    resetItems,
    getExplore,
    getAlbum,
  };
});

export const useArtist = defineStore('artist', () => {
  const state = reactive({
    playlistId: null,
    title: null,
    description: null,
    subscriberCount: 0,
    hash: null,
    thumbnails: [],
  });

  const results = useResults();

  function reset() {
    for (let i in state) state[i] = undefined;
  }

  function set(obj) {
    for (let i in obj) {
      if (i in state) state[i] = obj[i];
    }
  }

  async function getArtist(e) {
    e = e.replace('/channel/', '');

    const json = await getJsonHyp('/channel/' + e);

    results.resetItems();

    for (let i in json.items)
      results.setItem(i, {
        items: json.items[i],
        more: json.more[i],
      });

    json.items = undefined;
    json.hash = e;

    reset();
    set(json);

    useTitle(json.title);
    useRoute('/channel/' + e);
    useNav().state.page = 'home';
    document.body.scrollIntoView();
  }

  async function getArtistNext(i, { id, params, click, visit }) {
    if (!id || !params || !click || !visit) return;

    const json = await getJsonHyp(
      `/next/channel/${id}/${params}?ct=${click}&v=${visit}`,
    );

    results.resetItems();
    reset();

    results.setItem(i, {
      items: json.items,
    });
  }

  return { state, set, reset, getArtist, getArtistNext };
});

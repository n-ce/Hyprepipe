import { reactive } from 'vue';
import { defineStore } from 'pinia';

import { getJsonPiped, getJsonHyp } from '@/scripts/fetch.js';
import { useStore, useMetadata, AMP } from '@/scripts/util.js';

const store = useStore();

export const useData = defineStore('data', () => {
  const state = reactive({
      title: '',
      description: '',
      artist: '',
      art: '',
      url: '',
      artistUrl: '',
      lyrics: '',
      src: [],
      urls: [],
    }),
    player = usePlayer();

  async function getSong(e) {
    const hash = new URLSearchParams(e.substring(e.indexOf('?'))).get('v'),
      json = await getJsonPiped('/streams/' + hash, {}, store.getItem('streamapi')),
      unamp = txt => txt?.replace(AMP, '&');

    state.art = unamp(json.thumbnailUrl);
    state.description = json.description;
    state.title = unamp(json.title);
    state.artist = unamp(json.uploader.replace(' - Topic', ''));
    state.artistUrl = json.uploaderUrl;
    player.state.currentTime = 0;
    player.setTime(0);
    player.state.duration = json.duration;
    player.state.hls = unamp(json.hls);
    player.state.streams = json.audioStreams;
    state.url = e;

    await getNext(hash);
  }

  async function play(song) {
    if (song.offlineUri) {
      state.art = song.thumbnail;
      player.state.duration = song.duration;
      for (let i of ['title', 'artist', 'artistUrl', 'url']) state[i] = song[i];
      window.audioPlayer.load(song.offlineUri);
    } else await getSong(song.url);
  }

  async function getNext(hash) {
    if (
      store.getItem('next') !== 'false' &&
      (!state.urls ||
        state.urls.findIndex(s => s.url == state.url) < 0 ||
        state.urls.length == 1)
    ) {
      const json = await getJsonHyp('/next/' + hash);

      state.lyrics = json.lyricsId;

      state.url = json.songs[0]
        ? '/watch?v=' + json.songs[0].id
        : '/watch?v=' + hash;

      json.songs[0] &&= { ...json.songs[0], ...{ lyrics: json.lyricsId } };

      state.urls =
        json.songs.length > 0
          ? json.songs.map(i => ({
              ...i,
              ...{
                url: '/watch?v=' + i.id,
                id: undefined,
              },
            }))
          : state.urls;
    } else if (state.urls.length == 0)
      state.urls = [
        {
          title: state.title,
          url: state.url,
        },
      ];

    useMetadata(current(), {
      title: state.title,
      artist: state.artist,
      art: state.art,
    });
  }

  function playNext() {
    const i = state.urls.findIndex(s => s.url === state.url);

    if (player.state.loop == 2) play(state.urls[i]);
    else if (
      state.urls.length > i &&
      state.urls.length != 0 &&
      state.urls[i + 1]
    )
      play(state.urls[i + 1]);
    else if (player.state.loop == 1) {
      state.url = state.urls[0].url;
      play(state.urls[0]);
    } else state.urls = [];
  }

  function prevTrack() {
    const i = state.urls.findIndex(s => s.url === state.url);

    if (state.urls[i - 1]) play(state.urls[i - 1]);
    else if (player.state.loop == 1) {
      state.url = state.urls.at(-1).url;
      play(state.urls.at(-1));
    } else state.urls = [];
  }

  function nextTrack() {
    const i = state.urls.findIndex(s => s.url === state.url);

    if (state.urls[i + 1]) play(state.urls[i + 1]);
    else if (player.state.loop == 1) {
      state.url = state.urls[0].url;
      play(state.urls[0]);
    } else state.urls = [];
  }

  function current() {
    return state.urls.find(i => i.url == state.url);
  }

  return { state, getSong, play, playNext, prevTrack, nextTrack, current };
});

export const usePlayer = defineStore('player', () => {
  const state = reactive({
    loop: 0,
    play: false,
    status: 'play',
    hls: '',
    streams: [],
    duration: 0,
    time: 0,
    realTime: 0,
    currentTime: 0,
    playlist: false,
    lyrics: false,
    info: false,
    add: false,
    vol: store.vol ? store.vol / 100 : 1,
    rate: parseFloat(store.rate || '1'),
    pitch: store.pitch != 'false',
  });

  function toggle(i) {
    if (typeof state[i] == 'boolean') state[i] = !state[i];
  }

  function setTime(t) {
    state.time = Math.floor((t / state.duration) * 100);
    state.realTime = Math.floor(t);
  }

  return { state, toggle, setTime };
});

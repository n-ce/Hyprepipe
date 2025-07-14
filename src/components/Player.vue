<script setup>
import { ref, watch, onMounted, onBeforeUnmount, onUnmounted } from 'vue';

import { useStore, useRoute, useManifest } from '@/scripts/util.js';
import { useData, usePlayer } from '@/stores/player.js';
import { useAlert } from '@/stores/misc';

let shaka;

const player = usePlayer(),
  data = useData(),
  store = useStore(),
  a = useAlert();

const audio = ref(null);

const dataOf = Object.freeze({
  1002: 1,
  3016: 0,
});

function audioCanPlay() {
  player.state.status = 'pause';
  audio.value.play().catch(err => {
    console.error(err);
    player.state.status = 'play';
  });

  if (location.pathname != '/playlist') useRoute(data.state.url);

  document.title = `Playing: ${data.state.title} by ${data.state.artist}`;
}

async function Stream() {
  const res = player.state;

  shaka ??= await import('shaka-player/dist/shaka-player.compiled.js').then(
    mod => mod.default,
  );

  const { url, mime } = await useManifest(res);

  if (!window.audioPlayer) {
    shaka.polyfill.installAll();

    if (shaka.Player.isBrowserSupported()) {
      const audioPlayer = new shaka.Player(),
        codecs = store.getItem('codec');

      await audioPlayer.attach(audio.value);
      audioPlayer.addEventListener('error', ({ detail }) => errorLog(detail));

      for (const [playerKey, audioKey] of Object.entries({
        rate: 'playbackRate',
        pitch: 'preservesPitch',
      }))
        audio.value[audioKey] = player.state[playerKey];

      audioPlayer.getNetworkingEngine().registerRequestFilter((_type, req) => {
        const headers = req.headers;

        let url = new URL(req.uris[0]);

        if (url.pathname.indexOf('/videoplayback') > -1) {
          if (headers.Range) {
            url.searchParams.set('range', headers.Range.split('=')[1]);
            req.headers = {};
            req.uris[0] = url.toString();
          }
        }
      });

      window.offline = new shaka.offline.Storage(audioPlayer);
      window.offline.configure({
        offline: {
          progressCallback: ({ appMetadata: { title, artist } }, prog) =>
            a.add(`${title} by ${artist}: ${Math.floor(prog * 100)}%`),
          trackSelectionCallback: tracks => [
            tracks
              .sort((a, b) => a.bandwidth - b.bandwidth)
              .find(i => i.type == 'variant'),
          ],
        },
      });

      audioPlayer.configure({
        preferredAudioCodecs: codecs ? codecs.split(':') : ['opus', 'mp4a'],
        manifest: {
          disableVideo: true,
          disableThumbnails: true,
          dash: {
            ignoreMinBufferTime: true,
          },
        },
        streaming: {
          segmentPrefetchLimit: 2,
          bufferBehind: 50,
          rebufferingGoal: 1,
          stallEnabled: true,
        },
      });

      window.audioPlayer = audioPlayer;
    }
  }

  const quality = store.getItem('quality');

  if (url)
    window.audioPlayer
      .load(url, 0, mime)
      .then(() => {
        window.audioPlayer.configure('abr.enabled', true);

        if (quality && quality != 'auto') {
          window.audioPlayer.configure('abr.enabled', false);

          const tracks = window.audioPlayer.getVariantTracks(),
            bandwidths = tracks.map(i => i.bandwidth);

          let sel;

          if (quality == 'best') sel = Math.max(...bandwidths);
          else if (quality == 'worst') sel = Math.min(...bandwidths);

          if (sel)
            window.audioPlayer.selectVariantTrack(
              tracks[bandwidths.indexOf(sel)],
            );
        }
      })
      .catch(errorLog);
}

function errorLog(err) {
  a.add(err.data[dataOf[err.code]] ?? 'Shaka Error: ' + err.code);
  console.error(err);
}

function destroy() {
  if (window.audioPlayer) {
    window.audioPlayer.destroy();
    window.audioPlayer = undefined;
  }
  if (window.offline) {
    window.offline.destroy();
    window.offline = undefined;
  }
}

const titleState = ['Playing', 'Paused'];

watch(
  () => player.state.play,
  () => {
    if (audio.value.paused) {
      player.state.status = 'pause';
      audio.value.play().catch(err => {
        console.error(err);
        player.state.status = 'play';
      });
    } else {
      player.state.status = 'play';
      audio.value.pause();
    }

    if (
      document.title.startsWith(titleState[0] + ':') ||
      document.title.startsWith(titleState[1] + ':')
    )
      document.title = audio.value.paused
        ? document.title.replace(titleState[0], titleState[1])
        : document.title.replace(titleState[1], titleState[0]);
  },
);

watch(() => player.state.streams, Stream);

watch(
  () => player.state.currentTime,
  cur => (audio.value.currentTime = cur),
);

watch(
  () => player.state.rate,
  cur => (audio.value.playbackRate = cur),
);

watch(
  () => player.state.pitch,
  cur => (audio.value.preservesPitch = cur),
);

onMounted(() => {
  Stream();

  if ('mediaSession' in navigator) {
    navigator.mediaSession.setActionHandler('play', () => {
      player.state.status = 'pause';
      document.title = document.title.replace(titleState[1], titleState[0]);

      audio.value.play().catch(err => {
        console.log(err);
        player.state.status = 'play';
        document.title = document.title.replace(titleState[0], titleState[1]);
      });
    });

    navigator.mediaSession.setActionHandler('pause', () => {
      audio.value.pause();
      player.state.status = 'play';
      document.title = document.title.replace(titleState[0], titleState[1]);
    });

    navigator.mediaSession.setActionHandler('previoustrack', () =>
      data.prevTrack(),
    );
    navigator.mediaSession.setActionHandler('nexttrack', () =>
      data.nextTrack(),
    );

    navigator.mediaSession.setActionHandler('seekbackward', () => {
      audio.value.currentTime -= 10;
    });

    navigator.mediaSession.setActionHandler('seekforward', () => {
      audio.value.currentTime += 10;
    });
  }
});

onBeforeUnmount(() => {
  destroy();
});
onUnmounted(() => {
  destroy();
});
</script>

<template>
  <audio
    id="audio"
    ref="audio"
    :volume="player.state.vol"
    @canplay="audioCanPlay"
    @timeupdate="player.setTime($event.target.currentTime)"
    @ended="data.playNext"
    autoplay></audio>
</template>

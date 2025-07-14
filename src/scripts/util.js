import { useI18n } from '@/stores/misc.js';

export const AMP = /&amp;/g;
export const PL_EXP =
  /[\da-fA-F]{8}-[\da-fA-F]{4}-[\da-fA-F]{4}-[\da-fA-F]{4}-[\da-fA-F]{12}/;

export function useSanitize(txt) {
  return ('' + txt)
    .replace(AMP, '&')
    .replace(/&/g, '&amp;')
    .replace(/>/g, '&gt;')
    .replace(/</g, '&lt;');
}

export function useAutoTheme(t) {
  return t == 'auto'
    ? matchMedia('(prefers-color-scheme: light)').matches
      ? 'light'
      : 'dark'
    : t;
}

export const parseThumb = (url, proxy) =>
  url && proxy
    ? `${proxy}/vi_webp/${url.replace(
        '/watch?v=',
        '',
      )}/maxresdefault.webp?host=i.ytimg.com`
    : '/1x1.png';

export function useVerifyAuth(hash) {
  return PL_EXP.test(hash);
}

export function useTitle(t) {
  document.title = `${t} - Hyperpipe`;
}

export function useRoute(l) {
  history.pushState({}, '', l);
}

export function useWrap(txt) {
  return encodeURIComponent(txt).replaceAll('%20', '+');
}

export function useUnwrap(txt) {
  return decodeURIComponent(txt).replaceAll('+', ' ');
}

export function useStore() {
  try {
    return window.localStorage;
  } catch (err) {
    console.error(err);

    const error = i => {
      console.error("Failed to Access '" + i + "' Please Enable localStorage.");
    };

    return {
      getItem: error,
      setItem: error,
    };
  }
}

export function useShare(data) {
  if ('share' in navigator) {
    navigator.share(data).catch(err => {
      console.error(err);
    });
  } else {
    const { t } = useI18n();
    navigator.clipboard.writeText(data.url).then(
      () => {
        alert(t('info.copied_to_clipboard'));
      },
      err => {
        alert(err);
      },
    );
  }
}

export function useMetadata(now, data) {
  if ('mediaSession' in navigator) {
    let artwork = [],
      album = undefined;

    if (now) {
      album = now.subtitle;

      if (now.thumbnails) {
        artwork = [...now.thumbnails].reverse().map(t => ({
          sizes: t.width && t.height ? t.width + 'x' + t.height : '512x512',
          src: t.url,
          type: 'image/webp',
        }));
      } else if (data.art) artwork = [{ src: data.art, type: 'image/webp' }];
    }

    navigator.mediaSession.metadata = new MediaMetadata({
      title: data.title,
      artist: data.artist,
      album: album,
      artwork: artwork,
    });
  }
}

export async function useManifest({ streams, duration, hls }) {
  let url, mime;

  const mse = window.MediaSource || window.ManagedMediaSource !== undefined;
  const prefer_hls = useStore().getItem('hls') == 'true';

  if (mse && streams.length > 0 && !prefer_hls) {
    const { useDash } = await import('./dash.js');
    const dash = useDash(streams, duration);

    url = 'data:application/dash+xml;charset=utf-8;base64,' + btoa(dash);
    mime = 'application/dash+xml';
  } else if (hls) {
    url = hls;
    mime = 'application/x-mpegURL';
  } else if (streams.length > 0) {
    url = streams[0].url;
    mime = streams[0].mimeType;
  }

  return { url, mime };
}

export function useShuffle(songs) {
  let copy = [],
    nos = songs.length;

  while (nos) {
    const i = Math.floor(Math.random() * nos--);

    copy.push(songs[i]);

    songs[i] = songs[nos];
    delete songs[nos];
  }

  return copy;
}

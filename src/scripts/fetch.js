import { useStore, useSanitize } from './util.js';
import { useAlert } from '@/stores/misc.js';

export const PIPED_INSTANCE = 'pipedapi.kavin.rocks';
export const HYPERPIPE_INSTANCE = 'hyperpipeapi.onrender.com';

export function getPipedQuery() {
  const papi = new URLSearchParams(location.search).get('pipedapi');

  if (!papi) return '';

  return '?pipedapi=' + useSanitize(papi);
}

export async function getJson(url, opts) {
  const errs = useAlert(),
    res = await fetch(url, opts)
      .then(res => res.json())
      .catch(err => {
        console.error(err);
        errs.add(err);
      });

  if (res && res.error) {
    errs.add(
      res.message
        ? res.message
            .replaceAll('Video', 'Audio')
            .replaceAll('video', 'audio')
            .replaceAll('watched', 'heard')
        : res.error,
    );

    console.error(res.message);
    return {};
  }

  return res;
}

export async function getJsonPiped(path, opts, root) {
  root ??=
    new URLSearchParams(location.search).get('pipedapi') ||
    useStore().getItem('pipedapi') ||
    PIPED_INSTANCE;

  return await getJson('https://' + root + path, opts);
}

export async function getJsonHyp(path) {
  const root = useStore().getItem('api') || HYPERPIPE_INSTANCE;

  return await getJson('https://' + root + path);
}

export async function getJsonAuth(path, opts = {}) {
  const root = useStore().getItem('authapi') || PIPED_INSTANCE;

  return await getJson('https://' + root + path, opts);
}

function useAuthToken() {
  return useStore().getItem('auth');
}

export async function useAuthCreatePlaylist(name) {
  const auth = useAuthToken();

  if (auth && name)
    return await getJsonAuth('/user/playlists/create', {
      method: 'POST',
      body: JSON.stringify({
        name: `Playlist - ${name}`,
      }),
      headers: {
        Authorization: auth,
        'Content-Type': 'application/json',
      },
    });
}

export async function getAuthPlaylists() {
  const auth = useAuthToken();

  if (auth) {
    const res = await getJsonAuth('/user/playlists', {
      headers: {
        Authorization: auth,
      },
    });

    return res.filter(i => i.name.startsWith('Playlist - '));
  }
}

export async function useAuthAddToPlaylist(id, path) {
  const auth = useAuthToken();

  if (auth && id && path)
    return getJsonAuth('/user/playlists/add', {
      method: 'POST',
      headers: {
        Authorization: auth,
      },
      body: JSON.stringify({
        playlistId: id,
        videoId: new URL('https://example.com' + path).searchParams.get('v'),
      }),
    });
}

export async function useAuthRemovePlaylist(id) {
  const auth = useAuthToken();

  if (auth && id) {
    return getJsonAuth('/user/playlists/delete', {
      method: 'POST',
      headers: {
        Authorization: auth,
      },
      body: JSON.stringify({
        playlistId: id,
      }),
    });
  }
}

export async function useAuthLogout() {
  const auth = useAuthToken(),
    ctrl = new AbortController(),
    id = setTimeout(() => ctrl.abort(), 1000);

  if (auth)
    return await getJsonAuth('/logout', {
      method: 'POST',
      signal: ctrl.signal,
      headers: {
        Authorization: auth,
      },
    }).then(res => {
      clearTimeout(id);
      return res;
    });
}

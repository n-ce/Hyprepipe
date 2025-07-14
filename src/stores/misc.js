import { reactive, ref } from 'vue';
import { defineStore } from 'pinia';

import en from '@/locales/en.json';

export const SUPPORTED_LOCALES = Object.freeze([
  {
    code: 'ar',
    name: 'عربي',
  },
  {
    code: 'az',
    name: 'Azərbaycan',
  },
  {
    code: 'bn',
    name: 'বাংলা',
    cc: 'IN',
  },
  {
    code: 'bs',
    name: 'Bosanski',
  },
  {
    code: 'cs',
    name: 'Čeština',
  },
  {
    code: 'de',
    name: 'Deutsch',
    cc: 'DE',
  },
  {
    code: 'en',
    name: 'English',
  },
  {
    code: 'es',
    name: 'Español',
  },
  {
    code: 'et',
    name: 'Eesti keel',
  },
  {
    code: 'fa',
    name: 'فارسی',
  },
  {
    code: 'fr',
    name: 'Français',
    cc: 'FR',
  },
  {
    code: 'gl',
    name: 'Galego',
  },
  {
    code: 'hr',
    name: 'Hrvatski',
  },
  {
    code: 'hu',
    name: 'Magyar nyelv',
  },
  {
    code: 'id',
    name: 'Bahasa Indonesia',
    cc: 'ID',
  },
  {
    code: 'it',
    name: 'Italiano',
    cc: 'IT',
  },
  {
    code: 'ja',
    name: '日本語',
    cc: 'JP',
  },
  {
    code: 'ko',
    name: '한국어',
    cc: 'KR',
  },
  {
    code: 'ml',
    name: 'മലയാളം',
    cc: 'IN',
  },
  {
    code: 'nb_NO',
    name: 'Norwegian Bokmål',
    cc: 'NO',
  },
  {
    code: 'nl',
    name: 'Nederlands',
    cc: 'NL',
  },
  {
    code: 'pt_pt',
    name: 'Português (Portugal)',
    cc: 'PT',
  },
  {
    code: 'pt_br',
    name: 'Português (Brasil)',
    cc: 'BR',
  },
  {
    code: 'ro',
    name: 'Română',
    cc: 'RO',
  },
  {
    code: 'ru',
    name: 'Pyccкий',
    cc: 'RU',
  },
  {
    code: 'sr',
    name: 'Српски',
    cc: 'RS',
  },
  {
    code: 'ta',
    name: 'தமிழ்',
    cc: 'IN',
  },
  {
    code: 'tr',
    name: 'Türkçe',
    cc: 'TR',
  },
  {
    code: 'uk',
    name: 'Yкраїнська',
    cc: 'UA',
  },
  {
    code: 'vi',
    name: 'Tiếng Việt',
  },
  {
    code: 'zh_Hans',
    name: '中文 (简体)',
    cc: 'CN',
  },
  {
    code: 'zh_Hant',
    name: '中文（繁體)',
  },
]);

export const useNav = defineStore('nav', () => {
  const state = reactive({
    search: '',
    page: 'home',
    show: false,
  });

  const show = () => (state.show = !state.show);

  return { state, show };
});

export const useI18n = defineStore('i18n', () => {
  const locale = ref('en'),
    map = ref({ en });

  function t(path) {
    const msgs = map.value[locale.value],
      fallback = map.value['en'],
      keys = path.split('.'),
      translate = msg => keys.reduce((obj, i) => obj?.[i], msg),
      translated = translate(msgs) || translate(fallback);

    return translated || path;
  }

  function setupLocale(code) {
    if (code === 'en') return setLocale('en');
    if (SUPPORTED_LOCALES.some(i => i.code == code))
      import(`@/locales/${code}.json`)
        .then(mod => mod.default)
        .then(mod => {
          map.value[code] = mod;
          setLocale(code);
        });
  }

  function setLocale(code) {
    locale.value = code;
    localStorage.locale ??= code;
    if (localStorage.to_cc == 'true') {
      const cc = SUPPORTED_LOCALES.find(i => i.code == code).cc;
      if (cc) localStorage.cc = cc;
      else localStorage.cc = '';
    } else localStorage.cc = '';
  }

  return { locale, map, t, setupLocale };
});

export const useAlert = defineStore('alert', () => {
  const msg = ref('');

  function add(m) {
    if (!m) return;

    msg.value = m;

    setTimeout(() => {
      if (msg.value == m) msg.value = '';
    }, 4000);
  }

  return { msg, add };
});

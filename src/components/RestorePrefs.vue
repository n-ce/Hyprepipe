<script setup>
import { useI18n } from '@/stores/misc.js';
import { ref, onMounted } from 'vue';

const { t } = useI18n();

const params = ref({}),
  restorePrefs = () => {
    for (const i in params.value)
      window.localStorage.setItem(i, params.value[i]);

    window.location.href = '/';
  };

onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search);
  urlParams.forEach((val, key) => {
    params.value[key] = val;
  });
});
</script>

<template>
  <h2>{{ t('title.restore_prefs') }}</h2>

  <table>
    <thead>
      <th>Key</th>
      <th>Value</th>
    </thead>
    <tbody>
      <tr v-for="(val, key) in params">
        <td>{{ key }}</td>
        <td>{{ val }}</td>
      </tr>
    </tbody>
  </table>

  <button @click="restorePrefs" class="input">
    {{ t('title.restore_prefs') }}
  </button>
</template>

<style>
h2 {
  text-align: center;
}

table {
  width: 100%;
  margin: 2rem 0;
}
td {
  width: 50%;
}
</style>

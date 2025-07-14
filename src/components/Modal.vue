<script setup>
import { ref, watch } from 'vue';

const props = defineProps(['display', 'title', 'n']),
  emit = defineEmits(['show']);

const show = ref(props.display);

watch(
  () => props.display,
  n => {
    console.log(n, props.display);
    show.value = n;
  },
);

watch(show, () => {
  emit('show', show.value);
});
</script>
<template>
  <Transition name="fade">
    <div class="modal" v-if="show">
      <span
        class="bi bi-x-lg clickable modal-close"
        @click="show = false"></span>
      <div class="modal-box">
        <div class="modal-title">{{ title }}</div>
        <div class="modal-content">
          <slot name="content"></slot>
        </div>
        <div class="modal-buttons">
          <slot name="buttons"></slot>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style>
.bi-x {
  font-size: 2rem;
}
.modal {
  display: flex;
  position: fixed;
  align-items: center;
  justify-content: center;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: #00000066;
  z-index: 99999;
}
.modal-box {
  width: 50vw;
  max-height: 90vh;
  border-radius: 0.5rem;
  background-color: var(--color-background-soft);
}
.modal-title {
  font-size: 1.25rem;
  padding: 1rem 2rem;
  border-bottom: 1px solid var(--color-shadow);
}
.modal-title::after {
  content: '...';
}
.modal-content {
  padding: 1rem;
  max-height: calc(90vh - 8rem);
  overflow-y: auto;
  overscroll-behavior: contain;
}
.modal-content * {
  width: 100%;
  padding: 0.5rem 1rem;
}
.modal-close {
  color: var(--color-background);
  font-size: 2rem;
  position: fixed;
  top: 1rem;
  right: 1rem;
}
.modal-buttons {
  width: 100%;
  border-top: 1px solid var(--color-shadow);
}
.modal-buttons button {
  padding: 1rem 2rem;
  color: var(--color-foreground);
  width: calc(100% / v-bind('n'));
  font-size: 1rem;
  font-weight: bold;
  border-right: 1px solid var(--color-shadow);
  transition: background-color 0.1s ease;
}
.modal-buttons button:first-child {
  color: indianred;
}
.modal-buttons button:first-child,
.modal-buttons button:first-child:focus {
  border-bottom-left-radius: 0.5rem;
}
.modal-buttons button:last-child {
  border: none;
}
.modal-buttons button:last-child,
.modal-buttons button:last-child:focus {
  border-bottom-right-radius: 0.5rem;
}
.modal-buttons button:hover {
  background-color: var(--color-background-mute);
}
.modal-buttons button:focus {
  background-color: var(--color-background-mute);
  border-radius: 0;
}
.modal-buttons button:active {
  background-color: var(--color-border);
}

@media (max-width: 530px) {
  .modal-box {
    width: 80vw;
  }
}
@media (min-width: 1024px) {
  .modal-box {
    width: 40vw;
  }
}
</style>

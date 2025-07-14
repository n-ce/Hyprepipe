<template>
  <div class="txt-modal">
    <slot name="content"></slot>
  </div>
</template>

<style scoped>
.txt-modal {
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 2rem;
  bottom: calc(5rem + 5vh);
  right: 1rem;
  width: 30rem;
  max-width: calc(100% - 2rem);
  background: var(--color-background-mute);
  border-radius: 0.5rem;
  z-index: 9999;
  box-shadow: 0.1rem 0.1rem 1rem var(--color-shadow);
  padding: 1rem;
  overflow-y: auto;
  overscroll-behavior: contain;
}
:deep(pre) {
  font-family: inherit;
  font-size: 1.25rem;
  letter-spacing: 0.125rem;
  white-space: pre-wrap;
  word-break: break-word;
}

:deep(div) {
  padding: 1rem;
  letter-spacing: 0.1rem;
  font-weight: 600;
}
</style>

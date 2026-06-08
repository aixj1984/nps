<script setup lang="ts">
defineProps<{
  open: boolean
  title: string
  canClose?: boolean
}>()

const emit = defineEmits<{
  close: []
}>()
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="desk-modal-overlay" @click.self="canClose && emit('close')">
      <div class="desk-modal" role="dialog" aria-modal="true" :aria-label="title">
        <header class="desk-modal-header">
          <h2>{{ title }}</h2>
          <button
            v-if="canClose"
            type="button"
            class="close-btn"
            aria-label="关闭"
            @click="emit('close')"
          >
            ×
          </button>
        </header>
        <div class="desk-modal-body">
          <slot />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.desk-modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(16, 22, 26, 0.65);
  display: grid;
  place-items: center;
  padding: 24px;
}

.desk-modal {
  width: min(720px, 100%);
  max-height: calc(100vh - 48px);
  background: #fff;
  border-radius: 6px;
  box-shadow:
    0 0 0 1px rgba(16, 22, 26, 0.1),
    0 18px 46px rgba(16, 22, 26, 0.35);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.desk-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
  background: linear-gradient(180deg, #f5f8fa, #fff);
}

.desk-modal-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--bp-blue);
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 3px;
  background: transparent;
  font-size: 22px;
  line-height: 1;
  color: var(--muted);
  cursor: pointer;
}

.close-btn:hover {
  background: rgba(16, 107, 163, 0.1);
  color: var(--bp-blue);
}

.desk-modal-body {
  overflow-y: auto;
  padding: 0;
}
</style>

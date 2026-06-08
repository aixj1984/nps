<script setup lang="ts">
import type { HostStatus } from '../../types'

defineProps<{
  status: HostStatus
  previewThumb?: string
}>()

const emit = defineEmits<{ confirm: [] }>()
</script>

<template>
  <div class="confirm-step">
    <p class="hint">请确认预览画面，点击开始后向查看端投屏。</p>

    <div class="preview card" v-if="previewThumb">
      <img :src="previewThumb" alt="preview" />
    </div>
    <div v-else class="preview card placeholder">已选择共享源</div>

    <p v-if="status.room.selectedSource?.name" class="source-name">
      共享源：<strong>{{ status.room.selectedSource.name }}</strong>
    </p>

    <div class="actions">
      <button type="button" class="bp-success confirm-btn" @click="emit('confirm')">开始共享</button>
    </div>
  </div>
</template>

<style scoped>
.confirm-step {
  padding: 16px 24px 28px;
  text-align: center;
}

.hint {
  margin: 0 0 16px;
  font-size: 14px;
  color: var(--muted);
}

.preview {
  max-width: 520px;
  margin: 0 auto 16px;
  padding: 8px;
}

.preview img {
  width: 100%;
  max-height: 220px;
  object-fit: contain;
  background: #000;
  border-radius: 2px;
}

.preview.placeholder {
  height: 140px;
  display: grid;
  place-items: center;
  color: var(--muted);
}

.source-name {
  font-size: 15px;
  margin-bottom: 20px;
}

.actions {
  display: flex;
  justify-content: center;
}

.confirm-btn {
  min-width: 200px;
  padding: 12px 28px !important;
  font-size: 16px !important;
}
</style>

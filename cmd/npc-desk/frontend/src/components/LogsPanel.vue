<script setup lang="ts">
defineProps<{
  logs: string
}>()

const emit = defineEmits<{
  refresh: []
  clear: []
}>()
</script>

<template>
  <div class="card logs-panel">
    <div class="panel-head">
      <h2>运行日志</h2>
      <div class="actions">
        <button type="button" class="bp-minimal" @click="emit('refresh')">刷新</button>
        <button type="button" class="bp-minimal" @click="emit('clear')">清空</button>
      </div>
    </div>
    <pre>{{ logs || '暂无日志' }}</pre>
  </div>
</template>

<style scoped>
.logs-panel {
  margin-bottom: 16px;
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.actions {
  display: flex;
  gap: 8px;
}

h2 {
  margin: 0;
  font-size: 18px;
  color: var(--bp-blue);
}

pre {
  margin: 0;
  min-height: 160px;
  max-height: 280px;
  overflow: auto;
  padding: 12px;
  background: #182026;
  color: #d8e1e8;
  border-radius: 3px;
  font-size: 12px;
  line-height: 1.45;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>

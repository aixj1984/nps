<script setup lang="ts">
import type { HostStatus } from '../types'

defineProps<{
  status: HostStatus
}>()

const emit = defineEmits<{
  allow: []
  deny: []
}>()
</script>

<template>
  <div class="approval">
    <p class="lead">有设备请求将本机屏幕投屏到其浏览器，是否允许？</p>

    <table class="device-table">
      <tr><td>设备类型</td><td>{{ status.room.pendingDevice!.deviceType || '—' }}</td></tr>
      <tr><td>系统</td><td>{{ status.room.pendingDevice!.os || '—' }}</td></tr>
      <tr><td>浏览器</td><td>{{ status.room.pendingDevice!.browser || '—' }}</td></tr>
      <tr><td>IP</td><td><code>{{ status.room.pendingDevice!.ip }}</code></td></tr>
      <tr>
        <td>屏幕</td>
        <td>
          {{ status.room.pendingDevice!.screenWidth || '?' }}×{{
            status.room.pendingDevice!.screenHeight || '?'
          }}
        </td>
      </tr>
    </table>

    <div class="actions">
      <button type="button" class="bp-success" @click="emit('allow')">允许</button>
      <button type="button" class="bp-danger" @click="emit('deny')">拒绝</button>
    </div>
  </div>
</template>

<style scoped>
.approval {
  padding: 20px 24px 24px;
}

.lead {
  margin: 0 0 16px;
  font-size: 14px;
  color: var(--text);
  line-height: 1.5;
}

.device-table {
  width: 100%;
  font-size: 14px;
  margin-bottom: 24px;
}

.device-table td {
  padding: 8px 10px;
  border-bottom: 1px solid var(--border);
}

.device-table td:first-child {
  color: var(--muted);
  width: 88px;
}

.actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.actions .bp-success {
  min-width: 120px;
}
</style>

<script setup lang="ts">
import type { NpcStatus } from '../types'

defineProps<{
  status: NpcStatus
  busy?: boolean
}>()

const emit = defineEmits<{
  start: []
  stop: []
  pickConfig: []
  openFolder: []
}>()
</script>

<template>
  <div class="card connection-panel">
    <div class="panel-head">
      <h2>连接状态</h2>
      <div class="actions">
        <button type="button" class="bp-minimal" @click="emit('pickConfig')">打开配置</button>
        <button type="button" class="bp-minimal" @click="emit('openFolder')">配置目录</button>
      </div>
    </div>

    <div class="status-grid">
      <div class="stat">
        <span class="label">服务器</span>
        <strong>{{ status.serverAddr || '—' }}</strong>
      </div>
      <div class="stat">
        <span class="label">连接类型</span>
        <strong>{{ status.connType || 'tcp' }}</strong>
      </div>
      <div class="stat">
        <span class="label">验证密钥</span>
        <strong>{{ status.vkey || '—' }}</strong>
      </div>
      <div class="stat">
        <span class="label">自动重连</span>
        <strong>{{ status.autoReconnection ? '开启' : '关闭' }}</strong>
      </div>
    </div>

    <p class="config-path">配置文件：<code>{{ status.configPath || '未保存' }}</code></p>
    <p v-if="status.lastError" class="error">{{ status.lastError }}</p>

    <div class="control-row">
      <button
        v-if="!status.running"
        type="button"
        class="bp-primary"
        :disabled="busy"
        @click="emit('start')"
      >
        启动 NPC
      </button>
      <button v-else type="button" class="bp-danger" :disabled="busy" @click="emit('stop')">停止 NPC</button>
    </div>
  </div>
</template>

<style scoped>
.connection-panel {
  margin-bottom: 16px;
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

h2 {
  margin: 0;
  font-size: 18px;
  color: var(--bp-blue);
}

.actions {
  display: flex;
  gap: 8px;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.stat {
  padding: 12px;
  background: #f5f8fa;
  border: 1px solid var(--border);
  border-radius: 3px;
}

.label {
  display: block;
  font-size: 12px;
  color: var(--muted);
  margin-bottom: 4px;
}

.config-path {
  margin: 16px 0 0;
  font-size: 13px;
  color: var(--muted);
  word-break: break-all;
}

code {
  font-family: Consolas, monospace;
  color: var(--bp-blue);
}

.error {
  margin: 12px 0 0;
  padding: 10px 12px;
  background: rgba(219, 55, 55, 0.1);
  border: 1px solid rgba(219, 55, 55, 0.25);
  color: var(--bp-red);
  border-radius: 3px;
  font-size: 13px;
}

.control-row {
  margin-top: 16px;
  display: flex;
  gap: 10px;
}
</style>

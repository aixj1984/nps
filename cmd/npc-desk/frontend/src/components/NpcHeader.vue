<script setup lang="ts">
import type { NpcStatus } from '../types'

defineProps<{
  status: NpcStatus | null
}>()
</script>

<template>
  <header class="host-header">
    <div class="brand">
      <div class="logo">N</div>
      <div>
        <h1>NPC Desk</h1>
        <p class="tag">内网穿透客户端 · {{ status?.version || '—' }}</p>
      </div>
    </div>
    <div class="meta">
      <span
        class="badge"
        :class="status?.connected ? 'ok' : status?.reconnecting ? 'warn' : status?.running ? 'warn' : 'idle'"
      >
        {{
          status?.connected
            ? '已连接'
            : status?.reconnecting
              ? '重连中'
              : status?.running
                ? '启动中'
                : '未连接'
        }}
      </span>
      <span class="port">{{ status?.serverAddr || '未配置服务器' }}</span>
    </div>
  </header>
</template>

<style scoped>
.host-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background: linear-gradient(180deg, #137cbd, #106ba3);
  color: #fff;
  box-shadow: 0 1px 0 rgba(16, 22, 26, 0.2);
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.2);
  display: grid;
  place-items: center;
  font-weight: 800;
  font-size: 1.25rem;
}

h1 {
  margin: 0;
  font-size: 1.35rem;
  font-weight: 600;
}

.tag {
  margin: 2px 0 0;
  font-size: 0.75rem;
  opacity: 0.85;
}

.meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
}

.port {
  font-size: 0.8rem;
  opacity: 0.9;
  font-family: Consolas, monospace;
}

.host-header .badge.ok {
  background: rgba(15, 153, 96, 0.35);
  color: #fff;
}

.host-header .badge.warn {
  background: rgba(217, 130, 43, 0.45);
  color: #fff;
}

.host-header .badge.idle {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}
</style>

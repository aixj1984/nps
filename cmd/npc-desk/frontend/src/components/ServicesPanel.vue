<script setup lang="ts">
import { computed } from 'vue'
import type { ServiceItem } from '../types'

const props = defineProps<{
  services: ServiceItem[]
  connected: boolean
  activeForwardCount: number
  totalForwardCount: number
  recentlyForwarded: boolean
}>()

const kindLabel: Record<string, string> = {
  tunnel: '隧道',
  host: '域名',
  local: '本地',
  health: '健康检查',
}

const isForwarding = computed(
  () => props.activeForwardCount > 0 || props.recentlyForwarded,
)

const channelLabel = computed(() => {
  if (props.activeForwardCount > 0) {
    return `转发活跃 (${props.activeForwardCount})`
  }
  if (props.recentlyForwarded) {
    return '最近活跃'
  }
  if (props.connected) {
    return '等待连接'
  }
  return '未连接'
})

const channelClass = computed(() => {
  if (isForwarding.value) {
    return 'ok'
  }
  if (props.connected) {
    return 'idle'
  }
  return 'warn'
})

const rowStatusLabel = (active: boolean) => {
  if (active && isForwarding.value) {
    return props.activeForwardCount > 0 ? '转发中' : '最近活跃'
  }
  if (props.connected) {
    return '空闲'
  }
  return '离线'
}
</script>

<template>
  <div class="card services-panel">
    <div class="panel-head">
      <h2>服务映射</h2>
      <div class="channel-status">
        <span class="count">转发数 {{ activeForwardCount }} · 累计 {{ totalForwardCount }}</span>
        <span class="badge" :class="channelClass">{{ channelLabel }}</span>
      </div>
    </div>

    <p class="hint">
      统计包含 TCP 隧道与 HTTP/HTTPS 域名转发；HTTP 请求通常极短，「累计」可持续观察是否有访问经过客户端。
    </p>

    <p v-if="services.length === 0" class="empty">当前配置中没有隧道、域名或本地服务段落。</p>

    <div v-else class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>名称</th>
            <th>类型</th>
            <th>模式</th>
            <th>端口/域名</th>
            <th>目标</th>
            <th>状态</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="svc in services" :key="svc.name + svc.kind + svc.port">
            <td>{{ svc.name }}</td>
            <td>{{ kindLabel[svc.kind] || svc.kind }}</td>
            <td>{{ svc.mode || '—' }}</td>
            <td><code>{{ svc.port || '—' }}</code></td>
            <td class="target">{{ svc.target || '—' }}</td>
            <td>
              <span class="badge" :class="svc.active ? 'ok' : connected ? 'idle' : 'warn'">
                {{ rowStatusLabel(svc.active) }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.services-panel {
  margin-bottom: 16px;
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

h2 {
  margin: 0;
  font-size: 18px;
  color: var(--bp-blue);
}

.channel-status {
  display: flex;
  align-items: center;
  gap: 10px;
}

.count {
  font-size: 13px;
  color: var(--muted);
  font-variant-numeric: tabular-nums;
}

.hint {
  margin: 0 0 12px;
  font-size: 12px;
  color: var(--muted);
}

.empty {
  color: var(--muted);
  font-size: 14px;
}

.table-wrap {
  overflow: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

th,
td {
  text-align: left;
  padding: 10px 8px;
  border-bottom: 1px solid var(--border);
}

th {
  color: var(--muted);
  font-weight: 600;
}

.target {
  max-width: 240px;
  word-break: break-all;
}

code {
  font-family: Consolas, monospace;
}

.badge.idle {
  background: rgba(92, 112, 128, 0.12);
  color: #5c7080;
  border: 1px solid rgba(92, 112, 128, 0.2);
}
</style>

<script setup lang="ts">
import { ref } from 'vue'
import type { HostStatus } from '../../types'

defineProps<{
  status: HostStatus
  qrDataUrl: string
  locked: boolean
  dimmed?: boolean
}>()

const emit = defineEmits<{
  copy: []
  open: []
}>()

const showQrLarge = ref(false)
</script>

<template>
  <div class="connect" :class="{ dimmed }">
    <div v-if="!status.lanAvailable" class="alert warn">
      Deskreen 仅在 WiFi / 局域网环境下工作。请检查网络连接。
    </div>

    <p v-if="locked" class="alert info">
      已有查看端占用（CE 单台）。请点击底部「断开全部连接」后重试。
    </p>

    <div class="connect-grid" :class="{ dim: locked }">
      <div class="qr-panel card">
        <h3>扫描二维码</h3>
        <div v-if="qrDataUrl" class="qr-box" @click="showQrLarge = true">
          <img :src="qrDataUrl" alt="QR" width="200" height="200" />
        </div>
        <p v-else class="hint">正在生成二维码…</p>
        <button type="button" class="bp-minimal" :disabled="locked || !qrDataUrl" @click="showQrLarge = true">
          放大二维码
        </button>
      </div>
      <div class="url-panel card">
        <h3>或打开链接</h3>
        <p class="url">{{ status.viewerUrl }}</p>
        <p class="hint">连接后会出现弹框：先「允许设备」，再「选择投屏窗口」</p>
        <div class="actions">
          <button type="button" class="bp-primary" :disabled="locked" @click="emit('copy')">复制链接</button>
          <button type="button" class="bp-minimal" :disabled="locked" @click="emit('open')">在浏览器中打开</button>
        </div>
        <p class="room-id">房间 ID: <code>{{ status.room.roomId }}</code></p>
      </div>
    </div>

    <div v-if="showQrLarge && qrDataUrl" class="modal" @click.self="showQrLarge = false">
      <div class="modal-body">
        <img :src="qrDataUrl" alt="QR" width="320" height="320" />
        <button type="button" class="bp-minimal" @click="showQrLarge = false">关闭</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.connect {
  padding: 20px;
  transition: opacity 0.2s;
}

.connect.dimmed {
  opacity: 0.35;
  pointer-events: none;
}

.alert {
  padding: 12px 16px;
  border-radius: 3px;
  margin-bottom: 16px;
  font-size: 14px;
}

.alert.warn {
  background: rgba(217, 130, 43, 0.15);
  border: 1px solid rgba(217, 130, 43, 0.35);
  color: #a66321;
}

.alert.info {
  background: rgba(16, 107, 163, 0.1);
  border: 1px solid rgba(16, 107, 163, 0.25);
}

.connect-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.connect-grid.dim {
  opacity: 0.45;
  pointer-events: none;
}

h3 {
  margin: 0 0 12px;
  font-size: 16px;
  color: var(--bp-blue);
}

.qr-box {
  display: flex;
  justify-content: center;
  padding: 12px;
  background: #fff;
  border: 1px solid var(--border);
  cursor: pointer;
}

.url {
  font-family: Consolas, monospace;
  font-size: 13px;
  word-break: break-all;
  padding: 10px;
  background: #f5f8fa;
  border: 1px solid var(--border);
  border-radius: 3px;
  color: var(--bp-blue);
}

.hint {
  font-size: 13px;
  color: var(--muted);
  margin: 12px 0;
}

.room-id {
  font-size: 12px;
  color: var(--muted);
  margin-top: 16px;
}

.actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.modal {
  position: fixed;
  inset: 0;
  background: rgba(16, 22, 26, 0.6);
  display: grid;
  place-items: center;
  z-index: 500;
}

.modal-body {
  background: #fff;
  padding: 24px;
  border-radius: 6px;
  text-align: center;
}
</style>
